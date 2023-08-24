package planner

import (
	"fmt"
	"kermoo/modules/fluent"
	"kermoo/modules/logger"
	"kermoo/modules/values"
	"time"

	"go.uber.org/zap"
)

type Plan struct {
	Percentage        *fluent.FluentFloat `json:"percentage"`
	Size              *fluent.FluentSize  `json:"size"`
	Interval          *values.Duration    `json:"interval"`
	Duration          *values.Duration    `json:"duration"`
	Name              *string             `json:"name"`
	SubPlans          []SubPlan           `json:"subPlans"`
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
