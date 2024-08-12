package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortKindsByPriority(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		kinds := []string{"CloudResources", "GcpNfsVolumeBackup", "AwsNfsVolume", "GcpNfsVolumeRestore", "IpRange", "NfsBackupSchedule", "AwsVpcPeering", "GcpRedisInstance"}
		SortKindsByPriority(kinds)
		assert.Equal(t, []string{"NfsBackupSchedule", "GcpNfsVolumeRestore", "GcpNfsVolumeBackup", "AwsVpcPeering", "GcpRedisInstance", "AwsNfsVolume", "IpRange", "CloudResources"}, kinds)
	})
}
