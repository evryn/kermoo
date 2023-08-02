package planner

type Plannable interface {
	GetUid() string
	HasCustomPlan() bool
	MakeCustomPlan() *Plan
	AssignPlans([]*Plan)
	GetDesiredPlanNames() []string
}

type PlannableTrait struct {
	assignedPlans []*Plan
}

func (p *PlannableTrait) AssignPlans(ap []*Plan) {
	p.assignedPlans = ap
}
