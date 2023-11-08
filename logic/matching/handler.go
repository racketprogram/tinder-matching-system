package matching

import (
	"encoding/json"
	"math"
	"net/http"
	"sync"
)

// Person struct defines the attributes of a single person.
type Person struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Height      int    `json:"height"`
	Gender      string `json:"gender"`
	WantedDates int    `json:"wanted_dates"`
}

// MatchingSystem struct contains all single individuals and provides matching functionality.
type MatchingSystem struct {
	Users     map[int]*Person // Use a map to store user information for quick retrieval.
	lastID    int             // Tracks the last assigned ID.
	mu        sync.Mutex      // Mutex to protect the Users map.
}

// PersonRequest is the structure of the HTTP request body for adding a new user.
type PersonRequest struct {
	Name        string `json:"name"`
	Height      int    `json:"height"`
	Gender      string `json:"gender"`
	WantedDates int    `json:"wanted_dates"`
}

// QuerySinglePeopleRequest is the structure of the HTTP request body for query single people.
type QuerySinglePeopleRequest struct {
	ID         int `json:"id"`
	MatchCount int `json:"match_count"`
}

// RemoveSinglePersonRequest is the structure of the HTTP request body for remvoe one people.
type RemoveSinglePersonRequest struct {
	ID int `json:"id"`
}

// NewMatchingSystem initializes the matching system.
func NewMatchingSystem() *MatchingSystem {
	return &MatchingSystem{
		Users: make(map[int]*Person),
	}
}

// generateID creates a new unique ID.
func (ms *MatchingSystem) generateID() int {
	ms.lastID++      // Increment lastID to generate a new unique ID.
	return ms.lastID // Return the newly generated ID.
}

// AddSinglePersonAndMatchHandler handles the HTTP request to add a new user and find matches.
func (ms *MatchingSystem) AddSinglePersonAndMatchHandler(w http.ResponseWriter, r *http.Request) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var req PersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new Person struct for the new user.
	newPerson := &Person{
		ID:          ms.generateID(), // Call generateID to generate a unique ID.
		Name:        req.Name,
		Height:      req.Height,
		Gender:      req.Gender,
		WantedDates: req.WantedDates,
	}

	ms.Users[newPerson.ID] = newPerson

	// Call the function to add and find matches.
	matches := ms.findMatches(newPerson, math.MaxInt)

	// Encode the match results into JSON and write to the response body.
	if err := json.NewEncoder(w).Encode(matches); err != nil {
		http.Error(w, "Failed to encode matches", http.StatusInternalServerError)
	}
}

// QuerySinglePeople finds the most N possible matched single people for a given user.
func (ms *MatchingSystem) QuerySinglePeople(w http.ResponseWriter, r *http.Request) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var req QuerySinglePeopleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := req.ID
	N := req.MatchCount

	person, exists := ms.Users[userID]
	if !exists {
		http.Error(w, "User not exist.", http.StatusBadRequest)
		return
	}

	// Call the function to add and find matches.
	matches := ms.findMatches(person, N)

	// Encode the match results into JSON and write to the response body.
	if err := json.NewEncoder(w).Encode(matches); err != nil {
		http.Error(w, "Failed to encode matches", http.StatusInternalServerError)
	}
}

// findMatches finds matches for the given user.
func (ms *MatchingSystem) findMatches(person *Person, N int) []*Person {
	var matches []*Person
	for _, potentialMatch := range ms.Users {
		// Avoid matching with oneself.
		if potentialMatch.ID == person.ID {
			continue
		}
		if len(matches) >= N {
			break
		}
		if ms.isMatch(person, potentialMatch) {
			matches = append(matches, potentialMatch)
			// Decrease WantedDates.
			person.WantedDates--
			potentialMatch.WantedDates--
			// Check if WantedDates has reached zero.
			if person.WantedDates == 0 {
				delete(ms.Users, person.ID)
			}
			if potentialMatch.WantedDates == 0 {
				delete(ms.Users, potentialMatch.ID)
			}
		}
	}
	return matches
}

// isMatch checks if two persons match each other based on the matching rules.
func (ms *MatchingSystem) isMatch(person1, person2 *Person) bool {
	if person1.Gender == "male" && person1.WantedDates > 0 && person1.Height > person2.Height && person2.WantedDates > 0 {
		return true
	} else if person1.Gender == "female" && person1.WantedDates > 0 && person1.Height < person2.Height && person2.WantedDates > 0 {
		return true
	}
	return false
}

// RemoveSinglePerson removes a user from the system.
func (ms *MatchingSystem) RemoveSinglePerson(w http.ResponseWriter, r *http.Request) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var req RemoveSinglePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	delete(ms.Users, req.ID)
}
