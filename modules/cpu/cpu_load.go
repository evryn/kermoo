package cpu

import (
	"context"
	"fmt"
	"kermoo/modules/planner"
	"kermoo/modules/values"
	"runtime"
	"time"
)

var _ planner.Plannable = &CpuLoader{}

type CpuLoader struct {
	planner.CanAssignPlan

	PlanRefs []string `json:"planRefs"`

	Percentage *values.MultiFloat `json:"percentage"`
	Interval   *values.Duration   `json:"interval"`
	Duration   *values.Duration   `json:"duration"`

	ctx    context.Context
	cancel context.CancelFunc
}

func (cu *CpuLoader) GetName() string {
	return "cpu-manager"
}

func (cu *CpuLoader) HasInlinePlan() bool {
	return cu.MakeInlinePlan() != nil
}

func (cu CpuLoader) GetDesiredPlanNames() []string {
	return cu.PlanRefs
}

func (cu CpuLoader) Validate() error {
	if len(cu.PlanRefs) == 0 && !cu.HasInlinePlan() {
		return fmt.Errorf("no load specifications or plan refs is set")
	}

	if len(cu.PlanRefs) > 1 {
		return fmt.Errorf("plan refs can not contain more than one element")
	}

	if cu.HasInlinePlan() {
		if err := cu.MakeInlinePlan().Validate(); err != nil {
			return fmt.Errorf("crafted plan validation failed: %v", err)
		}
	}

	return nil
}

func (cu *CpuLoader) GetPlanCycleHooks() planner.CycleHooks {
	preSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		percentage, _ := cu.GetAssignedPlans()[0].GetCurrentValue().Percentage.ToFloat()
		cu.Start(percentage)
		return planner.PLAN_SIGNAL_CONTINUE
	})

	postSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		cu.Stop()
		return planner.PLAN_SIGNAL_CONTINUE
	})

	return planner.CycleHooks{
		PreSleep:  &preSleep,
		PostSleep: &postSleep,
	}
}

func (cu *CpuLoader) MakeInlinePlan() *planner.Plan {
	if cu.Percentage == nil {
		return nil
	}

	plan := planner.NewPlan(planner.Plan{
		Percentage: cu.Percentage,
		Interval:   cu.Interval,
		Duration:   cu.Duration,
	})

	return &plan
}

func (cu *CpuLoader) MakeDefaultPlan() *planner.Plan {
	return nil
}

func (cu *CpuLoader) Start(usage float32) {
	cu.ctx, cu.cancel = context.WithCancel(context.Background())

	cu.runCpuLoad(runtime.NumCPU(), int(usage*100))
}

func (cu *CpuLoader) Stop() {
	cu.cancel()
	time.Sleep(1 * time.Millisecond)
}

func (cu *CpuLoader) GetContextAndCancel() (context.Context, context.CancelFunc) {
	return cu.ctx, cu.cancel
}

// runCpuLoad run CPU load in specify cores count and percentage
// Borrowed from: https://github.com/0Delta/gocpuload/blob/master/cpu_load.go
func (cu *CpuLoader) runCpuLoad(coresCount int, percentage int) {
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
