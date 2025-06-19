package rating

func AverageRatingToString(r float64) string {
	if r >= 1.5 {
		return "положительный"
	} else {
		return "отрицательный"
	}
}

func RatingToString(r int) string {
	switch r {
	case 1:
		return "отрицательный"
	case 2:
		return "положительный"
	default:
		return "не оценен"
	}
}
