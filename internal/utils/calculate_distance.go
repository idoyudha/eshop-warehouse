package utils

import (
	"log"
	"math"
	"strconv"
)

func CalculateZipCodeDistance(zipCode string, zipCodes map[string]string) (map[string]float64, error) {
	zipCodeDistances := make(map[string]float64)

	zipCodeFloat, err := strconv.ParseFloat(zipCode, 64)
	if err != nil {
		log.Println("Error parsing zipCode:", err)
		return nil, err
	}

	for zipcode, distance := range zipCodes {
		distanceFloat, err := strconv.ParseFloat(distance, 64)
		if err != nil {
			log.Printf("Error parsing distance for zip code %s: %v\n", zipcode, err)
			return nil, err
		}
		zipCodeDistances[zipcode] = math.Abs(zipCodeFloat - distanceFloat)
	}

	return zipCodeDistances, nil
}
