package main

import (
	"auth-server-proxy/src/application/usecases"
	"auth-server-proxy/src/presentation/handler"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func gettoken(item int) {
	godotenv.Load()
	_logger := handler.NewLoggerHandler()
	apiManagerUrl := os.Getenv("AUTH_SERVER_API_MANAGER")

	path := "/oauth2/token"
	params := make(map[string]string)
	params["grant_type"] = "client_credentials"
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	// headers["Authorization"] = "Basic Q2pzQVE1d2VwdzFxTlZ2UjdIZXhiT0Q0Q2Z3YTpibmlCWEJRZE5KQnNoaGpobHdNcEszcWNFdmth" //OK
	// headers["Authorization"] = "Basic Q2pzQVE1d2VwdzFxTlZ2UjdIZXhiT0Q0Q2Z3YTpibmlCWEJRZE5KQnNoaGpobHdNcEszcWNFdmsy" //Bad Passwd
	headers["Authorization"] = "Basic Q2pzQVE1d2VwdzFxTlZ2UjdIZXhiT0Q0Q2Z3YTpibmlCWEJRZE5KQnNoaGpobHdNcEszcWNFdmszxc" //bad base64

	httpClientApiManager := usecases.NewOauthLoginCase(apiManagerUrl + path)
	result, err := httpClientApiManager.Start(params, headers)

	if err != nil {
		_logger.Error(err, "Error registering token")
		return
	}
	token := ""
	if result.StatusCode == 200 {
		token = result.AccessToken[956:]
	}

	_logger.Info(fmt.Sprint("Server 1 | Worker-", item, "|  {StatusCode:", result.StatusCode, ", PartToken: ", token, "}"))
}

func gettoken2(item int) {
	godotenv.Load()
	_logger := handler.NewLoggerHandler()
	apiManagerUrl := os.Getenv("AUTH_SERVER_API_MANAGER")

	path := "/oauth2/token"
	params := make(map[string]string)
	params["grant_type"] = "client_credentials"
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Authorization"] = "Basic Q2pzQVE1d2VwdzFxTlZ2UjdIZXhiT0Q0Q2Z3YTpibmlCWEJRZE5KQnNoaGpobHdNcEszcWNFdmth" //OK
	// headers["Authorization"] = "Basic Q2pzQVE1d2VwdzFxTlZ2UjdIZXhiT0Q0Q2Z3YTpibmlCWEJRZE5KQnNoaGpobHdNcEszcWNFdmsy" //Bad Passwd
	// headers["Authorization"] = "Basic Q2pzQVE1d2VwdzFxTlZ2UjdIZXhiT0Q0Q2Z3YTpibmlCWEJRZE5KQnNoaGpobHdNcEszcWNFdmszxc" //bad base64

	httpClientApiManager := usecases.NewOauthLoginCase(apiManagerUrl + path)
	result, err := httpClientApiManager.Start(params, headers)

	if err != nil {
		_logger.Error(err, "Error registering token")
		return
	}
	token := ""
	if result.StatusCode == 200 {
		token = result.AccessToken[956:]
	}

	_logger.Info(fmt.Sprint("Server 2 | Worker-", item, "|  {StatusCode:", result.StatusCode, ", PartToken: ", token, "}"))
}

func main() {
	for i := 0; i < 20; i++ {
		go gettoken(i)
		go gettoken2(i)
	}
	exitChan := make(chan struct{})
	<-exitChan
	fmt.Println("Program exited.")
}
