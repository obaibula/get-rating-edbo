package main

import (
	"fmt"
)

const (
	alinaRating = 154.64
	maxQuota    = 10 // number of places for people with quota (will be occupied in any case)
)

func main() {
	students, err := getStudents()
	if err != nil {
		fmt.Println("Error getting students", err)
		return
	}
	currentRating := getCurrentRatingBasedOnFirstPriority(students)
	fmt.Println("Your current rating is:", currentRating)

}

func getCurrentRatingBasedOnFirstPriority(students []student) int {
	predicate := hasQuota().negate().
		and(notRejected()).
		and(ratingAboveOrEqual(alinaRating)).
		and(priorityAbove(1))

	return filter(students, predicate)
}
