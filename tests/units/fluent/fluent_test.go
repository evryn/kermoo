package fluent_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kermoo/modules/fluent"
)

func TestFluentFloat_GetValues(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		between []float64
		want    []float64
		wantErr bool
	}{
		{
			name:  "single float value",
			input: "2.5",
			want:  []float64{2.5},
		},
		{
			name:    "float range",
			input:   "1 to 3.5",
			between: []float64{1, 3.5},
		},
		{
			name:  "comma-separated floats",
			input: "1.5,2.5,3.5, 4",
			want:  []float64{1.5, 2.5, 3.5, 4},
		},
		{
			name:    "invalid float range format",
			input:   "1.5 to",
			wantErr: true,
		},
		{
			name:    "invalid float value",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "invalid float in comma-separated list",
			input:   "1.5,abc,3.5",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fl, err := fluent.NewFluentFloat(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if len(tt.between) == 2 {
				values := fl.GetArray()
				require.Len(t, values, 1)
				assert.Less(t, tt.between[0], values[0])
				assert.Greater(t, tt.between[1], values[0])
			} else {
				values := fl.GetArray()
				assert.Equal(t, tt.want, values)
			}
		})
	}
}

func TestFluentFloat_GetValue(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expected        float64
		expectedBetween []float64
		expectedIn      []float64
	}{
		{
			name:     "single float value",
			input:    "2.5",
			expected: float64(2.5),
		},
		{
			name:            "float range",
			input:           "1.5 to 3.5",
			expectedBetween: []float64{1.5, 3.5},
		},
		{
			name:       "comma-separated floats",
			input:      "1.5,2.5,3.5",
			expectedIn: []float64{1.5, 2.5, 3.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fl, err := fluent.NewFluentFloat(tt.input)
			require.NoError(t, err)
			value := fl.Get()

			if len(tt.expectedIn) > 0 {
				assert.Contains(t, tt.expectedIn, value)
			} else if len(tt.expectedBetween) > 0 {
				assert.Less(t, tt.expectedBetween[0], value)
				assert.Greater(t, tt.expectedBetween[1], value)
			} else {
				assert.Equal(t, tt.expected, value)
			}

		})
	}
}

func TestFluentFloat_GetCachedValue(t *testing.T) {
	t.Run("single float value", func(t *testing.T) {
		fl, err := fluent.NewFluentFloat("2.5")
		require.NoError(t, err)

		assert.Equal(t, fl.GetCached(), fl.GetCached())
		assert.Equal(t, fl.GetUpdatedCache(), fl.GetCached())
	})

	t.Run("float range", func(t *testing.T) {
		fl, err := fluent.NewFluentFloat("1 to 5.5")
		require.NoError(t, err)

		v1 := fl.GetCached()
		v2 := fl.GetCached()
		v3 := fl.GetUpdatedCache()
		v4 := fl.GetCached()

		assert.Equal(t, v1, v2)
		assert.NotEqual(t, v2, v3)
		assert.Equal(t, v3, v4)
	})

	t.Run("comma-separated floats", func(t *testing.T) {
		fl, err := fluent.NewFluentFloat("1, 2.5, 5, 6.5")
		require.NoError(t, err)

		v1 := fl.GetCached()
		v2 := fl.GetCached()
		v3 := fl.GetUpdatedCache()
		v4 := fl.GetCached()

		assert.Equal(t, v1, v2)
		assert.NotEqual(t, v2, v3)
		assert.Equal(t, v3, v4)
	})
}
