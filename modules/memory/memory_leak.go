package memory

import (
	"fmt"
	"kermoo/modules/planner"
	"kermoo/modules/values"
)

var _ planner.Plannable = &MemoryLeak{}

type MemoryLeak struct {
	planner.CanAssignPlan
	Plan       *planner.Plan `json:"plan"`
	PlanRefs   []string      `json:"planRefs"`
	leakedData []byte
}

func (mu *MemoryLeak) GetName() string {
	return "memory-leaker"
}

func (mu *MemoryLeak) HasInlinePlan() bool {
	return mu.Plan != nil
}

func (mu *MemoryLeak) GetDesiredPlanNames() []string {
	return mu.PlanRefs
}

func (mu *MemoryLeak) Validate() error {
	if len(mu.PlanRefs) == 0 && mu.Plan == nil {
		return fmt.Errorf("no plan or plan refs is set")
	}

	if len(mu.PlanRefs) > 1 {
		return fmt.Errorf("plan refs can not contain more than one element for memory leak")
	}

	if mu.Plan != nil {
		if err := mu.Plan.Validate(); err != nil {
			return fmt.Errorf("plan validation failed: %v", err)
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
	return mu.Plan
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
