package planner

import (
	"buggybox/modules/common"
	"time"
)

type PlanInternal struct {
	ExecutablePlans []*ExecutablePlan
	Callbacks       []Callbacks
}

type Plan struct {
	Value    *common.MixedValueF
	Interval *time.Duration
	Duration *time.Duration
	Name     *string
	SubPlans []SubPlan
	internal *PlanInternal
}

type Callbacks struct {
	PreSleep  func(ep *ExecutablePlan, ev *ExecutableValue) PlanSignal
	PostSleep func(startedAt time.Time, timeSpent time.Duration) PlanSignal
}

type PlanSignal uint32

const (
	PLAN_SIGNAL_CONTINUE  PlanSignal = iota
	PLAN_SIGNAL_TERMINATE PlanSignal = iota
)

func (p *Plan) SetInternal(pi *PlanInternal) {
	p.internal = pi
}

func (p *Plan) ToSubPlan() SubPlan {
	return SubPlan{
		Value:    p.Value,
		Interval: p.Interval,
		Duration: p.Duration,
	}
}

func (p *Plan) AddCallback(callback Callbacks) {
	p.internal.Callbacks = append(p.internal.Callbacks, callback)
}

func (p *Plan) GetExecutablePlans() ([]*ExecutablePlan, error) {
	if len(p.SubPlans) == 0 {
		subPlan := p.ToSubPlan()
		executablePlan, err := subPlan.ToExecutablePlan()

		if err != nil {
			return nil, err
		}

		return []*ExecutablePlan{executablePlan}, nil
	}

	// If the plans has SubPlans, generate coresponding executable SubPlans
	var ep = []*ExecutablePlan{}

	for _, sp := range p.SubPlans {
		executablePlan, err := sp.ToExecutablePlan()

		if err != nil {
			return nil, err
		}

		ep = append(ep, executablePlan)
	}

	return ep, nil
}

func (p *Plan) Execute(callbacks Callbacks) error {
	if len(p.internal.ExecutablePlans) == 0 {
		executablePlans, err := p.GetExecutablePlans()

		if err != nil {
			return err
		}

		p.internal.ExecutablePlans = executablePlans
	}

	for _, ep := range p.internal.ExecutablePlans {
		for ep.IsForever || ep.CurrentTries <= ep.TotalTries {
			for _, ev := range ep.Values {
				t := time.Now()

				if !ep.IsForever {
					ep.CurrentTries++

					if ep.CurrentTries > ep.TotalTries {
						break
					}
				}

				if callbacks.PreSleep(ep, ev) == PLAN_SIGNAL_TERMINATE {
					return nil
				}

				time.Sleep(ep.Interval)

				if callbacks.PostSleep(t, time.Since(t)) == PLAN_SIGNAL_TERMINATE {
					return nil
				}
			}
		}
	}

	return nil
}

func (p *Plan) ExecuteAll() error {
	if len(p.internal.ExecutablePlans) == 0 {
		executablePlans, err := p.GetExecutablePlans()

		if err != nil {
			return err
		}

		p.internal.ExecutablePlans = executablePlans
	}

	for _, ep := range p.internal.ExecutablePlans {
		for ep.IsForever || ep.CurrentTries <= ep.TotalTries {
			for _, ev := range ep.Values {
				t := time.Now()

				if !ep.IsForever {
					ep.CurrentTries++

					if ep.CurrentTries > ep.TotalTries {
						break
					}
				}

				for _, c := range p.internal.Callbacks {
					if c.PreSleep(ep, ev) == PLAN_SIGNAL_TERMINATE {
						return nil
					}
				}

				time.Sleep(ep.Interval)

				for _, c := range p.internal.Callbacks {
					if c.PostSleep(t, time.Since(t)) == PLAN_SIGNAL_TERMINATE {
						return nil
					}
				}

			}
		}
	}

	return nil
}

func InitPlan(p Plan) Plan {
	p.SetInternal(&PlanInternal{})

	if p.Value == nil {
		p.Value = &common.MixedValueF{}
	}

	return p
}
