package seasonal

import "time"

// GetSuggestion returns products based on the current month
func GetSuggestion() (string, []string) {
	month := time.Now().Month()

	switch month {
	case time.December, time.January, time.February:
		return "Summer", []string{"Fans", "Air Conditioning", "Sunscreen"}

	case time.June, time.July, time.August:
		return "Winter", []string{"Heaters", "Blankets", "Thermal Jackets"}

	default:
		return "Spring/Autumn", []string{"Smartphones", "Laptops", "Books"}
	}
}
