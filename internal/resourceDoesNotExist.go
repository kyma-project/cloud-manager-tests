package internal

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/onsi/gomega"
)

func eventuallyResourceDoesNotExist(ctx context.Context, ref string) error {
	return eventuallyResourceDoesNotExistWithOptions(ctx, ref, "")
}

func eventuallyResourceDoesNotExistWithOptions(ctx context.Context, ref string, withOpts string) error {
	timeout := DefaultEventuallyTimeout

	withOpts = strings.TrimSpace(withOpts)
	if len(withOpts) > 0 {
		opts := strings.Split(withOpts, ",")
		for _, opt := range opts {
			opt = strings.TrimSpace(opt)
			if opt == "" {
				continue
			}
			// ugly, but for now with just few timeout1-5X works, if you add more, try to find a better implementation.
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
	ok := gm.Eventually(func(ctx context.Context, ref string) error {
		return resourceDoesNotExist(ctx, ref)
	}, timeout).
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
