package internal

import (
	"context"
	"fmt"
	"github.com/onsi/gomega"
	"reflect"
	"strings"
)

func eventuallyValueAssertEqualsNoOptions(ctx context.Context, a string, b string) error {
	return eventuallyValueAssertEqualsWithOptions(ctx, a, b, "")
}

func eventuallyValueAssertEqualsWithOptions(ctx context.Context, a string, b string, withOpts string) error {
	timeout := DefaultEventuallyTimeout

	withOpts = strings.TrimSpace(withOpts)
	if len(withOpts) > 0 {
		opts := strings.Split(withOpts, ",")
		for _, opt := range opts {
			opt = strings.TrimSpace(opt)
			if opt == "" {
				continue
			}
			// ugly, but for now with just few timeout1-5X works, if you add more, try to find a better implementation
			switch opt {
			case "timeout2X":
				timeout = 2 * timeout
			case "timeout3X":
				timeout = 3 * timeout
			case "timeout4X":
				timeout = 4 * timeout
			case "timeout5X":
				timeout = 5 * timeout
			default:
				return fmt.Errorf("unknown option: %s", opt)
			}
		}
	}

	var errMsg string
	gm := gomega.NewGomega(func(message string, callerSkip ...int) {
		errMsg = message
	})
	ok := gm.Eventually(func(ctx context.Context, a string, b string) error {
		return valueAssertEquals(ctx, a, b)
	}, timeout).
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
