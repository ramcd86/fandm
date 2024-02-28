package user

import (
	"time"

	"github.com/google/uuid"
)

type IncomingUser struct {
	Username         string `json:"username"`
	DisplayName      string `json:"display_name"`
	Password         string `json:"password"`
	RegistrationDate string
	Email            string `json:"email"`
	UUID             string
}

func (u *IncomingUser) SetUUID() {
	u.UUID = uuid.New().String()
}

func (u *IncomingUser) SetRegistrationDate() {
	u.RegistrationDate = time.Now().Format(time.RFC3339)
}

type InternalUser struct {
	Username         string
	Displayname      string
	Hash             string
	Salt             string
	Authkey          string
	AuthkeyExpiry    string
	Email            string
	Uuid             string
	RegistrationDate string
}
