package handlers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"database/sql"
)

var db *sql.DB

func SetDB(database *sql.DB) {
    db = database
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Home Page!")
}

func GetMembersHandler(w http.ResponseWriter, r *http.Request) {

    // Read query parameter
    spaceID := r.URL.Query().Get("space_id")

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

    // Call database function
    members, err := GetMembers(id)
    if err != nil {
        http.Error(w, "Failed to fetch members", http.StatusInternalServerError)
        return
    }

    // Return JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Members fetched successfully",
        "members": members,
    })
}

func GetMemberByNameHandler(w http.ResponseWriter, r *http.Request) {

    // Read query parameter
    name := r.URL.Query().Get("name")

    if name == "" {
        http.Error(w, "Member name is required", http.StatusBadRequest)
        return
    }

    // Call database function
    member, err := GetMembersByName(name)

    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Member not found", http.StatusNotFound)
            return
        }

        http.Error(w, "Failed to fetch member", http.StatusInternalServerError)
        return
    }

    // Return JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Member fetched successfully",
        "member": member,
    })
}
