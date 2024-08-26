package internal

import (
	"github.com/onsi/gomega"
	"time"
)

const (
	DefaultEventuallyTimeout           = 5 * time.Minute
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
