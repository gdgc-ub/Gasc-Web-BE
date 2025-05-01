package searchServices

import (
	"ProjectGolang/internal/api/gemini"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)


func (s *searchService) AnalyzeMedicineImage(ctx context.Context, imageFile *multipart.FileHeader) (gemini.GeminiResponse, error) {
	file , err := imageFile.Open()
	if err != nil{
		return gemini.GeminiResponse{}, err
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil { 
		return gemini.GeminiResponse{}, err
	}

	base64Image := base64.StdEncoding.EncodeToString(fileData)

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" { 
		return gemini.GeminiResponse{}, fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	geminiModel := os.Getenv("GEMINI_MODEL")
	if geminiModel == "" {
		return gemini.GeminiResponse{}, fmt.Errorf("GEMINI_MODEL environment variable not set ")
	}

	mimeType := imageFile.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"

		filename := imageFile.Filename
		ext := strings.ToLower(filepath.Ext(filename))
		switch ext { 
		case ".png": 
		mimeType = "image/png"
		case ".jpg", ".jpeg": 
		mimeType = "image/jpeg"
		case ".gif": 
		mimeType = "image/gif"
		case ".webp": 
		mimeType = "image/webp"
		}
	}

	geminiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", geminiModel, geminiAPIKey)

	requestBody := map[string]interface{}{ 
		"contents": [] map[string]interface{} { 
			{ 
				"parts": []map[string]interface{}{ 
					{ 
						"text": `
You are an expert pharmaceutical image analyzer. Examine this image of a medicine or supplement package/container carefully.

TASK: Identify and extract key medicine information from the image.

Look for:
- Brand/product name (on the front of packaging, usually prominent)
- Active ingredients (may be on front or side panels)
- Dosage information (strength per unit, units per container)
- Expiry date (usually formatted as YYYY-MM-DD, MM/YYYY, or similar)
- Other identifiable details (manufacturer, drug class, etc.)

If information is partially visible or unclear, make your best guess and indicate lower confidence.

FORMAT YOUR RESPONSE AS A VALID JSON OBJECT with these fields:
{
  \"medicineName\": \"[product name]\",
  \"activeIngredient\": \"[main active ingredients]\",
  \"dosage\": \"[dosage information]\",
  \"expiryDate\": \"[date in YYYY-MM-DD format if possible]\",
  \"confidenceScore\": [number between 0-1],
  \"notes\": \"[brief notes on any ambiguities or missing information]\"
}

Provide ONLY the JSON object without additional text, explanations, or markdown formatting`,
					},
					{
						"inline_data": map[string]interface{}{ 
							"mime_type": mimeType, 
							"data": base64Image,
						},
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{ 
			"temperature": 0.1, 
			"topP" : 0.8, 
			"topK": 40, 
		},
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil{
		return gemini.GeminiResponse{}, err
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "POST", geminiURL, bytes.NewBuffer(requestJSON))
	if err != nil { 
		return gemini.GeminiResponse{}, err
	}
	req.Header.Set("Content-Type", "application/jsson")

	resp, err := httpClient.Do(req)
	if err != nil { 
		return gemini.GeminiResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { 
		bodyBytes, _ := io.ReadAll(resp.Body)
		return gemini.GeminiResponse{}, fmt.Errorf("Gemini API error: %s", string(bodyBytes))
	}

	var geminiResp struct { 
		Candidates []struct { 
			Content struct {
				Parts [] struct { 
					Text string `json:"text"`
				} `json:"parts"`
			}`json:"content"`
		}`json:"candidates"`
	}


	if err:= json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil{ 
		return gemini.GeminiResponse{}, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 { 
		return gemini.GeminiResponse{}, gemini.ErrGeminiProcessingFailed
	}

	responseText := geminiResp.Candidates[0].Content.Parts[0].Text

	jsonPattern := regexp.MustCompile(`(?s)\{.*\}`)
	matches := jsonPattern.FindString(responseText)
	if matches != "" { 
		responseText = matches
	}

	responseText = strings.TrimSpace(responseText)
	if strings.HasPrefix(responseText, "```json"){ 
		responseText = strings.TrimPrefix(responseText, "```json")
		responseText = strings.TrimSuffix(responseText, "```")
	}else if strings.HasPrefix(responseText, "```"){ 
		responseText = strings.TrimPrefix(responseText, "```")
		responseText = strings.TrimSuffix(responseText, "```")
	}

	var medicineAnalysis gemini.GeminiResponse
	if err := json.Unmarshal([]byte(responseText), &medicineAnalysis); err != nil{ 
		type MedicineAltResponse struct { 
			MedicineName string `json:"medicineName"`
			ActiveIngredient string `json:"activeIngredient"`
			Dosage string `json:"dosage"`
			ExpiryDate string `json:"expiryDate"`
			ConfidenceScore float64 `json:"confidenceScore"`
		}

		var altResponse MedicineAltResponse
		if altErr := json.Unmarshal([]byte(responseText), &altResponse); altErr != nil {
			return gemini.GeminiResponse{}, fmt.Errorf("failed to parse Gemini response: %v - Raw response: %s", err, responseText)
		}

		var expiryDate time.Time
		if altResponse.ExpiryDate != "" {
			parsedDate, dateErr := time.Parse("2006-01-02", altResponse.ExpiryDate)
			if dateErr == nil {
				expiryDate = parsedDate
			} else {
				if altResponse.ExpiryDate != "" {
					altRespInt, intErr := strconv.Atoi(altResponse.ExpiryDate)
					if intErr == nil {
						expiryDate = time.Now().AddDate(0, 0, altRespInt)
					} else {
						expiryDate = time.Now().AddDate(0, 0, 30)
					}
				} else {
					expiryDate = time.Now().AddDate(0, 0, 30)
				}
			}
		} else {
			expiryDate = time.Now().AddDate(0, 0, 30)
		}

		medicineAnalysis = gemini.GeminiResponse{
			MedicineName:     altResponse.MedicineName,
			ActiveIngredient: altResponse.ActiveIngredient,
			Dosage:           altResponse.Dosage,
			ExpiryDate:       expiryDate,
			ConfidenceScore:  altResponse.ConfidenceScore,
		}
	}

	if medicineAnalysis.MedicineName == "" {
		medicineAnalysis.MedicineName = "Unknown Medicine"
	}

	if medicineAnalysis.ExpiryDate.IsZero() {
		medicineAnalysis.ExpiryDate = time.Now().AddDate(0, 0, 30)
	}

	if medicineAnalysis.ConfidenceScore < 0 || medicineAnalysis.ConfidenceScore > 1 {
		medicineAnalysis.ConfidenceScore = 0.5
	}

	return medicineAnalysis, nil
}
