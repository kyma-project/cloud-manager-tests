package internal

import (
	"context"
	"fmt"
	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func moduleRemoved(ctx context.Context, moduleName string) error {
	if moduleName == "" {
		moduleName = "cloud-manager"
	}

	rd, err := getReloadedKyma(ctx)
	if err != nil {
		return err
	}

	status, err := getKymaModuleStateInStatus(moduleName, rd.Value)
	if err != nil {
		return err
	}
	if status == "" {
		return nil
	}

	if err := removeKymaModuleInSpec(moduleName, rd.Value); err != nil {
		return err
	}
	if err := rd.Apply(ctx); err != nil {
		return err
	}

	var errMsg string
	gm := gomega.NewGomega(func(message string, callerSkip ...int) {
		errMsg = message
	})
	ok := gm.Eventually(func(ctx context.Context) error {
		if err := rd.Reload(ctx); err != nil {
			return err
		}
		status, err := getKymaModuleStateInStatus(moduleName, rd.Value)
		if err != nil {
			return err
		}
		if status == "" {
			return nil
		}
		return fmt.Errorf("module still exists: %q", status)
	}).
		WithArguments(ctx).
		Should(gomega.Succeed())

	if !ok || len(errMsg) > 0 {
		return fmt.Errorf("failed removing module: %s", errMsg)
	}
	return nil
}

func moduleAdded(ctx context.Context, moduleName string) error {
	if moduleName == "" {
		moduleName = "cloud-manager"
	}

	rd, err := getReloadedKyma(ctx)
	if err != nil {
		return err
	}

	status, err := getKymaModuleStateInStatus(moduleName, rd.Value)
	if err != nil {
		return err
	}
	if status == "Ready" {
		return nil
	}

	if err := addKymaModuleInSpec(moduleName, rd.Value); err != nil {
		return err
	}
	if err := rd.Apply(ctx); err != nil {
		return err
	}

	var errMsg string
	gm := gomega.NewGomega(func(message string, callerSkip ...int) {
		errMsg = message
	})
	ok := gm.Eventually(func(ctx context.Context) error {
		if err := rd.Reload(ctx); err != nil {
			return err
		}
		status, err := getKymaModuleStateInStatus(moduleName, rd.Value)
		if err != nil {
			return err
		}
		if status == "Ready" {
			return nil
		}
		return fmt.Errorf("module not ready: %q", status)
	}).
		WithArguments(ctx).
		Should(gomega.Succeed())

	if !ok || len(errMsg) > 0 {
		return fmt.Errorf("failed removing module: %s", errMsg)
	}
	return nil
}

func getReloadedKyma(ctx context.Context) (*ResourceDefn, error) {
	kfrCtx := KfrFromContext(ctx)
	rd := kfrCtx.Get("kyma")
	if rd == nil {
		rd = kfrCtx.Set("kyma", "Kyma")
		rd.Namespace = "kyma-system"
		rd.Name = "default"
		rd.NamesEvaluated = true
	}

	if err := rd.Reload(ctx); err != nil {
		return nil, fmt.Errorf("error loading kyma: %w", err)
	}

	return rd, nil
}

func getKymaModuleIndexInSpec(moduleName string, obj map[string]interface{}) (int, []interface{}, error) {
	modules, exists, err := unstructured.NestedSlice(obj, "spec", "modules")
	if !exists || err != nil {
		return -1, nil, fmt.Errorf("error reading kyma spec modules: %w", err)
	}

	index := -1
	for i, m := range modules {
		mm, ok := m.(map[string]interface{})
		if !ok {
			return -1, nil, fmt.Errorf("kyma spec module is not a map, but: %T", m)
		}
		name, exists, err := unstructured.NestedString(mm, "name")
		if !exists || err != nil {
			return -1, nil, fmt.Errorf("kyma spec module does not have a name: %w", err)
		}

		if name == moduleName {
			index = i
			break
		}
	}

	return index, modules, nil
}

func removeKymaModuleInSpec(moduleName string, obj map[string]interface{}) error {
	index, modules, err := getKymaModuleIndexInSpec(moduleName, obj)
	if err != nil {
		return err
	}

	if index == -1 {
		return nil
	}

	modules[index] = modules[len(modules)-1]
	modules = modules[:len(modules)-1]
	if err := unstructured.SetNestedSlice(obj, modules, "spec", "modules"); err != nil {
		return fmt.Errorf("error setting new kyma spec modules: %w", err)
	}

	return nil
}

func addKymaModuleInSpec(moduleName string, obj map[string]interface{}) error {
	index, modules, err := getKymaModuleIndexInSpec(moduleName, obj)
	if err != nil {
		return err
	}

	if index > -1 {
		return nil
	}

	val := map[string]interface{}{
		"name": moduleName,
	}
	modules = append(modules, val)
	if err := unstructured.SetNestedSlice(obj, modules, "spec", "modules"); err != nil {
		return fmt.Errorf("error setting new kyma spec modules: %w", err)
	}

	return nil
}

func getKymaModuleStateInStatus(moduleName string, obj map[string]interface{}) (string, error) {
	modules, exists, err := unstructured.NestedSlice(obj, "status", "modules")
	if !exists {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("error reading kyma status modules: %w", err)
	}

	for _, m := range modules {
		mm, ok := m.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("kyma status module is not a map, but: %T", m)
		}
		name, exists, err := unstructured.NestedString(mm, "name")
		if !exists || err != nil {
			return "", fmt.Errorf("kyma status module does not have a name: %w", err)
		}
		if name == moduleName {
			val, exists, err := unstructured.NestedString(mm, "state")
			if !exists || err != nil {
				return "", fmt.Errorf("kyma status module does not have a state: %w", err)
			}
			return val, nil
		}
	}

	return "", nil
}
