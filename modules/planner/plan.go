package planner

import (
	"fmt"
	"kermoo/modules/fluent"
	"kermoo/modules/logger"
	"time"

	"go.uber.org/zap"
)

type Plan struct {
	// Percentage determines the percentage. Each module which reference this plan
	// may consider the percentage in their own language in the future. Currently,
	// all modules consider the percentage as the chance of failing.
	//
	// For specific and ranged declearations, it's going to use that but when an array of
	// percentages are specified, it'll act like a graph of bars and iterate over them.
	Percentage *fluent.FluentFloat `json:"percentage"`

	// Size determines the digital storage size. Currently, only memory leak module uses it.
	//
	// For specific and ranged declearations, it's going to use that but when an array of
	// sizes are specified, it'll act like a graph of bars and iterate over them.
	Size *fluent.FluentSize `json:"size"`

	// Interval decides how long each plan cycle should last. A value above one second is recommended
	// but you're free  to use any interval. Default is one second.
	Interval *fluent.FluentDuration `json:"interval"`

	// Duration defines the duration of the entire plan. Leave it empty for
	// life-long running or specify one to end the modules. Each module will end with their
	// own rules.
	// In fact, Duration/Interval determines the number of cycle, if defined. Default is empty
	// for unlimited activity.
	Duration *fluent.FluentDuration `json:"duration"`

	// Name defines the name of plan which later can be used as the reference.
	// It must be unique among all of the plans.
	Name *string `json:"name"`

	// SubPlans defines even more detailed steps of the plan. If the given values, interval
	// and the duration are not enough, leave them empty and focus on this field and use
	// sub-plans to define more complex steps. Plan will run each of the sub-plans with
	// respect to their order in a serial manner.
	SubPlans []SubPlan `json:"subPlans"`

	plannables        []*Plannable
	currentCycleValue *CycleValue
	isDedicated       bool
}

type Cycle struct {
	Value     CycleValue
	StartedAt time.Time
	TimeSpent time.Duration
}

type HookFunc func(cycle Cycle) PlanSignal

type CycleHooks struct {
	PreSleep  *HookFunc
	PostSleep *HookFunc
}

type PlanSignal uint32

const (
	PLAN_SIGNAL_CONTINUE  PlanSignal = iota
	PLAN_SIGNAL_TERMINATE PlanSignal = iota
)

func (p *Plan) ToSubPlan() SubPlan {
	return SubPlan{
		Percentage: p.Percentage,
		Size:       p.Size,
		Interval:   p.Interval,
		Duration:   p.Duration,
	}
}

func (p *Plan) Assign(plannable Plannable) {
	p.plannables = append(p.plannables, &plannable)
	plannable.AssignPlan(p)
}

func (p *Plan) MakePrivate() {
	p.isDedicated = true
}

func (p *Plan) Validate() error {
	_, err := p.GetPreparedSubPlans()

	if err != nil {
		return fmt.Errorf("unable to prepare sub-plans: %v", err)
	}

	return nil
}

func (p *Plan) GetCurrentValue() *CycleValue {
	return p.currentCycleValue
}

func (p *Plan) SetCurrentValue(cv CycleValue) {
	p.currentCycleValue = &cv
}

func (p *Plan) GetPreparedSubPlans() ([]*SubPlan, error) {
	subPlans := []*SubPlan{}

	if len(p.SubPlans) == 0 {
		subPlans = append(subPlans, &SubPlan{
			Percentage: p.Percentage,
			Size:       p.Size,
			Interval:   p.Interval,
			Duration:   p.Duration,
		})
	} else {
		for i := 0; i < len(p.SubPlans); i++ {
			subPlans = append(subPlans, &p.SubPlans[i])
		}
	}

	for _, subPlan := range subPlans {
		subPlan.SetPlan(p)
		if err := subPlan.Prepare(); err != nil {
			return nil, err
		}

	}

	return subPlans, nil
}

func (p *Plan) Start() {
	if logger.Log.Level() == zap.InfoLevel {
		logger.Log.Info("executing plan...", zap.String("name", *p.Name))
	} else {
		plannableNames := []string{}
		for _, pl := range p.plannables {
			plr := *pl
			plannableNames = append(plannableNames, plr.GetName())
		}
		logger.Log.Debug("executing plan...", zap.String("name", *p.Name), zap.Any("plan", *p), zap.Any("plannables", plannableNames))
	}

	subPlans, _ := p.GetPreparedSubPlans()

	for _, subPlan := range subPlans {
		subPlan.Execute()
	}
}

func NewPlan(p Plan) Plan {
	return p
}
