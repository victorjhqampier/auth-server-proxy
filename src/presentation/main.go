package main

/******************************************************************
*******************************************************************
*	Oauth2 Authentication Server Proxy for Client Credentiales Flow
*
*	Copyright © 2025 Victor JCaxi
*
*	Este código está licenciado bajo la Licencia Pública General Affero de GNU (AGPL), versión 3.0 o posterior.
*	Permite usar, modificar y distribuir este código, pero cualquier cambio o mejora debe ser contribuido de vuelta al proyecto original y compartido bajo la misma licencia.
*
*	Para obtener más detalles sobre esta licencia, consulta: https://www.gnu.org/licenses/agpl-3.0.html
********************************************************************
********************************************************************/

import (
	"auth-server-proxy/src/presentation/controllers"
	"auth-server-proxy/src/presentation/handler"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	_logger := handler.NewLoggerHandler()

	if err := godotenv.Load(); err != nil {
		log.Println("No se pudieron cargar las variables de entorno")
	}

	r := mux.NewRouter()
	controllers.NewIndexController(r)
	controllers.NewOauthController(r)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "5000"
	}

	_logger.Info("Run server in http://localhost:" + port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
