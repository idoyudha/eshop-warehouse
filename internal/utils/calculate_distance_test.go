package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateZipCodeDistance(t *testing.T) {
	// allow this function run in parallel with other test function
	// t.Parallell()

	tests := []struct {
		name          string
		zipCode       string
		zipCodes      map[string]string
		expectedDist  map[string]float64
		expectedError bool
	}{
		{
			name:    "valid zip codes with positive distances",
			zipCode: "10000",
			zipCodes: map[string]string{
				"20000": "20000",
				"30000": "30000",
			},
			expectedDist: map[string]float64{
				"20000": 10000,
				"30000": 20000,
			},
			expectedError: false,
		},
		{
			name:    "valid zip codes with negative differences",
			zipCode: "30000",
			zipCodes: map[string]string{
				"10000": "10000",
				"20000": "20000",
			},
			expectedDist: map[string]float64{
				"10000": 20000,
				"20000": 10000,
			},
			expectedError: false,
		},
		{
			name:    "invalid input zipcode",
			zipCode: "abc",
			zipCodes: map[string]string{
				"10000": "10000",
			},
			expectedDist:  nil,
			expectedError: true,
		},
		{
			name:    "zero distance",
			zipCode: "10000",
			zipCodes: map[string]string{
				"10000": "10000",
			},
			expectedDist: map[string]float64{
				"10000": 0,
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test case will run in parallel
			// t.Parallell()

			distances, err := CalculateZipCodeDistance(tt.zipCode, tt.zipCodes)

			if (err != nil) != tt.expectedError {
				t.Errorf("calculateZipCodeDistance() error = %v, expected error %v", err, tt.expectedError)
				return
			}

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, distances)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedDist, distances, "distances should match expected values")
		})
	}
}
