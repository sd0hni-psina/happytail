package models

import "errors"

type RoleType string

const (
	RoleAdmin        RoleType = "admin"
	RoleShelterAdmin RoleType = "shelter_admin"
	RoleUser         RoleType = "user"
)

var validRole = map[RoleType]struct{}{
	RoleAdmin:        {},
	RoleShelterAdmin: {},
	RoleUser:         {},
}

func (r RoleType) IsValid() bool {
	_, ok := validRole[r]
	return ok
}

type Role struct {
	ID        int      `json:"id"`
	UserID    int      `json:"user_id"`
	RoleType  RoleType `json:"role_type"`
	ShelterID *int     `json:"shelter_id"`
}

type RoleInput struct {
	UserID    int      `json:"user_id"`
	RoleType  RoleType `json:"role_type"`
	ShelterID *int     `json:"shelter_id"`
}

func (ri *RoleInput) Validate() error {
	if !ri.RoleType.IsValid() {
		return errors.New("Invalid role type")
	}
	if ri.RoleType == RoleShelterAdmin && ri.ShelterID == nil {
		return errors.New("shelter_id is required for shelter_admin")
	}
	if ri.RoleType != RoleShelterAdmin && ri.ShelterID != nil {
		return errors.New("shelter_id is only allowed for shelter_admin")
	}
	return nil
}
