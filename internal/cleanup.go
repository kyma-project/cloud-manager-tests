package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/onsi/gomega"
	"strings"
)

func cleanup(ctx context.Context, args string) error {
	var resourcesToCleanup []string
	for _, s := range strings.Split(args, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			resourcesToCleanup = append(resourcesToCleanup, s)
		}
	}

	kfrCtx := KfrFromContext(ctx)

	for _, name := range resourcesToCleanup {
		rd := kfrCtx.Get(name)
		if rd == nil {
			return fmt.Errorf("resource not declared: %s", name)
		}
		err := resourceDeleted(ctx, rd.Var)
		if errors.Is(err, NotFoundError) {
			return fmt.Errorf("resource already deleted: %s", name)
		}
		if err != nil {
			return fmt.Errorf("error deleting resource %s: %w", name, err)
		}

		var errMsg string
		gm := gomega.NewGomega(func(message string, callerSkip ...int) {
			errMsg = message
		})
		ok := gm.Eventually(func(ctx context.Context) error {
			err := rd.Reload(ctx)
			if errors.Is(err, NotFoundError) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("error reloading resource %s: %w", name, err)
			}
			return errors.New("resource still exists")
		}).
			WithArguments(ctx).
			Should(gomega.Succeed())

		if !ok || len(errMsg) > 0 {
			return fmt.Errorf("failed: %s", errMsg)
		}
	}

	return nil
}
