package authentication

import (
	"fmt"
	"net/http"
)

const StatusOnline = "online"

func Status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, StatusOnline)
}
