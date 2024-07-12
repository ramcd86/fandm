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

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println("ok")

		// Next
		next.ServeHTTP(w, r)
		return
	})
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
