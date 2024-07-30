package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ReadJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("Missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}
func Error(w http.ResponseWriter, status int, err error) {
	fmt.Println("too lazy to finish this")
}
