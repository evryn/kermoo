package planner

type Plannable interface {
	GetUid() string
	HasCustomPlan() bool
	MakeCustomPlan() *Plan
	AssignPlan(*Plan)
	GetDesiredPlanNames() []string
	GetPlanCallbacks() Callbacks
}

type PlannableTrait struct {
	assignedPlans []*Plan
}

func (p *PlannableTrait) AssignPlan(plan *Plan) {
	p.assignedPlans = append(p.assignedPlans, plan)
}
