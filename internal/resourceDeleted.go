package internal

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func resourceDeleted(ctx context.Context, ref string) error {
	kfrCtx := KfrFromContext(ctx)
	rd := kfrCtx.Get(ref)
	if rd == nil {
		return fmt.Errorf("resource %s is not declated", ref)
	}

	params := []string{
		"delete",
		"--wait=false",
	}
	if len(rd.Namespace) > 0 {
		params = append(params, "-n", rd.Namespace)
	}
	params = append(params, rd.Kind, rd.Name)

	cmd := exec.CommandContext(ctx, "kubectl", params...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		if strings.Contains(string(out), "(NotFound)") {
			rd.Value = nil
			return nil
		}
		return fmt.Errorf("error deleting resource: %w: %s", err, string(out))
	}

	return nil
}
