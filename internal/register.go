package internal

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
)

func Register(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		for _, tag := range sc.Tags {
			if tag.Name == "@skip" {
				fmt.Printf("Skipping scenario: %s %s\n", sc.Name, sc.Uri)
				return ctx, godog.ErrSkip
			}
		}
		ctx = KfrToContext(ctx, &KfrContext{
			Resources: map[string]*ResourceDefn{},
			K8S:       &k8sClient{},
			Values: KfrValues{
				Env: os.Getenv("ENV"),
			},
		})
		kfrCtx := KfrFromContext(ctx)
		cm, err := kfrCtx.K8S.Get(ctx, "ConfigMap", "shoot-info", "kube-system")
		if err != nil {
			return ctx, fmt.Errorf("error loading shoot-info configmap: %w", err)
		}

		s, f, err := unstructured.NestedString(cm, "data", "shootName")
		if !f || err != nil {
			return ctx, fmt.Errorf("error getting shootName: %w", err)
		}
		kfrCtx.Values.Shoot = s

		s, f, err = unstructured.NestedString(cm, "data", "provider")
		if !f || err != nil {
			return ctx, fmt.Errorf("error getting provider: %w", err)
		}
		kfrCtx.Values.Provider = s

		return ctx, nil
	})
	ctx.Step("^resource declaration:?$", resourceDeclaration)
	ctx.Step("^resource ([^ ]+)? ?is applied:?$", resourceApplied)
	ctx.Step("^value (.*) equals (.*)$", valueAssertEquals)
	ctx.Step("^eventually value (.*) equals (.*)$", eventuallyValueAssertEquals)
	ctx.Step(`^resource (.*) is deleted`, resourceDeleted)
	ctx.Step(`^cleanup (.*)$`, cleanup)
	ctx.Step(`^resource (.*) does not exist$`, resourceDoesNotExist)
	ctx.Step(`^eventually resource (.*) does not exist$`, eventuallyResourceDoesNotExist)
	ctx.Step(`^there are no cloud resources$`, noCloudResources)
	ctx.Step(`^module ([^ ]+)? ?is removed$`, moduleRemoved)
	ctx.Step(`^module ([^ ]+)? ?is added`, moduleAdded)
	ctx.Step(`^CRDs are loaded$`, crdsAreLoaded)
	ctx.Step(`^CRDs exist:$`, crdsExist)
}
