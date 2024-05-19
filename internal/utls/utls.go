package utls

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Catch(err error) {
	if err != nil {
		log.Fatal(err)
		return
	}
}

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	// (*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func StringToMap(str string) map[string]interface{} {
	m := make(map[string]interface{})
	if str == "" {
		return m
	}
	items := strings.Split(str, ",")
	for _, item := range items {
		if item == "" {
			continue
		}
		parts := strings.Split(item, ":")
		if len(parts) == 2 {
			key := parts[0]
			value, err := strconv.Atoi(parts[1])
			if err == nil {
				m[key] = value
			}
		}
	}
	return m
}

func MapToString(m map[string]interface{}) string {
	var str strings.Builder
	for k, v := range m {
		str.WriteString(fmt.Sprintf("%s:%d,", k, v))
	}
	return strings.TrimRight(str.String(), ",")
}
