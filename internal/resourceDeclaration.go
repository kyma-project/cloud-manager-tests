package internal

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
)

func resourceDeclaration(ctx context.Context, tbl *godog.Table) error {
	kfrCtx := KfrFromContext(ctx)
	for i, row := range tbl.Rows {
		if len(row.Cells) < 2 {
			return fmt.Errorf("resource declaration table must have at least two cells, but row %d has %d", i, len(row.Cells))
		}
		name := row.Cells[0].Value
		if _, exists := kfrCtx.Resources[name]; exists {
			return fmt.Errorf("resource declaration name already exists: %s", name)
		}
		kind := row.Cells[1].Value
		rd := kfrCtx.Set(name, kind)
		if len(row.Cells) >= 4 {
			rd.Name = row.Cells[2].Value
			rd.Namespace = row.Cells[3].Value
		}
	}
	return nil
}
