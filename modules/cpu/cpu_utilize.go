package cpu

import (
	"context"
	"kermoo/modules/logger"
	"kermoo/modules/planner"
	"kermoo/modules/utils"
	"runtime"
	"time"

	"go.uber.org/zap"
)

type CpuUtilize struct {
	Plan     *planner.Plan `json:"plan"`
	PlanRefs []string      `json:"planRefs"`

	ctx          context.Context
	cancel       context.CancelFunc
	currentUsage float32
}

func (cu *CpuUtilize) Start(usage float32) {
	cu.ctx, cu.cancel = context.WithCancel(context.Background())

	go cu.updateCpuUsage(cu.ctx)

	coreCount := runtime.NumCPU()

	go func() {
		for {
			logger.Log.Debug("usage", zap.Float32("usage", cu.currentUsage))
			time.Sleep(5 * time.Millisecond)
		}
	}()

	for i := 0; i < coreCount; i++ {
		time.Sleep(1 * time.Millisecond)
		go cu.utilize(cu.ctx, usage)
	}
}

func (cu *CpuUtilize) Stop() {
	cu.cancel()
	time.Sleep(1 * time.Millisecond)
}

func (cu *CpuUtilize) utilize(ctx context.Context, targetUsage float32) {
	logger.Log.Debug("started utilization goroutine")

	for {
		if cu.currentUsage < targetUsage {
			for i := 0; i < 100000; i++ {
			}
		} else {
			time.Sleep(1 * time.Millisecond)
		}

		select {
		case <-ctx.Done():
			logger.Log.Debug("ending utilization goroutine")
			return
		default:
		}
	}
}

func (cu *CpuUtilize) updateCpuUsage(ctx context.Context) {
	var err error

	for {
		cu.currentUsage, err = utils.GetCpuUsage(0)

		if err != nil {
			logger.Log.Fatal("error getting cpu usage", zap.Error(err))
		}

		time.Sleep(10 * time.Millisecond)

		select {
		case <-ctx.Done():
			logger.Log.Debug("ending cpu usage updater")
			return
		default:
		}
	}
}
