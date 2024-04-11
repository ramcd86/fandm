package relationshiphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Relationship struct {
	ActorName     string `json:"actor_name"`
	ConditionName string `json:"condition_name"`
	TreatmentName string `json:"treatment_name"`
}

func CreateNewRelationship(w http.ResponseWriter, r *http.Request) {
	var incomingRelationship Relationship
	err := json.NewDecoder(r.Body).Decode(&incomingRelationship)
	if err != nil {
		http.Error(w, string("No data sent."), http.StatusBadRequest)
		return
	}

	fmt.Println("incomingRelationship: ", incomingRelationship)

	response, err := json.Marshal(incomingRelationship)
	if err != nil {
		http.Error(w, string("Error marshalling data."), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// func Register(w http.ResponseWriter, r *http.Request) {
// 	var user user.IncomingUser
// 	var errorSlice []string
//
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		http.Error(w, string("No data sent."), http.StatusBadRequest)
// 		return
// 	}
//
// 	errorSlice = validateUser(&user)
// 	if len(errorSlice) > 0 {
// 		errListJSON, _ := json.Marshal(errorSlice)
// 		http.Error(w, string(errListJSON), http.StatusBadRequest)
// 		return
// 	}
//
// 	user.SetUUID()
// 	user.SetRegistrationDate()
//
// 	insertuser.InsertUser(&user)
//
// 	w.WriteHeader(http.StatusOK)
// }
