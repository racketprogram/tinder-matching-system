package matching

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// TestAddSinglePersonAndMatchHandler tests the HTTP handler for adding a single person and finding matches.
func TestAddSinglePersonAndMatchHandler(t *testing.T) {
	ms := NewMatchingSystem()

	// Create a recorder to capture HTTP responses.
	handler := http.HandlerFunc(ms.AddSinglePersonAndMatchHandler)

	femaleCount := 50
	// Add some initial female users.
	for i := 0; i < femaleCount; i++ {
		// Create a request body for a new female user.
		newFemalePerson := PersonRequest{
			Name:        "Female " + strconv.Itoa(i),
			Height:      150 + i, // Assume height starts from 160cm and increases.
			Gender:      "female",
			WantedDates: 3,
		}
		requestBody, _ := json.Marshal(newFemalePerson)
		// Create a simulated HTTP request.
		req, err := http.NewRequest("POST", "/add_single_person_and_match", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		// Serve the HTTP request using our handler.
		handler.ServeHTTP(rr, req)
	}

	maleHeight := 220
	wantedDates := 100
	// Create a request body for a new male user.
	newMalePerson := PersonRequest{
		Name:        "Bob",
		Height:      maleHeight, // The male user's height is 170cm, can match females shorter than 170cm.
		Gender:      "male",
		WantedDates: wantedDates,
	}
	requestBody, _ := json.Marshal(newMalePerson)

	// Create a simulated HTTP request.
	req, err := http.NewRequest("POST", "/add_single_person_and_match", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	// Serve the HTTP request using our handler.
	handler.ServeHTTP(rr, req)

	// Check if the returned status code is http.StatusOK.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse the response body to check the match results.
	var matches []*Person
	if err := json.Unmarshal(rr.Body.Bytes(), &matches); err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}

	// Ensure the number of matches is as expected.
	if len(matches) != femaleCount {
		t.Fatalf("Expected match with %d females, got only: %d", femaleCount, len(matches))
	}

	// Validate match results, ensuring each matched female is shorter than 170cm.
	for _, match := range matches {
		if match.Height >= maleHeight {
			t.Errorf("Expected match with height less than 170, got height: %d", match.Height)
		}
	}
}
