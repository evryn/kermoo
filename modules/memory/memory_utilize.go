package memory

import (
	"kermoo/modules/planner"
)

type MemoryUtilize struct {
	Plan       *planner.Plan `json:"plan"`
	PlanRefs   []string      `json:"planRefs"`
	leakedData []byte
}

func (mu *MemoryUtilize) Start(bytes int64) {
	mu.leakedData
}

func (mu *MemoryUtilize) Stop() {
	mu.leakedData = []byte("")
}
