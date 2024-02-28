package main

import (
	"fmt"
	"strconv"
	"strings"

	dbutils "fandm/internal/database/_dbutils"
	"fandm/internal/routes"
)

func main() {
	str := "item_1:12,item_2:32,item_3:78"
	items := strings.Split(str, ",")
	m := make(map[string]int)
	for _, item := range items {
		parts := strings.Split(item, ":")
		if len(parts) == 2 {
			key := parts[0]
			value, err := strconv.Atoi(parts[1])
			if err == nil {
				m[key] = value
			}
		}
	}
	fmt.Println(m)
	dbutils.SetUpDatabase()
	routes.Routes()
}
