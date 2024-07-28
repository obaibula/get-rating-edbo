package main

type studentPredicate func(student) bool

func (sp studentPredicate) and(other studentPredicate) studentPredicate {
	return func(s student) bool {
		return sp(s) && other(s)
	}
}

func (sp studentPredicate) or(other studentPredicate) studentPredicate {
	return func(s student) bool {
		return sp(s) || other(s)
	}
}

func (sp studentPredicate) negate() studentPredicate {
	return func(s student) bool {
		return !sp(s)
	}
}

func filter(students []student, predicate studentPredicate) int {
	result := make([]student, 0, len(students))
	for _, s := range students {
		if predicate(s) {
			result = append(result, s)
		}
	}
	return len(result) + maxQuota
}

func hasQuota() studentPredicate {
	return func(s student) bool {
		for _, q := range s.Quota {
			if q.Quota != "" {
				return true
			}
		}
		return false
	}
}

func notRejected() studentPredicate {
	return func(s student) bool {
		return s.Status != Rejected
	}
}

func ratingAboveOrEqual(threshold float64) studentPredicate {
	return func(s student) bool {
		return s.Rating >= threshold
	}
}

func priorityAbove(priority priority) studentPredicate {
	return func(s student) bool {
		return s.Priority <= priority && s.Priority > 0 // 0 priority is reserved for contract
	}
}
