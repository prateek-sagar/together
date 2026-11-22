package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"together/models"
)

func CreateMedia(w http.ResponseWriter, r *http.Request) {

    var req models.AddMediaRequest

    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request body", 400)
        return
    }

    if req.SpaceID == uuid.Nil ||
        req.MemberID <= 0 ||
        req.Location == "" ||
        req.Hash == "" {

        http.Error(w, "Missing required fields", 400)
        return
    }

	hashID, err := GetOrCreateMediaHash(req.Hash)

	if err != nil {
		http.Error(w, "Failed to save hash", 500)
		return
	}

    // Create media
    media := models.Media{
        ID:       uuid.New(),
        SpaceID:  req.SpaceID,
        MemberID: req.MemberID,
        Location: req.Location,
        HashID:   hashID,
    }

    // Save media
    err = AddMedia(media)
    if err != nil {
        http.Error(w, "Failed to save media", 500)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Media added successfully",
        "media":   media,
    })
}

func AddMedia(media models.Media) error {

    _, err := db.Exec(`
        INSERT INTO media (
            id, space_id, member_id, location, hash_id
        )
        VALUES (?, ?, ?, ?, ?)
    `,
        media.ID.String(),
        media.SpaceID.String(),
        media.MemberID,
        media.Location,
        media.HashID,
    )

    return err
}

func AddMediaHash(hash string) (string, error) {

    hashID := uuid.New().String()

    _, err := db.Exec(`
        INSERT INTO media_hash (hash_id, hash)
        VALUES (?, ?)
    `, hashID, hash)

    if err != nil {
        return "", err
    }

    return hashID, nil
}


func IsHashExists(hash string) (bool, error) {

	var exists bool	

	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM media_hash WHERE hash = ?
		)
	`, hash).Scan(&exists)

	return exists, err

}

func GetHashIDByHash(hash string) (string, error) {

	var hashID string

	err := db.QueryRow(`
		SELECT hash_id FROM media_hash WHERE hash = ?
	`, hash).Scan(&hashID)

	return hashID, err
}

func GetOrCreateMediaHash(hash string) (string, error) {

    hashID := uuid.New().String()

    _, err := db.Exec(`
        INSERT INTO media_hash (hash_id, hash)
        VALUES (?, ?)
        ON CONFLICT(hash) DO NOTHING
    `, hashID, hash)

    if err != nil {
        return "", err
    }

    err = db.QueryRow(`
        SELECT hash_id
        FROM media_hash
        WHERE hash = ?
    `, hash).Scan(&hashID)

    if err != nil {
        return "", err
    }

    return hashID, nil
}
