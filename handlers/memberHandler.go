package handlers

import (
	"log"
	"database/sql"
	"github.com/google/uuid"

	"together/models"
)


func CreateMember(space_id uuid.UUID, name string, isHost bool) {

	member := models.Member{
		SpaceID: space_id,
		Name:    name,
		IsHost:  isHost,
	}

	_, err := db.Exec("INSERT INTO members (space_id, name, is_host) VALUES (?, ?, ?)",
		member.SpaceID,
		member.Name,
		member.IsHost,
	)

	if err != nil {
		log.Println("Failed to create member:", err)
	}

}


func GetMembers(space_id uuid.UUID) ([]models.Member, error) {

    var members []models.Member

    rows, err := db.Query(`
        SELECT id, space_id, name, is_host
        FROM members
        WHERE space_id = ?
    `, space_id.String())

    if err != nil {
        return nil, err
    }

    defer rows.Close()

    for rows.Next() {

        var member models.Member

        err := rows.Scan(
            &member.ID,
            &member.SpaceID,
            &member.Name,
            &member.IsHost,
        )

        if err != nil {
            return nil, err
        }

        members = append(members, member)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return members, nil
}

func GetMembersByName(name string) (models.Member, error) {

	var member models.Member
	
	err := db.QueryRow("SELECT id, space_id, name, is_host FROM members WHERE name = ? limit 1", name).Scan(
		&member.ID,
		&member.SpaceID,
		&member.Name,
		&member.IsHost,
	)

	if err != nil {
		return models.Member{}, err
	}


	return member, nil

}

func IsExistingMember(space_id uuid.UUID, name string) (bool, error) {

	var member models.Member
	
	err := db.QueryRow("SELECT id, space_id, name, is_host FROM members WHERE space_id = ? and name = ? limit 1", space_id, name).Scan(
		&member.ID,
		&member.SpaceID,
		&member.Name,
		&member.IsHost,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil

}
