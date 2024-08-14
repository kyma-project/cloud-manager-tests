package internal

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"os"
	"os/exec"
	"strings"
)

type K8sClient interface {
	Apply(ctx context.Context, txt string) error
	Get(ctx context.Context, kind, name, namespace string) (map[string]interface{}, error)
	Delete(ctx context.Context, kind, name, namespace string, all bool) error
	Logs(ctx context.Context, name, namespace string) (string, error)
	CloudResourceKinds(ctx context.Context) ([]string, error)
	List(ctx context.Context, kind, namespace, labelSelector string) ([]*ResourceDefn, error)
}

type k8sClient struct{}

func (c *k8sClient) Apply(ctx context.Context, txt string) error {
	f, err := os.CreateTemp("", "kfr-*.yaml")
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}
	defer os.Remove(f.Name())

	if err := os.WriteFile(f.Name(), []byte(txt), 0644); err != nil {
		return fmt.Errorf("error writing temp file: %w", err)
	}

	params := []string{
		"apply",
		"-f", f.Name(),
	}
	cmd := exec.CommandContext(ctx, "kubectl", params...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("error applying resource: %w: %s", err, string(out))
	}

	return nil
}

func (c *k8sClient) Get(ctx context.Context, kind, name, namespace string) (map[string]interface{}, error) {
	params := []string{
		"get",
	}
	if len(namespace) > 0 {
		params = append(params, "--namespace", namespace)
	}
	params = append(
		params,
		kind,
		name,
		"-o",
		"json",
	)
	cmd := exec.CommandContext(ctx, "kubectl", params...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "(NotFound)") {
			return nil, NotFoundError
		}
		return nil, fmt.Errorf("error getting resource: %w: %s", err, string(out))
	}

	val := map[string]interface{}{}
	err = json.Unmarshal(out, &val)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling yaml: %w", err)
	}

	return val, nil
}

func (c *k8sClient) Delete(ctx context.Context, kind, name, namespace string, all bool) error {
	params := []string{
		"delete",
		"--wait=false",
	}
	if namespace == "-A" {
		params = append(params, "-A")
	} else if len(namespace) > 0 {
		params = append(params, "-n", namespace)
	}
	if all {
		params = append(params, "--all")
	}

	params = append(params, kind)

	if len(name) > 0 {
		params = append(params, name)
	}

	cmd := exec.CommandContext(ctx, "kubectl", params...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		if strings.Contains(string(out), "(NotFound)") {
			return NotFoundError
		}
		return fmt.Errorf("error deleting resource: %w: %s", err, string(out))
	}

	return nil
}

func (c *k8sClient) Logs(ctx context.Context, name, namespace string) (string, error) {
	params := []string{
		"logs",
	}
	if len(namespace) > 0 {
		params = append(params, "--namespace", namespace)
	}
	params = append(
		params,
		name,
	)
	cmd := exec.CommandContext(ctx, "kubectl", params...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "(NotFound)") {
			return "", NotFoundError
		}
		return "", fmt.Errorf("error getting logs: %w: %s", err, string(out))
	}

	return string(out), nil
}

func (c *k8sClient) CloudResourceKinds(ctx context.Context) ([]string, error) {
	params := []string{
		"api-resources",
		"--api-group=cloud-resources.kyma-project.io",
	}
	cmd := exec.CommandContext(ctx, "kubectl", params...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error getting cloud resource kinds: %w: %s", err, string(out))
	}
	lines := strings.Split(string(out), "\n")
	var result []string
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		chunks := strings.Fields(line)
		if len(chunks) < 4 {
			return result, fmt.Errorf("error parsing cloud resource kinds: line %q: invalid format", line)
		}
		kind := chunks[len(chunks)-1]
		result = append(result, kind)
	}
	return result, nil
}

func (c *k8sClient) List(ctx context.Context, kind, namespace, labelSelector string) ([]*ResourceDefn, error) {
	params := []string{
		"get",
		kind,
	}
	if namespace == "-A" {
		params = append(params, "-A")
	} else if len(namespace) > 0 {
		params = append(params, "-n", namespace)
	}
	if labelSelector != "" {
		params = append(params, "-l", "'"+labelSelector+"'")
	}
	params = append(params, "-o", "json")

	cmd := exec.CommandContext(ctx, "kubectl", params...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error listing resources: %w: %s", err, string(out))
	}

	list := map[string]interface{}{}
	err = json.Unmarshal(out, &list)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling yaml: %w", err)
	}

	arr, found, err := unstructured.NestedSlice(list, "items")
	if err != nil {
		return nil, fmt.Errorf("error getting resource list: %w", err)
	}
	if !found {
		return nil, nil
	}
	var result []*ResourceDefn
	for x, item := range arr {
		obj, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("item %d is not a map but %T", x, item)
		}
		rd := &ResourceDefn{
			Kind:  obj["kind"].(string),
			Value: obj,
		}
		if err := rd.ExtractNames(); err != nil {
			return nil, fmt.Errorf("error extracting item %d names: %w", x, err)
		}
		result = append(result, rd)
	}

	return result, nil
}
