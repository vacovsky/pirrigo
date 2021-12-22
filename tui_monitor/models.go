package main

import (
	"log"
	"strconv"
)

type Chart struct {
	Series []string
	Labels []float64
	Data   [][]string
}

func convertStrToFloats(input []string) []float64 {
	result := []float64{}
	for _, v := range input {
		converted, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Println(err)
		}
		result = append(result, converted)
	}
	return result
}
