package main

import (
	"flag"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/kyma-project/cloud-manager-tests/internal"
	"github.com/onsi/gomega"
	"os"
	"time"
)

var opts = godog.Options{
	Output:      colors.Colored(os.Stdout),
	Concurrency: 1,
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func main() {
	gomega.Default.SetDefaultEventuallyTimeout(5 * time.Minute)
	gomega.Default.SetDefaultEventuallyPollingInterval(10 * time.Second)
	gomega.Default.SetDefaultConsistentlyDuration(20 * time.Second)
	gomega.Default.SetDefaultConsistentlyPollingInterval(5 * time.Second)

	flag.Parse()
	o := opts
	status := godog.TestSuite{
		Name:                 "k8sFeatureRunner",
		Options:              &o,
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
	}.Run()

	os.Exit(status)
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() { fmt.Println("Get the party started!") })
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	internal.Register(ctx)
}
