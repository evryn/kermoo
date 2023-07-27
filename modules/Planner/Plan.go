package Planner

import (
	"time"
)

type PlanInternal struct {
	ExecutablePlans []ExecutablePlan
}

type Plan struct {
	Value    *Value
	Interval *time.Duration
	Duration *time.Duration
	Name     *string
	SubPlans []SubPlan
	internal *PlanInternal
}

type Callbacks struct {
	PreSleep  func(ep ExecutablePlan, ev ExecutableValue) PlanSignal
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

func (p *Plan) GetExecutablePlans() []ExecutablePlan {
	if len(p.SubPlans) == 0 {
		sp := p.ToSubPlan()
		return []ExecutablePlan{
			sp.ToExecutablePlan(),
		}
	}

	// If the plans has SubPlans, generate coresponding executable SubPlans
	var ep = []ExecutablePlan{}

	for _, sp := range p.SubPlans {
		ep = append(ep, sp.ToExecutablePlan())
	}

	return ep
}

func (p *Plan) Execute(callbacks Callbacks) {
	if len(p.internal.ExecutablePlans) == 0 {
		p.internal.ExecutablePlans = p.GetExecutablePlans()
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
					return
				}

				time.Sleep(ep.Interval)

				if callbacks.PostSleep(t, time.Since(t)) == PLAN_SIGNAL_TERMINATE {
					return
				}
			}
		}
	}
}
