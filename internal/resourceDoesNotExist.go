package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/onsi/gomega"
)

func eventuallyResourceDoesNotExist(ctx context.Context, ref string) error {
	var errMsg string
	gm := gomega.NewGomega(func(message string, callerSkip ...int) {
		errMsg = message
	})
	ok := gm.Eventually(func(ctx context.Context, ref string) error {
		return resourceDoesNotExist(ctx, ref)
	}).
		WithArguments(ctx, ref).
		Should(gomega.Succeed())
	if !ok || len(errMsg) > 0 {
		return fmt.Errorf("failed: %s", errMsg)
	}
	return nil
}

func resourceDoesNotExist(ctx context.Context, ref string) error {
	kfrCtx := KfrFromContext(ctx)
	rd := kfrCtx.Get(ref)
	if rd == nil {
		return fmt.Errorf("resource not declared: %s", ref)
	}
	err := rd.Reload(ctx)
	if errors.Is(err, NotFoundError) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error loading resource %s: %w", ref, err)
	}
	return fmt.Errorf("resource exists: %s", ref)
}
