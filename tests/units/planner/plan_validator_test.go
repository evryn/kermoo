package planner_test

// import (
// 	"testing"
// )

// func TestPlanValidator(t *testing.T) {
// 	t.Run("invalid name when required but not given", func(t *testing.T) {
// 		fixedValue := float32(1.0)
// 		p := &Plan{FixedValue: &fixedValue}
// 		err := p.Validate(true)
// 		if err == nil {
// 			t.Error("Expected an error when name is required but not provided")
// 		}
// 	})

// 	t.Run("valid name when required and given", func(t *testing.T) {
// 		fixedValue := float32(1.0)
// 		name := "test"
// 		p := &Plan{Name: &name, FixedValue: &fixedValue}
// 		err := p.Validate(true)
// 		if err != nil {
// 			t.Error("Expected no error when name is required and provided")
// 		}
// 	})

// 	t.Run("no name required", func(t *testing.T) {
// 		fixedValue := float32(1.0)
// 		p := &Plan{FixedValue: &fixedValue}
// 		err := p.Validate(false)
// 		if err != nil {
// 			t.Error("Expected no error when name is not required")
// 		}
// 	})

// 	t.Run("FixedValue used with Maximum", func(t *testing.T) {
// 		fixedValue := float32(10)
// 		maximum := float32(20)
// 		p := &Plan{FixedValue: &fixedValue, Maximum: &maximum}
// 		err := p.Validate(false)
// 		if err == nil {
// 			t.Error("Expected an error when FixedValue is used with Maximum")
// 		}
// 	})

// 	t.Run("FixedValue used with Minimum", func(t *testing.T) {
// 		fixedValue := float32(10)
// 		minimum := float32(5)
// 		p := &Plan{FixedValue: &fixedValue, Minimum: &minimum}
// 		err := p.Validate(false)
// 		if err == nil {
// 			t.Error("Expected an error when FixedValue is used with Minimum")
// 		}
// 	})

// 	t.Run("valid plan with FixedValue", func(t *testing.T) {
// 		fixedValue := float32(10)
// 		p := &Plan{FixedValue: &fixedValue}
// 		err := p.Validate(false)
// 		if err != nil {
// 			t.Error("Expected no error for a valid plan with FixedValue")
// 		}
// 	})

// 	t.Run("valid plan with Minimum", func(t *testing.T) {
// 		minimum := float32(5)
// 		p := &Plan{Minimum: &minimum}
// 		err := p.Validate(false)
// 		if err != nil {
// 			t.Error("Expected no error for a valid plan with Minimum")
// 		}
// 	})

// 	t.Run("valid plan with Maximum", func(t *testing.T) {
// 		maximum := float32(20)
// 		p := &Plan{Maximum: &maximum}
// 		err := p.Validate(false)
// 		if err != nil {
// 			t.Error("Expected no error for a valid plan with Maximum")
// 		}
// 	})

// 	t.Run("valid plan with Phases", func(t *testing.T) {
// 		Phases := []PlanPhase{{Chart: PlanPhaseChart{Bars: []float32{1, 2}}}}
// 		p := &Plan{Phases: Phases}
// 		err := p.Validate(false)
// 		if err != nil {
// 			t.Error("Expected no error for a valid plan with Phases")
// 		}
// 	})

// 	t.Run("Maximum is less than Minimum", func(t *testing.T) {
// 		maximum := float32(5)
// 		minimum := float32(10)
// 		p := &Plan{Maximum: &maximum, Minimum: &minimum}
// 		err := p.Validate(false)
// 		if err == nil {
// 			t.Error("Expected an error when Maximum is less than Minimum")
// 		}
// 	})

// 	t.Run("Maximum is equal to Minimum", func(t *testing.T) {
// 		value := float32(10)
// 		p := &Plan{Maximum: &value, Minimum: &value}
// 		err := p.Validate(false)
// 		if err != nil {
// 			t.Error("Expected no error when Maximum is equal to Minimum")
// 		}
// 	})

// 	t.Run("Maximum is greater than Minimum", func(t *testing.T) {
// 		maximum := float32(20)
// 		minimum := float32(10)
// 		p := &Plan{Maximum: &maximum, Minimum: &minimum}
// 		err := p.Validate(false)
// 		if err != nil {
// 			t.Error("Expected no error when Maximum is greater than Minimum")
// 		}
// 	})
// }
