package memory

// import (
// 	"fmt"
// 	"kermoo/modules/planner"
// )

// // Ensure that Process is implementing Plannable
// var _ planner.Plannable = &Memory{}

// type Memory struct {
// 	planner.PlannableTrait
// 	Leak MemoryUtilize `json:"utilize"`
// }

// func (m *Memory) GetUid() string {
// 	return "memory-manager"
// }

// func (m *Memory) HasCustomPlan() bool {
// 	return m.Utilize.Plan != nil
// }

// func (m Memory) GetDesiredPlanNames() []string {
// 	return c.Utilize.PlanRefs
// }

// func (m Memory) Validate() error {
// 	if len(c.Utilize.PlanRefs) == 0 && c.Utilize.Plan == nil {
// 		return fmt.Errorf("no plan or plan refs is set")
// 	}

// 	if len(c.Utilize.PlanRefs) > 1 {
// 		return fmt.Errorf("plan refs can not contain more than one element for memory utilization")
// 	}

// 	if c.Utilize.Plan != nil {
// 		if err := c.Utilize.Plan.Validate(); err != nil {
// 			return fmt.Errorf("plan validation failed: %v", err)
// 		}
// 	}

// 	return nil
// }

// func (m *Memory) GetPlanCycleHooks() planner.CycleHooks {
// 	preSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
// 		c.Utilize.Start(c.GetAssignedPlans()[0].GetCurrentValue().GetValue())
// 		return planner.PLAN_SIGNAL_CONTINUE
// 	})

// 	postSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
// 		c.Utilize.Stop()
// 		return planner.PLAN_SIGNAL_CONTINUE
// 	})

// 	return planner.CycleHooks{
// 		PreSleep:  &preSleep,
// 		PostSleep: &postSleep,
// 	}
// }

// func (m *Memory) MakeCustomPlan() *planner.Plan {
// 	return c.Utilize.Plan
// }

// func (m *Memory) MakeDefaultPlan() *planner.Plan {
// 	return nil
// }
