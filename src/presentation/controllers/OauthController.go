package controllers

import (
	"auth-server-proxy/src/application/adapters/internals"
	enums "auth-server-proxy/src/application/enums"
	helper "auth-server-proxy/src/application/helpers"
	"auth-server-proxy/src/application/usecases"
	"auth-server-proxy/src/presentation/handler"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type OauthController struct {
	_logger handler.LoggerHandler
	_app    usecases.OauthLoginCase
}

func NewOauthController(r *mux.Router) {
	godotenv.Load()
	controller := &OauthController{
		_logger: *handler.NewLoggerHandler(),
		_app:    usecases.NewOauthLoginCase(os.Getenv("AUTH_SERVER_API_MANAGER") + "/oauth2/token"),
	}

	r.HandleFunc("/wso2apim/oauth2/token", controller.authWso2Apim).Methods("POST")
}

// Added enpoints
func (c *OauthController) authWso2Apim(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extraer los parámetros de la query string como map[string]string
	params := make(map[string]string)
	query := r.URL.Query()
	for key, values := range query {
		params[key] = values[0]
	}

	headers := make(map[string]string)
	for key, values := range r.Header {
		headers[key] = values[0]
	}

	delete(headers, "Accept-Encoding")

	result, err := c._app.Start(params, headers)

	if err != nil {
		c._logger.Error(err, "Error registering token")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.EasyErrorRespond("10099", string(enums.InternalError)))
		return
	}

	if result.StatusCode == 204 {
		w.WriteHeader(result.StatusCode)
		json.NewEncoder(w).Encode(helper.EasyEmptyRespond(204))
		return
	}

	if result.StatusCode != 200 {
		w.WriteHeader(result.StatusCode)
		json.NewEncoder(w).Encode(helper.EasyListErrorRespond([]internals.FieldErrorAdapter{
			{Code: "192", Message: string(enums.RequestError)},
		},
			result.StatusCode))
		return
	}

	w.WriteHeader(result.StatusCode)
	json.NewEncoder(w).Encode(result)
}
