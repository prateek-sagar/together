package handlers

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"time"
	"database/sql"
	

	"github.com/google/uuid"

	"together/models"
)

func generateSpaceCode() (string, error) {
	length := 7

	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	code := make([]byte, length)

	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}

		code[i] = chars[n.Int64()]
	}

	return string(code), nil
}


func CreateSpace(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	var name struct {
		Name string `json: "name"`
	}

	err := json.NewDecoder(r.Body).Decode(&name)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if name.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

    code, err := generateSpaceCode()
    if err != nil {
        http.Error(w, "Failed to generate space code", http.StatusInternalServerError)
        return
    }

    space := models.Space{
        ID:        uuid.New(),
        Code:      code,
        IsActive:  true,
        CreatedAt: time.Now().UTC(),
    }

    _, err = db.Exec(`
        INSERT INTO spaces (id, code, is_active, created_at)
        VALUES (?, ?, ?, ?)
    `,
        space.ID.String(),
        space.Code,
        space.IsActive,
        space.CreatedAt,	
    )

	CreateMember(space.ID, name.Name, true)

    if err != nil {
        http.Error(w, "Failed to save space", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Space created successfully",
        "space":   space,
    })
}


func GetSpace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get space ID from query parameter
	spaceID := r.URL.Query().Get("id")

	if spaceID == "" {
		http.Error(w, "Space ID is required", http.StatusBadRequest)
		return
	}

	// Validate UUID
	id, err := uuid.Parse(spaceID)
	if err != nil {
		http.Error(w, "Invalid space ID", http.StatusBadRequest)
		return
	}

	var space models.Space
	var dbID string

	err = db.QueryRow(`
		SELECT id, code, is_active, created_at
		FROM spaces
		WHERE id = ?
	`, id.String()).Scan(
		&dbID,
		&space.Code,
		&space.IsActive,
		&space.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Space not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to fetch space", http.StatusInternalServerError)
		return
	}

	space.ID, err = uuid.Parse(dbID)
	if err != nil {
		http.Error(w, "Invalid space ID in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Space fetched successfully",
		"space":   space,
	})
}

func EnterSpace(w http.ResponseWriter, r *http.Request) {
	// require code in the request body
	var requestData struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestData.Code == "" {
		http.Error(w, "Space code is required", http.StatusBadRequest)
		return
	}

	if requestData.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	var space models.Space
	var dbID string

	err = db.QueryRow(`
		SELECT id, code, is_active, created_at
		FROM spaces
		WHERE code = ?
	`, requestData.Code).Scan(
		&dbID,
		&space.Code,
		&space.IsActive,
		&space.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Space not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to fetch space", http.StatusInternalServerError)
		return
	}

	if !space.IsActive {
		http.Error(w, "Space is inactive", http.StatusForbidden)
		return
	}

	space.ID, err = uuid.Parse(dbID)
	if err != nil {
		http.Error(w, "Invalid space ID in database", http.StatusInternalServerError)
		return
	}


	IsExisting, err := IsExistingMember(space.ID, requestData.Name)
	if err != nil {
		http.Error(w, "Failed to check existing member", http.StatusInternalServerError)
		return
	}

	if !IsExisting {
		CreateMember(space.ID, requestData.Name, false)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Space entered successfully",
		"space":   space,
	})

}

func GetCode(w http.ResponseWriter, r *http.Request) {

}
