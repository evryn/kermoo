package cpu

import (
	"context"
	"kermoo/modules/planner"
	"runtime"
	"time"
)

type CpuUtilize struct {
	Plan     *planner.Plan `json:"plan"`
	PlanRefs []string      `json:"planRefs"`

	ctx    context.Context
	cancel context.CancelFunc
}

func (cu *CpuUtilize) Start(usage float32) {
	cu.ctx, cu.cancel = context.WithCancel(context.Background())

	cu.runCpuLoad(runtime.NumCPU(), int(usage*100))
}

func (cu *CpuUtilize) Stop() {
	cu.cancel()
	time.Sleep(1 * time.Millisecond)
}

// runCpuLoad run CPU load in specify cores count and percentage
// Borrowed from: https://github.com/0Delta/gocpuload/blob/master/cpu_load.go
func (cu *CpuUtilize) runCpuLoad(coresCount int, percentage int) {
	runtime.GOMAXPROCS(coresCount)

	// 1 unit = 100 ms may be the best
	unitHundresOfMicrosecond := 1000
	runMicrosecond := unitHundresOfMicrosecond * percentage
	sleepMicrosecond := unitHundresOfMicrosecond*100 - runMicrosecond

	for i := 0; i < coresCount; i++ {
		go func(ctx context.Context) {
			runtime.LockOSThread()
			// endless loop
			for {
				select {
				case <-ctx.Done():
					return
				default:
					begin := time.Now()
					for {
						// run 100%
						if time.Since(begin) > time.Duration(runMicrosecond)*time.Microsecond {
							break
						}
					}
					// sleep
					time.Sleep(time.Duration(sleepMicrosecond) * time.Microsecond)
				}
			}
		}(cu.ctx)
	}
}
