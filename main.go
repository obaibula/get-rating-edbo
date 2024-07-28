package main

import (
	"fmt"
	"log/slog"
)

const (
	alinaRating = 154.64
	maxQuota    = 10 // number of places for people with quota (will be occupied in any case)
)

func main() {
	students, err := getStudents()
	if err != nil {
		slog.Error("Error getting students", "msg", err)
		return
	}
	fmt.Println("Your place based on first priority is:", getCurrentPlaceBasedOnPriority(students, First))
	fmt.Println("Your place based on second priority is:", getCurrentPlaceBasedOnPriority(students, Second))
	fmt.Println("Your place based on third priority is:", getCurrentPlaceBasedOnPriority(students, Third))
	fmt.Println("Your place based on fourth priority is:", getCurrentPlaceBasedOnPriority(students, Fourth))
	fmt.Println("Your place based on fifth priority is:", getCurrentPlaceBasedOnPriority(students, Fifth))
}

func getCurrentPlaceBasedOnPriority(students []student, priority priority) int {
	predicate := hasQuota().negate().
		and(notRejected()).
		and(ratingAboveOrEqual(alinaRating)).
		and(priorityAbove(priority))

	return filter(students, predicate)
}
