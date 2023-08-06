package planner

type Plannable interface {
	GetUid() string
	HasCustomPlan() bool
	MakeCustomPlan() *Plan
	MakeDefaultPlan() *Plan
	AssignPlan(*Plan)
	GetDesiredPlanNames() []string
	GetPlanCycleHooks() CycleHooks
}

type PlannableTrait struct {
	assignedPlans []*Plan
}

func (p *PlannableTrait) AssignPlan(plan *Plan) {
	p.assignedPlans = append(p.assignedPlans, plan)
}

func (p *PlannableTrait) GetAssignedPlans() []*Plan {
	return p.assignedPlans
}
