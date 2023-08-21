package memory

import (
	"fmt"
	"kermoo/modules/planner"
	"kermoo/modules/values"
)

var _ planner.Plannable = &MemoryLeak{}

type MemoryLeak struct {
	planner.CanAssignPlan

	PlanRefs []string `json:"planRefs"`

	Size     *values.MultiSize `json:"size"`
	Interval *values.Duration  `json:"interval"`
	Duration *values.Duration  `json:"duration"`

	leakedData []byte
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
		size, _ := mu.GetAssignedPlans()[0].GetCurrentValue().Size.ToSize()
		mu.StartLeaking(size)
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

func (mu *MemoryLeak) StartLeaking(size values.Size) {
	mu.leakedData = make([]byte, size.ToBytes())
}

func (mu *MemoryLeak) StopLeaking() {
	mu.leakedData = make([]byte, 0)
}
