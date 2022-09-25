package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// getPathNumber extracts the number from the path.
func getPathNumber(req *http.Request) (int, error) {
	vars := mux.Vars(req)
	page := vars["page"]

	num, err := strconv.Atoi(page)
	if err != nil {
		return 0, fmt.Errorf("converting the page number to int: %w", err)
	}

	return num, nil
}
