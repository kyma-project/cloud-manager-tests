package main

import (
	"flag"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/kyma-project/cloud-manager-tests/internal"
	"os"
)

var opts = godog.Options{
	Output:        colors.Colored(os.Stdout),
	Concurrency:   1,
	StopOnFailure: true,
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func main() {
	internal.InitGomegaDefaults()

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
