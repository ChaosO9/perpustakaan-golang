package utils

import "time"

func CalculateFine(dueDate, returnedAt time.Time) int64 {
	diff := returnedAt.Sub(dueDate)
	if diff > 0 {
		return int64(diff.Hours() * 1000)
	}
	return 0
}