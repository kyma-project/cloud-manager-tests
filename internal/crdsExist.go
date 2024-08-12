package internal

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
)

func crdsAreLoaded(ctx context.Context) error {
	kfrCtx := KfrFromContext(ctx)
	kinds, err := kfrCtx.K8S.CloudResourceKinds(ctx)
	if err != nil {
		return err
	}
	kfrCtx.Values.LoadedCRDs = make(map[string]struct{}, len(kinds))
	for _, kind := range kinds {
		kfrCtx.Values.LoadedCRDs[kind] = struct{}{}
	}

	return nil
}

func crdsExist(ctx context.Context, tbl *godog.Table) error {
	kfrCtx := KfrFromContext(ctx)
	kinds, err := kfrCtx.K8S.CloudResourceKinds(ctx)
	if err != nil {
		return err
	}
	kindIndex := make(map[string]struct{}, len(kinds))
	for _, kind := range kinds {
		kindIndex[kind] = struct{}{}
	}

	for x, row := range tbl.Rows {
		if len(row.Cells) < 1 {
			return fmt.Errorf("in row %d expected at least one cell, but got %d", x, len(row.Cells))
		}
		kind := row.Cells[0].Value
		shouldCheckExistance, err := func() (bool, error) {
			if len(row.Cells) > 1 {
				cond := row.Cells[1].Value
				if cond == "" || cond == "true" {
					return true, nil
				}
				r, err := kfrCtx.Eval(ctx, cond)
				if err != nil {
					return false, fmt.Errorf("error evaluating expression %q in row %d: %w", cond, x, err)
				}
				b, ok := r.(bool)
				if !ok {
					return false, fmt.Errorf("expected bool but expression %q in row %d evaluates to %T", cond, x, r)
				}
				if !b {
					return false, nil
				}
			}
			return true, nil
		}()
		if err != nil {
			return err
		}
		if !shouldCheckExistance {
			continue
		}

		if _, ok := kindIndex[kind]; !ok {
			return fmt.Errorf("kind %s does not exist", kind)
		}
	}

	return nil
}
