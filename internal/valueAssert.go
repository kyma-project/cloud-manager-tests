package internal

import (
	"context"
	"fmt"
	"github.com/onsi/gomega"
	"reflect"
)

func eventuallyValueAssertEquals(ctx context.Context, a string, b string) error {
	var errMsg string
	gm := gomega.NewGomega(func(message string, callerSkip ...int) {
		errMsg = message
	})
	ok := gm.Eventually(func(ctx context.Context, a string, b string) error {
		return valueAssertEquals(ctx, a, b)
	}).
		WithArguments(ctx, a, b).
		Should(gomega.Succeed())
	if !ok || len(errMsg) > 0 {
		return fmt.Errorf("failed: %s", errMsg)
	}
	return nil
}

func valueAssertEquals(ctx context.Context, a string, b string) error {
	kfrCtx := KfrFromContext(ctx)
	aa, err := kfrCtx.Eval(ctx, a)
	if err != nil {
		return fmt.Errorf("error evaluating value %s: %w", a, err)
	}
	bb, err := kfrCtx.Eval(ctx, b)
	if err != nil {
		return fmt.Errorf("error evaluating value %s: %w", b, err)
	}

	if !reflect.DeepEqual(aa, bb) {
		return fmt.Errorf("value %s (%v) and value %s (%v) are not equal", a, aa, b, bb)
	}

	return nil
}
