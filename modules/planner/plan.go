package planner

import (
	"fmt"
	"kermoo/modules/logger"
	"kermoo/modules/utils"
	"kermoo/modules/values"
	"time"

	"go.uber.org/zap"
)

type Plan struct {
	Value                  *values.MultiFloat `json:"value"`
	Size                   *values.MultiSize  `json:"size"`
	Interval               *values.Duration   `json:"interval"`
	Duration               *values.Duration   `json:"duration"`
	Name                   *string            `json:"name"`
	SubPlans               []SubPlan          `json:"subPlans"`
	plannables             []*Plannable
	currentExecutableValue *ExecutableValue
	currentStateByChance   bool
	isDedicated            bool
	executablePlans        []*ExecutablePlan
}

type Cycle struct {
	ExecutablePlan  *ExecutablePlan
	ExecutableValue *ExecutableValue
	StartedAt       time.Time
	TimeSpent       time.Duration
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
		Value:    p.Value,
		Size:     p.Size,
		Interval: p.Interval,
		Duration: p.Duration,
	}
}

func (p *Plan) Assign(plannable Plannable) {
	p.plannables = append(p.plannables, &plannable)
	plannable.AssignPlan(p)
}

func (p *Plan) MakePrivate() {
	p.isDedicated = true
}

func (p *Plan) GetExecutablePlans() ([]*ExecutablePlan, error) {
	if len(p.SubPlans) == 0 {
		subPlan := p.ToSubPlan()
		executablePlan, err := subPlan.ToExecutablePlan()

		if err != nil {
			return nil, fmt.Errorf("failed to convert plan itself to executable plan: %v", err)
		}

		return []*ExecutablePlan{executablePlan}, nil
	}

	// If the plans has SubPlans, generate coresponding executable SubPlans
	var ep = []*ExecutablePlan{}

	for _, sp := range p.SubPlans {
		executablePlan, err := sp.ToExecutablePlan()

		if err != nil {
			return nil, fmt.Errorf("failed to convert subplan to executable plan: %v", err)
		}

		ep = append(ep, executablePlan)
	}

	return ep, nil
}

func (p *Plan) Validate() error {
	_, err := p.GetExecutablePlans()

	if err != nil {
		return fmt.Errorf("unable to generate executable plans: %v", err)
	}

	return nil
}

func (p *Plan) GetCurrentValue() *ExecutableValue {
	return p.currentExecutableValue
}

func (p *Plan) GetCurrentStateByChance() bool {
	return p.currentStateByChance
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

	if len(p.executablePlans) == 0 {
		executablePlans, _ := p.GetExecutablePlans()

		p.executablePlans = executablePlans
	}

	for _, ep := range p.executablePlans {
		for ep.IsForever || ep.CurrentTries <= ep.TotalTries {
			for _, ev := range ep.Values {
				p.currentExecutableValue = ev
				p.currentStateByChance = utils.IsSuccessByChance(ev.GetValue())
				startedAt := time.Now()

				logger.Log.Info("executing pre-sleep hooks...", zap.String("plan", *p.Name), zap.Float32("value", ev.GetValue()), zap.Bool("success_by_chance", p.GetCurrentStateByChance()))

				if !ep.IsForever {
					ep.CurrentTries++

					if ep.CurrentTries > ep.TotalTries {
						break
					}
				}

				for _, pl := range p.plannables {
					plannable := *pl
					hook := plannable.GetPlanCycleHooks().PreSleep
					if hook != nil {
						executable := *hook
						value := executable(Cycle{
							ExecutablePlan:  ep,
							ExecutableValue: ev,
							StartedAt:       startedAt,
							TimeSpent:       time.Since(startedAt),
						})

						if value == PLAN_SIGNAL_TERMINATE {
							logger.Log.Info("terminating plan by signal", zap.String("plan", *p.Name), zap.String("cause", plannable.GetName()))
							return
						}
					}
				}

				time.Sleep(ep.Interval)

				logger.Log.Info("executing post-sleep hooks...", zap.String("plan", *p.Name))

				for _, pl := range p.plannables {
					plannable := *pl
					hook := plannable.GetPlanCycleHooks().PostSleep
					if hook != nil {
						executable := *hook
						value := executable(Cycle{
							ExecutablePlan:  ep,
							ExecutableValue: ev,
							StartedAt:       startedAt,
							TimeSpent:       time.Since(startedAt),
						})

						if value == PLAN_SIGNAL_TERMINATE {
							logger.Log.Info("terminating plan by signal", zap.String("plan", *p.Name), zap.String("cause", plannable.GetName()))
							return
						}
					}
				}

				if ep.Interval == 0 {
					logger.Log.Info("pausing plan due to zero interval", zap.String("plan", *p.Name))
					return
				}
			}
		}
	}
}

func InitPlan(p Plan) Plan {
	return p
}
