package internal

import (
	"time"

	"github.com/onsi/gomega"
)

const (
	DefaultEventuallyTimeout           = 10 * time.Minute
	DefaultEventuallyPollingInterval   = 10 * time.Second
	DefaultConsistentlyDuration        = 20 * time.Second
	DefaultConsistentlyPollingInterval = 5 * time.Second
)

func InitGomegaDefaults() {
	gomega.Default.SetDefaultEventuallyTimeout(DefaultEventuallyTimeout)
	gomega.Default.SetDefaultEventuallyPollingInterval(DefaultEventuallyPollingInterval)
	gomega.Default.SetDefaultConsistentlyDuration(DefaultConsistentlyDuration)
	gomega.Default.SetDefaultConsistentlyPollingInterval(DefaultConsistentlyPollingInterval)
}
