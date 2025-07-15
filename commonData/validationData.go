package commonData

import (
	"fmt"
	"time"
)

func ValueExists(value string, name string) int8 {
	if value == "" {
		fmt.Println("Error, value not found", name)
		return 0
	}
	return 1
}

func ValueNumberExits(value float64, name string) int8 {
	if value <= 0 {
		fmt.Println("Error, value is 0 or less", name)
		return 0
	}
	return 1
}

func ValueTimeExists(value time.Time, format string, name string) int8 {
	if value.Year() <= 1800 {
		fmt.Println("Error, the date cannot possibly exist", name)
		return 0
	}
	return 1
}
