package internal

import (
	"sort"
	"strings"
)

type cloudResourceByPriority []string

func resourcePriority(kind string) int {
	kind = strings.ToLower(kind)
	switch {
	case strings.Contains(kind, "backupschedule"):
		return 100
	case strings.Contains(kind, "restore"):
		return 200
	case strings.Contains(kind, "backup"):
		return 300
	case strings.Contains(kind, "peer"):
		return 400
	case strings.Contains(kind, "redis"):
		return 500
	case strings.Contains(kind, "nfsvolume"):
		return 800
	case kind == "iprange":
		return 9000
	case kind == "cloudresources":
		return 10000
	default:
		return 0
	}
}

func (p cloudResourceByPriority) Len() int {
	return len(p)
}

func (p cloudResourceByPriority) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p cloudResourceByPriority) Less(i, j int) bool {
	// Lower priority comes first
	return resourcePriority(p[i]) < resourcePriority(p[j])
}

func SortKindsByPriority(kinds []string) {
	sort.Stable(cloudResourceByPriority(kinds))
}
