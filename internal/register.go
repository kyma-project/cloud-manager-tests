package internal

import (
	"context"
	"github.com/cucumber/godog"
)

func Register(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		return KfrToContext(ctx, &KfrContext{
			Resources: map[string]*ResourceDefn{},
		}), nil
	})
	ctx.Step("^resource declaration:?$", resourceDeclaration)
	ctx.Step("^resource ([^ ]+)? ?is applied:?$", resourceApplied)
	ctx.Step("^value (.*) equals (.*)$", valueAssertEquals)
	ctx.Step("^eventually value (.*) equals (.*)$", eventuallyValueAssertEquals)
	ctx.Step(`^resource (.*) is deleted`, resourceDeleted)
	ctx.Step(`^cleanup (.*)$`, cleanup)
	ctx.Step(`^resource (.*) does not exist$`, resourceDoesNotExist)
	ctx.Step(`^eventually resource (.*) does not exist$`, eventuallyResourceDoesNotExist)
	ctx.Step(`^(.*)$`, script)
}
