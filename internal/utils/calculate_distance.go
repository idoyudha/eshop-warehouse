package utils

import (
	"math"
	"strconv"
)

func CalculateZipCodeDistance(zipCode string, zipCodes map[string]string) (map[string]float64, error) {
	zipCodeDistances := make(map[string]float64)

	zipCodeFloat, err := strconv.ParseFloat(zipCode, 64)
	if err != nil {
		return nil, err
	}

	for zipcode, distance := range zipCodes {
		distanceFloat, err := strconv.ParseFloat(distance, 64)
		if err != nil {
			return nil, err
		}
		zipCodeDistances[zipcode] = math.Abs(zipCodeFloat - distanceFloat)
	}

	return zipCodeDistances, nil
}
