package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/dop251/goja"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"math/rand"
	"strings"
)

type kfrCtxKey struct{}

type KfrValues struct {
	Provider   string
	Shoot      string
	LoadedCRDs map[string]struct{}
	Env        string
}

type KfrContext struct {
	Resources map[string]*ResourceDefn
	K8S       K8sClient
	Values    KfrValues
}

var NotFoundError = errors.New("not found")

type ResourceDefn struct {
	Var       string
	Kind      string
	Name      string
	Namespace string
	Value     map[string]interface{}

	NamesEvaluated bool
	KfrCtx         *KfrContext
}

func KfrToContext(ctx context.Context, kfr *KfrContext) context.Context {
	return context.WithValue(ctx, kfrCtxKey{}, kfr)
}

func KfrFromContext(ctx context.Context) *KfrContext {
	g, _ := ctx.Value(kfrCtxKey{}).(*KfrContext)
	return g
}

func (k *KfrContext) Get(varName string) *ResourceDefn {
	r, ok := k.Resources[varName]
	if !ok {
		return nil
	}
	return r
}

func (k *KfrContext) Set(varName string, kind string) *ResourceDefn {
	d := &ResourceDefn{
		Var:    varName,
		Kind:   kind,
		KfrCtx: k,
	}
	k.Resources[varName] = d
	return d
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func (k *KfrContext) Eval(ctx context.Context, exp string) (any, error) {
	vm := goja.New()
	if err := vm.GlobalObject().Set("namespace", k.Namespace()); err != nil {
		return nil, fmt.Errorf("error registering namespace: %w", err)
	}
	for _, rd := range k.Resources {
		if err := vm.GlobalObject().Set(rd.Var, rd.Value); err != nil {
			return nil, fmt.Errorf("error registering resource %s: %w", rd.Var, err)
		}
	}

	if err := vm.GlobalObject().Set("shoot", k.Values.Shoot); err != nil {
		return nil, fmt.Errorf("error registering shoot: %w", err)
	}
	if err := vm.GlobalObject().Set("provider", k.Values.Provider); err != nil {
		return nil, fmt.Errorf("error registering provider: %w", err)
	}
	if err := vm.GlobalObject().Set("env", k.Values.Env); err != nil {
		return nil, fmt.Errorf("error registering provider: %w", err)
	}

	if err := vm.GlobalObject().Set("declare", func(ref, kind, name, namespace string, r *goja.Runtime) (any, error) {
		if ref == "" {
			return nil, errors.New("declare() requires mandatory first argument resource declatation name")
		}
		if kind == "" {
			return nil, errors.New("declare() requires mandatory second argument kind")
		}
		if _, exists := k.Resources[ref]; exists {
			return nil, fmt.Errorf("resource %s already defined", ref)
		}
		rd := k.Set(ref, kind)
		rd.Name = name
		rd.Namespace = namespace
		return nil, nil
	}); err != nil {
		return nil, fmt.Errorf("error registering declare: %w", err)
	}

	if err := vm.GlobalObject().Set("load", func(ref string, r *goja.Runtime) (map[string]interface{}, error) {
		rd := k.Get(ref)
		if rd == nil {
			chunks := strings.Split(ref, "/")
			if len(chunks) <= 1 {
				return nil, fmt.Errorf("resource %s is not registered", ref)
			}
			if len(chunks) > 3 {
				return nil, fmt.Errorf("resource %s malformed name spec not matching kind(/namespace)?/name", ref)
			}
			rd = k.Set(ref, chunks[0])
			if len(chunks) == 2 {
				rd.Name = chunks[1]
			} else {
				rd.Namespace = chunks[1]
				rd.Name = chunks[2]
			}
			rd.NamesEvaluated = true
		}
		err := rd.Reload(ctx)
		if errors.Is(err, NotFoundError) {
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("error loading resource %s: %w", ref, err)
		}
		return rd.Value, nil
	}); err != nil {
		return nil, fmt.Errorf("error registering load(): %w", err)
	}

	if err := vm.GlobalObject().Set("logs", func(ref string, r *goja.Runtime) (string, error) {
		rd := k.Get(ref)
		if rd == nil {
			return "", fmt.Errorf("resource %s is not registered", ref)
		}
		txt, err := rd.Logs(ctx)
		if errors.Is(err, NotFoundError) {
			return "", nil
		}
		if err != nil {
			return "", fmt.Errorf("error loading resource %s logs: %w", ref, err)
		}
		return txt, nil
	}); err != nil {
		return nil, fmt.Errorf("error registering logs(): %w", err)
	}

	if err := vm.GlobalObject().Set("rndStr", func(l int, r *goja.Runtime) (string, error) {
		if l < 1 || l > 100 {
			l = 8
		}
		val := randString(l)

		return val, nil
	}); err != nil {
		return nil, fmt.Errorf("error registering rndStr(): %w", err)
	}

	if err := vm.GlobalObject().Set("apply", func(ref string, obj any) (any, error) {
		var doc string
		if txt, ok := obj.(string); ok {
			doc = txt
		} else if o, ok := obj.(map[string]interface{}); ok {
			b, err := yaml.Marshal(o)
			if err != nil {
				return nil, fmt.Errorf("error marshalling resource %s to yaml: %w", ref, err)
			}
			doc = string(b)
		} else {
			return nil, fmt.Errorf("unsupported resource type %T", obj)
		}

		return nil, resourceApplied(ctx, ref, &godog.DocString{Content: doc})
	}); err != nil {
		return nil, fmt.Errorf("error registering apply(): %w", err)
	}

	val, err := vm.RunString(exp)
	if err != nil {
		return nil, fmt.Errorf("error evaluating script %s: %w", exp, err)
	}

	return val.Export(), nil
}

func (k *KfrContext) Namespace() string {
	// TODO: namespace should be configurable
	return "default"
}

func (rd *ResourceDefn) EvaluateNames(ctx context.Context) error {
	if len(rd.Name) == 0 {
		return errors.New("resource name is empty")
	}
	if rd.NamesEvaluated {
		return nil
	}

	// Name evaluation
	v, err := rd.KfrCtx.Eval(ctx, rd.Name)
	if err != nil {
		return fmt.Errorf("error evaluating resource name %s: %w", rd.Name, err)
	}
	if v == nil {
		return errors.New("resource name evaluated to nil")
	}
	vv := fmt.Sprintf("%v", v)
	if len(vv) == 0 {
		return errors.New("resource name evaluated to empty string")
	}
	rd.Name = vv

	// Namespace evaluation
	if len(rd.Namespace) > 0 {
		v, err := rd.KfrCtx.Eval(ctx, rd.Namespace)
		if err != nil {
			return fmt.Errorf("error evaluating resource namespace %s: %w", rd.Namespace, err)
		}
		if v == nil {
			return errors.New("resource namespace evaluated to nil")
		}
		vv := fmt.Sprintf("%v", v)
		rd.Namespace = vv
	}

	rd.NamesEvaluated = true

	return nil
}

func (rd *ResourceDefn) Logs(ctx context.Context) (string, error) {
	if err := rd.EvaluateNames(ctx); err != nil {
		return "", fmt.Errorf("resource logs can not be loaded due to name evaluation error: %w", err)
	}

	out, err := rd.KfrCtx.K8S.Logs(ctx, rd.Name, rd.Namespace)
	if errors.Is(err, NotFoundError) {
		rd.Value = nil
		return "", err
	}
	if err != nil {
		return "", err
	}

	return out, nil
}

func (rd *ResourceDefn) Reload(ctx context.Context) error {
	if err := rd.EvaluateNames(ctx); err != nil {
		return fmt.Errorf("resource can not be loaded due to name evaluation error: %w", err)
	}

	val, err := rd.KfrCtx.K8S.Get(ctx, rd.Kind, rd.Name, rd.Namespace)
	if errors.Is(err, NotFoundError) {
		rd.Value = nil
		return err
	}
	if err != nil {
		return err
	}
	rd.Value = val

	return nil
}

func (rd *ResourceDefn) Apply(ctx context.Context) error {
	if err := rd.EvaluateNames(ctx); err != nil {
		return fmt.Errorf("resource can not be applied due to name evaluation error: %w", err)
	}
	if rd.Value == nil {
		return errors.New("resource has no value and can not be applied")
	}

	b, err := yaml.Marshal(rd.Value)
	if err != nil {
		return fmt.Errorf("error marshalling resource %s to yaml: %w", rd.Var, err)
	}

	if err := rd.KfrCtx.K8S.Apply(ctx, string(b)); err != nil {
		return fmt.Errorf("error applying resource %s: %w", rd.Var, err)
	}

	return nil
}

func (rd *ResourceDefn) HasValue() bool {
	return rd.Value != nil
}

func (rd *ResourceDefn) ExtractNames() error {
	if !rd.HasValue() {
		return errors.New("resource has no value")
	}

	n, found, err := unstructured.NestedString(rd.Value, "metadata", "name")
	if err != nil {
		return fmt.Errorf("error extracting resource name: %w", err)
	}
	if !found {
		return errors.New("resource has no metadata name")
	}
	rd.Name = n

	n, found, err = unstructured.NestedString(rd.Value, "metadata", "namespace")
	if err != nil {
		return fmt.Errorf("error extracting resource namespace: %w", err)
	}
	if found {
		rd.Namespace = n
	} else {
		rd.Namespace = rd.KfrCtx.Namespace()
	}

	rd.NamesEvaluated = true

	return nil
}
