package gemini

import "time"

type GeminiResponse struct {
	MedicineName     string    `json:"medicineName"`
	ActiveIngredient string    `json:"activeIngredient"`
	Dosage           string    `json:"dosage"`
	ExpiryDate       time.Time `json:"expiryDate"`
	ConfidenceScore  float64   `json:"confidenceScore"`
}
