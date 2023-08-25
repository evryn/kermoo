package memory

import (
	"fmt"
	"kermoo/modules/fluent"
	"kermoo/modules/planner"
)

var _ planner.Plannable = &MemoryLeak{}

type MemoryLeak struct {
	planner.CanAssignPlan

	// PlanRefs is an optional list of plan names. It can used to avoid redundant
	// re-declearing of plans in large-scale configurations.
	// PlanRefs overrides Size, Interval and Duration fields are overrided in favor
	// of the one defined in the referenced plan.
	PlanRefs []string `json:"planRefs"`

	// Size determines the size of the memory leak (memory consumption). This memory will be
	// used in addition to the amount used by the Kermoo application itself. So the actual
	// total memory usage is not guaranteed to be accurate.
	//
	// For specific and ranged declearations, it's going to use that but when an array of
	// sizes are specified, it'll act like a graph of bars and iterate over them.
	Size *fluent.FluentSize `json:"size"`

	// Interval decides how long each leak cycle should last. A value above one second is recommended
	// but you're free  to use any interval. Default is one second.
	Interval *fluent.FluentDuration `json:"interval"`

	// Duration defines the duration of the entire memory leak module. Leave it empty for
	// life-long running or specify one to end the module completely after that and won't
	// consume the specified memory.
	// In fact, Duration/Interval determines the number of cycle, if defined. Default is empty
	// for unlimited activity.
	Duration *fluent.FluentDuration `json:"duration"`

	leakedData []byte
}

func (mu *MemoryLeak) GetLeakedData() []byte {
	return mu.leakedData
}

func (mu *MemoryLeak) GetName() string {
	return "memory-leaker"
}

func (mu *MemoryLeak) HasInlinePlan() bool {
	return mu.MakeInlinePlan() != nil
}

func (mu *MemoryLeak) GetDesiredPlanNames() []string {
	return mu.PlanRefs
}

func (mu *MemoryLeak) Validate() error {
	if len(mu.PlanRefs) == 0 && !mu.HasInlinePlan() {
		return fmt.Errorf("no leak specifications or plan refs is set")
	}

	if len(mu.PlanRefs) > 1 {
		return fmt.Errorf("plan refs can not contain more than one element")
	}

	if mu.HasInlinePlan() {
		if err := mu.MakeInlinePlan().Validate(); err != nil {
			return fmt.Errorf("crafted plan validation failed: %v", err)
		}
	}

	return nil
}

func (mu *MemoryLeak) GetPlanCycleHooks() planner.CycleHooks {
	preSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		mu.StartLeaking(
			mu.GetAssignedPlans()[0].GetCurrentValue().Size,
		)
		return planner.PLAN_SIGNAL_CONTINUE
	})

	postSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		mu.StopLeaking()
		return planner.PLAN_SIGNAL_CONTINUE
	})

	return planner.CycleHooks{
		PreSleep:  &preSleep,
		PostSleep: &postSleep,
	}
}

func (mu *MemoryLeak) MakeInlinePlan() *planner.Plan {
	if mu.Size == nil {
		return nil
	}

	plan := planner.NewPlan(planner.Plan{
		Size:     mu.Size,
		Interval: mu.Interval,
		Duration: mu.Duration,
	})

	return &plan
}

func (mu *MemoryLeak) MakeDefaultPlan() *planner.Plan {
	return nil
}

func (mu *MemoryLeak) StartLeaking(size int64) {
	mu.leakedData = make([]byte, size)
}

func (mu *MemoryLeak) StopLeaking() {
	mu.leakedData = make([]byte, 0)
}
