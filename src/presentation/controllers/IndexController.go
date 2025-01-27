package controllers

import (
	helper "auth-server-proxy/src/application/helpers"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type IndexController struct{}

func NewIndexController(r *mux.Router) {
	controller := IndexController{}

	r.HandleFunc("/", controller.getMessage).Methods("GET")
}

// Added enpoints
func (c *IndexController) getMessage(w http.ResponseWriter, r *http.Request) {
	arrData := map[string]string{
		"server":    "Oauth2 Authentication Server Proxy",
		"copyright": "© 2025 Victor JCaxi",
		"license":   "GNU Affero General Public License (AGPL) v3.0",
		"details":   "https://www.gnu.org/licenses/agpl-3.0.html",
	}
	result := helper.EasySuccessRespond(arrData, 200)

	// Responder con JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.StatusCode)
	json.NewEncoder(w).Encode(result)
}
