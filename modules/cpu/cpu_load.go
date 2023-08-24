package cpu

import (
	"context"
	"fmt"
	"kermoo/modules/fluent"
	"kermoo/modules/planner"
	"runtime"
	"time"
)

var _ planner.Plannable = &CpuLoader{}

type CpuLoader struct {
	planner.CanAssignPlan

	// PlanRefs is an optional list of plan names. It can used to avoid redundant
	// re-declearing of plans in large-scale configurations.
	// PlanRefs overrides Percentage, Interval and Duration fields are overrided in favor
	// of the one defined in the referenced plan.
	PlanRefs []string `json:"planRefs"`

	// Percentage determines CPU load in percentage. 0 means no additional load and 100 means
	// to use all cores as much as possible. The percentage is not guaranteed to be accurate.
	//
	// For specific and ranged declearations, it's going to use that but when an array of
	// percentages are specified, it'll act like a graph of bars and iterate over them.
	Percentage *fluent.FluentFloat `json:"percentage"`

	// Interval decides how long each load cycle should last. A value above one second is recommended
	// but you're free  to use any interval. Default is one second.
	Interval *fluent.FluentDuration `json:"interval"`

	// Duration defines the duration of the entire CPU load module. Leave it empty for
	// life-long running or specify one to end the module completely after that and won't do
	// any additional load.
	// In fact, Duration/Interval determines the number of cycle, if defined. Default is empty
	// for unlimited activity.
	Duration *fluent.FluentDuration `json:"duration"`

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
		cu.Start(
			cu.GetAssignedPlans()[0].GetCurrentValue().Percentage,
		)
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

func (cu *CpuLoader) Start(usagePercentage float64) {
	cu.ctx, cu.cancel = context.WithCancel(context.Background())

	cu.runCpuLoad(runtime.NumCPU(), int(usagePercentage))
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
