package internal

import (
	"context"
	"fmt"
	"github.com/onsi/gomega"
	"strings"
)

func noCloudResources(ctx context.Context) error {
	kfrCtx := KfrFromContext(ctx)

	// First deleting pods that are created for E2E tests so they don't block PVC, PV and NfsVolume deletion
	podsForE2ETests, err := kfrCtx.K8S.List(ctx, "pods", "-A", "context=cloud-manager-tests")
	if err != nil {
		return err
	}
	for _, pod := range podsForE2ETests {
		if err := kfrCtx.K8S.Delete(ctx, "pods", pod.Name, pod.Namespace, false); err != nil {
			return fmt.Errorf("error deleting pod %s/%s: %w", pod.Namespace, pod.Name, err)
		}
	}

	kinds, err := kfrCtx.K8S.CloudResourceKinds(ctx)
	if err != nil {
		return err
	}

	SortKindsByPriority(kinds)

	for _, kind := range kinds {
		if strings.ToLower(kind) == "cloudresources" {
			continue
		}

		if err := kfrCtx.K8S.Delete(ctx, kind, "", "-A", true); err != nil {
			return fmt.Errorf("error deleting %s: %w", kind, err)
		}
	}

	var errMsg string
	gm := gomega.NewGomega(func(message string, callerSkip ...int) {
		errMsg = message
	})
	ok := gm.Eventually(func(ctx context.Context) error {
		for _, kind := range kinds {
			if strings.ToLower(kind) == "cloudresources" {
				continue
			}
			arr, err := kfrCtx.K8S.List(ctx, kind, "-A", "")
			if err != nil {
				return err
			}
			if len(arr) > 0 {
				return fmt.Errorf("resources of kind %s still exist", kind)
			}
		}
		return nil
	}, 5*DefaultEventuallyTimeout).
		WithArguments(ctx).
		Should(gomega.Succeed())

	if !ok || len(errMsg) > 0 {
		return fmt.Errorf("failed removing all cloud resources: %s", errMsg)
	}
	return nil
}
