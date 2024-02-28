package register

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
	"unicode"

	user "fandm/internal/_dto/user"
	insertuser "fandm/internal/database/insertuser"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user user.IncomingUser
	var errorSlice []string

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, string("No data sent."), http.StatusBadRequest)
		return
	}

	errorSlice = validateUser(&user)
	if len(errorSlice) > 0 {
		errListJSON, _ := json.Marshal(errorSlice)
		http.Error(w, string(errListJSON), http.StatusBadRequest)
		return
	}

	user.SetUUID()
	user.SetRegistrationDate()

	insertuser.InsertUser(&user)

	w.WriteHeader(http.StatusOK)
}

func validateUser(incomingUser *user.IncomingUser) []string {
	errList := []string{}
	checkField := func(field string, key string) {
		if field == "" || len(field) < 5 {
			errList = append(errList, "Property '"+key+"' is either missing, or is not at least 5 characters long")
		}
		if key == "email" {
			_, err := mail.ParseAddress(field)
			if err != nil {
				errList = append(errList, "Property 'email' is not a valid email")
			}
		}
		if key == "password" && !validatePassword(field) {
			errList = append(errList, "Property 'password' must contain at least one uppercase letter and one special character")
		}
	}
	checkField(incomingUser.Username, "username")
	checkField(incomingUser.DisplayName, "display_Name")
	checkField(incomingUser.Password, "password")
	checkField(incomingUser.Email, "email")
	return errList
}

func validatePassword(s string) bool {
	hasUpper := strings.IndexFunc(s, unicode.IsUpper) >= 0
	hasSpecial := strings.IndexFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	}) >= 0

	return hasUpper && hasSpecial
}
