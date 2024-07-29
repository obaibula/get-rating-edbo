package main

import (
	"fmt"
	"log/slog"
)

const (
	alinaRating = 164.64
	maxQuota    = 10 // number of places for people with quota (will be occupied in any case)
)

func main() {
	students, err := getStudents()
	if err != nil {
		slog.Error("Error getting students", "msg", err)
		return
	}
	quotaWithFirstPriority := getNumberOfQuotaWithPriority(students, First)
	placesOccupiedByQuota := min(maxQuota, quotaWithFirstPriority)

	fmt.Println("Your place based on first priority is:", getCurrentPlaceBasedOnPriority(students, First)+placesOccupiedByQuota)
	fmt.Println("Your place based on second priority is:", getCurrentPlaceBasedOnPriority(students, Second)+placesOccupiedByQuota)
	fmt.Println("Your place based on third priority is:", getCurrentPlaceBasedOnPriority(students, Third)+placesOccupiedByQuota)
	fmt.Println("Your place based on fourth priority is:", getCurrentPlaceBasedOnPriority(students, Fourth)+placesOccupiedByQuota)
	fmt.Println("Your place based on fifth priority is:", getCurrentPlaceBasedOnPriority(students, Fifth)+placesOccupiedByQuota)
}

func getNumberOfQuotaWithPriority(students []student, priority priority) int {
	predicate := hasQuota().
		and(priorityAbove(priority)).
		and(notRejected())

	return len(filter(students, predicate))
}

func getCurrentPlaceBasedOnPriority(students []student, priority priority) int {
	predicate := hasQuota().negate().
		and(notRejected()).
		and(ratingAboveOrEqual(alinaRating)).
		and(priorityAbove(priority))

	return len(filter(students, predicate))
}
