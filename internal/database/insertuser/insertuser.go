package insertuser

import (
	"fmt"

	user "fandm/internal/_dto/user"
	dbutils "fandm/internal/database/_dbutils"
)

func InsertUser(user *user.IncomingUser) {
	fmt.Println(user)
	dbutils.TestConnection()
}
