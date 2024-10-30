package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"regexp"
	"strings"
)

var appliedRegex = regexp.MustCompile(`<\(.+\)>`)

func resourceApplied(ctx context.Context, ref string, txt *godog.DocString) error {
	kfrCtx := KfrFromContext(ctx)
	var rd *ResourceDefn
	val := map[string]interface{}{}
	var content string

	if txt != nil {
		content = txt.Content

		var err error
		content = string(appliedRegex.ReplaceAllFunc([]byte(content), func(b []byte) []byte {
			e := strings.TrimPrefix(string(b), "<(")
			e = strings.TrimSuffix(e, ")>")
			v, errE := kfrCtx.Eval(ctx, e)
			if errE != nil {
				err = fmt.Errorf(`error evaluating "%s": %w`, e, errE)
				return nil
			}
			return []byte(fmt.Sprintf("%v", v))
		}))
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal([]byte(content), &val); err != nil {
			return fmt.Errorf("invalid yaml: %w", err)
		}
	} else if len(ref) > 0 {
		rd = kfrCtx.Get(ref)
		if rd == nil {
			return fmt.Errorf("resource %s not declated and no yaml document provided", ref)
		}
		val = rd.Value
	}

	kind, _, err := unstructured.NestedString(val, "kind")
	if err != nil {
		return fmt.Errorf("missing kind: %w", err)
	}
	if len(kind) == 0 {
		return errors.New("zero length kind")
	}

	// set context label to all resources
	err = unstructured.SetNestedField(val, map[string]interface{}{"context": "cloud-manager-tests"}, "metadata", "labels")

	if len(ref) > 0 {
		rd = kfrCtx.Get(ref)
		if rd == nil {
			rd = kfrCtx.Set(ref, kind)
		} else if rd.Kind == "" {
			rd.Kind = kind
		} else if rd.Kind != kind {
			return fmt.Errorf("appied kind %s is different then registered kind %s", kind, rd.Kind)
		}
		rd.Value = val

		if len(rd.Name) == 0 {
			// name is not specified in resource declaration, then it must be in the specified yaml
			// manifest that will be applied
			if err = rd.ExtractNames(); err != nil {
				return fmt.Errorf("error finding resource name: %w", err)
			}
		} else {
			// name is specified in the resource declaration, then it must be evaluated and
			// injected into the specified yaml manifest before it's applied
			if err = rd.EvaluateNames(ctx); err != nil {
				return fmt.Errorf("resource %s name %s can not be evaluated: %w", rd.Var, rd.Name, err)
			}

			err = unstructured.SetNestedField(val, rd.Name, "metadata", "name")
			if err != nil {
				return fmt.Errorf("error setting resource name: %w", err)
			}
			if len(rd.Namespace) > 0 {
				err = unstructured.SetNestedField(val, rd.Namespace, "metadata", "namespace")
				if err != nil {
					return fmt.Errorf("error setting resource name: %w", err)
				}
			}

			b, err := yaml.Marshal(val)
			if err != nil {
				return fmt.Errorf("error marshalling resource: %w", err)
			}
			content = string(b)
		}

	} else {
		n, f, err := unstructured.NestedString(val, "metadata", "name")
		if err != nil {
			return fmt.Errorf("error finding resource name: %w", err)
		}
		if !f || len(n) == 0 {
			return errors.New("resource name is missing")
		}

		n, f, err = unstructured.NestedString(val, "metadata", "namespace")
		if err != nil {
			return fmt.Errorf("error finding resource namespace: %w", err)
		}

		if !f || len(n) == 0 {
			err = unstructured.SetNestedField(val, kfrCtx.Namespace(), "metadata", "namespace")
			if err != nil {
				return fmt.Errorf("error setting resource namespace: %w", err)
			}
			b, err := yaml.Marshal(val)
			if err != nil {
				return fmt.Errorf("error marshalling yaml with namespace set: %w", err)
			}
			content = string(b)
		}
	}

	return kfrCtx.K8S.Apply(ctx, content)
}
