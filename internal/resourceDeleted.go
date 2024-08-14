package internal

import (
	"context"
	"errors"
	"fmt"
)

func resourceDeleted(ctx context.Context, ref string) error {
	kfrCtx := KfrFromContext(ctx)
	rd := kfrCtx.Get(ref)
	if rd == nil {
		return fmt.Errorf("resource %s is not declated", ref)
	}

	err := kfrCtx.K8S.Delete(ctx, rd.Kind, rd.Name, rd.Namespace, false)
	if errors.Is(err, NotFoundError) {
		rd.Value = nil
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}
