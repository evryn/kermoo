package planner

import (
	"fmt"
	"kermoo/config"
	"kermoo/modules/logger"
	"kermoo/modules/utils"
	"kermoo/modules/values"
	"time"

	"go.uber.org/zap"
)

type SubPlan struct {
	Percentage *values.MultiFloat `json:"percentage"`
	Size       *values.MultiSize  `json:"size"`
	Interval   *values.Duration   `json:"interval"`
	Duration   *values.Duration   `json:"duration"`

	cycleValues  []CycleValue
	relatedPlan  *Plan
	totalCycles  uint64
	currentCycle uint64
}

type CycleValue struct {
	Percentage               values.SingleFloat
	Size                     values.SingleSize
	ComputedPercentageChance *bool
}

func (cv *CycleValue) ComputeStaticValues() {
	value, _ := cv.Percentage.ToFloat()
	computedPercentageChance := utils.IsSuccessByChance(value)
	cv.ComputedPercentageChance = &computedPercentageChance
}

func (s *SubPlan) computeCycleValues() ([]CycleValue, error) {
	count := 0
	var err error
	var cycleValues []CycleValue
	var singleSizes []values.SingleSize
	var singleValues []values.SingleFloat

	if s.Size != nil {
		singleSizes, err = s.Size.ToSingleSizes()

		if err != nil {
			return nil, fmt.Errorf("failed to convert size to single values: %v", err)
		}

		if len(singleSizes) > count {
			count = len(singleSizes)
		}
	}

	if s.Percentage != nil {
		singleValues, err = s.Percentage.ToSingleFloats()

		if err != nil {
			return nil, fmt.Errorf("failed to convert value to single values: %v", err)
		}

		if len(singleValues) > count {
			count = len(singleValues)
		}
	}

	if len(singleSizes) > 0 && len(singleValues) > 0 && len(singleSizes) != len(singleValues) {
		return nil, fmt.Errorf("both size and values are set while the count of individual steps does not match together")
	}

	for i := 0; i < count; i++ {
		percentage := values.NewZeroFloat()
		size := values.NewZeroSize()

		if len(singleValues) >= i+1 {
			percentage = singleValues[i]
		}

		if len(singleSizes) >= i+1 {
			size = singleSizes[i]
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
		return time.Duration(*s.Interval)
	}

	return config.Default.Planner.Interval
}

func (s *SubPlan) computeRequiredCycles() uint64 {
	dur := time.Duration(*s.Duration)
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
