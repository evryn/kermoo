package cpu

import (
	"fmt"
	"kermoo/modules/planner"
)

// Ensure that Process is implementing Plannable
var _ planner.Plannable = &Cpu{}

type Cpu struct {
	planner.PlannableTrait
	Load CpuLoader `json:"load"`
}

func (c *Cpu) GetName() string {
	return "cpu-manager"
}

func (c *Cpu) HasInlinePlan() bool {
	return c.Load.Plan != nil
}

func (c Cpu) GetDesiredPlanNames() []string {
	return c.Load.PlanRefs
}

func (c Cpu) Validate() error {
	if len(c.Load.PlanRefs) == 0 && c.Load.Plan == nil {
		return fmt.Errorf("no plan or plan refs is set")
	}

	if len(c.Load.PlanRefs) > 1 {
		return fmt.Errorf("plan refs can not contain more than one element for cpu utilization")
	}

	if c.Load.Plan != nil {
		if err := c.Load.Plan.Validate(); err != nil {
			return fmt.Errorf("plan validation failed: %v", err)
		}
	}

	return nil
}

func (c *Cpu) GetPlanCycleHooks() planner.CycleHooks {
	preSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		c.Load.Start(c.GetAssignedPlans()[0].GetCurrentValue().GetValue())
		return planner.PLAN_SIGNAL_CONTINUE
	})

	postSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		c.Load.Stop()
		return planner.PLAN_SIGNAL_CONTINUE
	})

	return planner.CycleHooks{
		PreSleep:  &preSleep,
		PostSleep: &postSleep,
	}
}

func (c *Cpu) MakeInlinePlan() *planner.Plan {
	return c.Load.Plan
}

func (c *Cpu) MakeDefaultPlan() *planner.Plan {
	return nil
}
