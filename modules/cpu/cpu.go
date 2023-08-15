package cpu

import (
	"fmt"
	"kermoo/modules/planner"
)

// Ensure that Process is implementing Plannable
var _ planner.Plannable = &Cpu{}

type Cpu struct {
	planner.PlannableTrait
	Utilize CpuUtilize `json:"utilize"`
}

func (c *Cpu) GetUid() string {
	return "cpu-manager"
}

func (c *Cpu) HasCustomPlan() bool {
	return c.Utilize.Plan != nil
}

func (c Cpu) GetDesiredPlanNames() []string {
	return c.Utilize.PlanRefs
}

func (c Cpu) Validate() error {
	if len(c.Utilize.PlanRefs) == 0 && c.Utilize.Plan == nil {
		return fmt.Errorf("no plan or plan refs is set")
	}

	if len(c.Utilize.PlanRefs) > 1 {
		return fmt.Errorf("plan refs can not contain more than one element for cpu utilization")
	}

	if c.Utilize.Plan != nil {
		if err := c.Utilize.Plan.Validate(); err != nil {
			return fmt.Errorf("plan validation failed: %v", err)
		}
	}

	return nil
}

func (c *Cpu) GetPlanCycleHooks() planner.CycleHooks {
	preSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		c.Utilize.Start(c.GetAssignedPlans()[0].GetCurrentValue().GetValue())
		return planner.PLAN_SIGNAL_CONTINUE
	})

	postSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		c.Utilize.Stop()
		return planner.PLAN_SIGNAL_CONTINUE
	})

	return planner.CycleHooks{
		PreSleep:  &preSleep,
		PostSleep: &postSleep,
	}
}

func (c *Cpu) MakeCustomPlan() *planner.Plan {
	return c.Utilize.Plan
}

func (c *Cpu) MakeDefaultPlan() *planner.Plan {
	return nil
}
