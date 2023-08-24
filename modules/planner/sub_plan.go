package planner

import (
	"fmt"
	"kermoo/config"
	"kermoo/modules/fluent"
	"kermoo/modules/logger"
	"kermoo/modules/utils"
	"time"

	"go.uber.org/zap"
)

type SubPlan struct {
	Percentage *fluent.FluentFloat    `json:"percentage"`
	Size       *fluent.FluentSize     `json:"size"`
	Interval   *fluent.FluentDuration `json:"interval"`
	Duration   *fluent.FluentDuration `json:"duration"`

	cycleValues  []CycleValue
	relatedPlan  *Plan
	totalCycles  uint64
	currentCycle uint64
}

type CycleValue struct {
	Percentage               float64
	Size                     int64
	ComputedPercentageChance *bool
}

func (cv *CycleValue) ComputeStaticValues() {
	computedPercentageState := utils.PercentageToBoolean(cv.Percentage)
	cv.ComputedPercentageChance = &computedPercentageState
}

func (s *SubPlan) computeCycleValues() ([]CycleValue, error) {
	var cycleValues []CycleValue

	var sizes []int64
	if s.Size != nil {
		sizes = s.Size.GetArray()
	}

	count := len(sizes)

	var percentages []float64
	if s.Percentage != nil {
		percentages = s.Percentage.GetArray()
	}

	if len(percentages) > count {
		count = len(percentages)
	}

	if len(sizes) > 0 && len(percentages) > 0 && len(sizes) != len(percentages) {
		return nil, fmt.Errorf("both size and percentage are set while the count of individual items does not match together")
	}

	for i := 0; i < count; i++ {
		percentage := float64(0)
		size := int64(0)

		if len(percentages) >= i+1 {
			percentage = percentages[i]
		}

		if len(sizes) >= i+1 {
			size = sizes[i]
		}

		cycleValues = append(cycleValues, CycleValue{
			Percentage: percentage,
			Size:       size,
		})
	}

	return cycleValues, nil
}

func (s *SubPlan) getInterval() time.Duration {
	if s.Interval != nil {
		return s.Interval.Get()
	}

	return config.Default.Planner.Interval
}

func (s *SubPlan) computeRequiredCycles() uint64 {
	dur := s.Duration.Get()
	return uint64(dur.Nanoseconds() / s.getInterval().Nanoseconds())
}

func (s *SubPlan) SetPlan(plan *Plan) {
	s.relatedPlan = plan
}

// Prepare makes the sub plan ready for execution
func (s *SubPlan) Prepare() error {
	var err error

	if !s.isEndless() {
		s.totalCycles = s.computeRequiredCycles()
	}
	s.cycleValues, err = s.computeCycleValues()

	if err != nil {
		return err
	}

	return nil
}

func (s *SubPlan) Execute() {
	for {
		for _, cycleValue := range s.cycleValues {
			if !s.NextCycle() {
				return
			}

			startedAt := time.Now()

			cycleValue.ComputeStaticValues()
			s.relatedPlan.SetCurrentValue(cycleValue)

			logger.Log.Info("executing preSleep hooks...", zap.String("plan", *s.relatedPlan.Name))
			if !s.RunPlannableHooks(startedAt, cycleValue, "preSleep") {
				logger.Log.Info("terminating plan by signal", zap.String("plan", *s.relatedPlan.Name))
				return
			}

			time.Sleep(s.getInterval())

			logger.Log.Info("executing postSleep hooks...", zap.String("plan", *s.relatedPlan.Name))
			if !s.RunPlannableHooks(startedAt, cycleValue, "postSleep") {
				logger.Log.Info("terminating plan by signal", zap.String("plan", *s.relatedPlan.Name))
				return
			}

			if s.getInterval() == 0 {
				logger.Log.Info("pausing plan due to zero interval", zap.String("plan", *s.relatedPlan.Name))
				return
			}
		}
	}
}

func (s *SubPlan) Validate() error {
	return s.Prepare()
}

func (s *SubPlan) isEndless() bool {
	return s.Duration == nil
}

func (s *SubPlan) NextCycle() bool {
	if s.isEndless() {
		return true
	}

	if s.currentCycle >= s.totalCycles {
		return false
	}

	s.currentCycle++

	return true
}

func (s *SubPlan) RunPlannableHooks(startedAt time.Time, cv CycleValue, hookType string) bool {
	for _, pl := range s.relatedPlan.plannables {
		plannable := *pl

		var hook *HookFunc
		if hookType == "preSleep" {
			hook = plannable.GetPlanCycleHooks().PreSleep
		} else {
			hook = plannable.GetPlanCycleHooks().PostSleep
		}

		if hook != nil {
			executable := *hook
			value := executable(Cycle{
				Value:     cv,
				StartedAt: startedAt,
				TimeSpent: time.Since(startedAt),
			})

			if value == PLAN_SIGNAL_TERMINATE {
				return false
			}
		}
	}

	return true
}
