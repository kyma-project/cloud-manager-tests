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
	for x, row := range tbl.Rows {
		if len(row.Cells) < 1 {
			return fmt.Errorf("in row %d expected at least one cell, but got %d", x, len(row.Cells))
		}
		kind := row.Cells[0].Value

		if _, ok := kfrCtx.Values.LoadedCRDs[kind]; !ok {
			return fmt.Errorf("kind %s does not exist, but is expect it does", kind)
		}
	}

	return nil
}

func crdsDoNotExist(ctx context.Context, tbl *godog.Table) error {
	kfrCtx := KfrFromContext(ctx)
	for x, row := range tbl.Rows {
		if len(row.Cells) < 1 {
			return fmt.Errorf("in row %d expected at least one cell, but got %d", x, len(row.Cells))
		}
		kind := row.Cells[0].Value

		if _, ok := kfrCtx.Values.LoadedCRDs[kind]; ok {
			return fmt.Errorf("kind %s exist, but it's expected it does not", kind)
		}
	}

	return nil
}
