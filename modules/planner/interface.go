package planner

type Plannable interface {
	GetName() string
	HasInlinePlan() bool
	MakeInlinePlan() *Plan
	MakeDefaultPlan() *Plan
	AssignPlan(*Plan)
	GetDesiredPlanNames() []string
	GetPlanCycleHooks() CycleHooks
}

type CanAssignPlan struct {
	assignedPlans []*Plan
}

func (p *CanAssignPlan) AssignPlan(plan *Plan) {
	p.assignedPlans = append(p.assignedPlans, plan)
}

func (p *CanAssignPlan) GetAssignedPlans() []*Plan {
	return p.assignedPlans
}
