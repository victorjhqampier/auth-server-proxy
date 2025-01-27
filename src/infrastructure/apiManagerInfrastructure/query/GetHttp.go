package apiManagerinfraquery

import (
	"auth-server-proxy/src/domain/entity"
	httpinfracollection "auth-server-proxy/src/infrastructure/apiManagerInfrastructure/collection"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type GetHttp struct {
	timeout int
	uri     string
	user    string
}

func NewGetHttp(uri string) GetHttp {
	return GetHttp{
		timeout: 10,
		uri:     uri,
		user:    "",
	}
}
func (_this *GetHttp) printLog(statusCode int, message string) {
	if strings.HasPrefix(_this.user, "Basic ") {
		base64Part := strings.TrimPrefix(_this.user, "Basic ")

		decodedBytes, err := base64.StdEncoding.DecodeString(base64Part)
		if err == nil {
			decodedString := string(decodedBytes)

			parts := strings.SplitN(decodedString, ":", 2)
			if len(parts) == 2 {
				_this.user = parts[0]
			}
		}
	}

	fmt.Println("[Error][Infrastructure]", time.Now().Format("2006-01-02 15:04:05"), statusCode, "url", _this.uri, "\n\t> for user", _this.user, "with message", message)
}

func (_this *GetHttp) GetJwtClientCredentials(params map[string]string, headers map[string]string) (entity.CacheValueEntity, error) {
	_this.user = headers["Authorization"]
	var uri string
	var queryParts []string

	for key, value := range params {
		queryParts = append(queryParts, key+"="+value)
	}
	if len(queryParts) > 0 {
		uri = _this.uri + "?" + strings.Join(queryParts, "&")
	}

	// Crear la solicitud HTTP
	req, err := http.NewRequest("POST", uri, nil)
	if err != nil {
		_this.printLog(500, err.Error())
		return entity.CacheValueEntity{
			StatusCode: 500,
			Message:    "failed to create HTTP request: " + err.Error(),
		}, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: time.Duration(_this.timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		_this.printLog(500, err.Error())
		return entity.CacheValueEntity{
			StatusCode: 500,
			Message:    "failed to send HTTP request: " + err.Error(),
		}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		_this.printLog(500, err.Error())
		return entity.CacheValueEntity{
			StatusCode: 500,
			Message:    "failed to read response body: " + err.Error(),
		}, err
	}

	// Si el estado no es 200, devolver directamente el error
	if resp.StatusCode != http.StatusOK {
		_this.printLog(resp.StatusCode, string(bodyBytes))
		return entity.CacheValueEntity{
			StatusCode: resp.StatusCode,
			Message:    string(bodyBytes), // Devolver el cuerpo como mensaje
		}, nil
	}

	// Deserializar el JSON del body
	var oauthResponse httpinfracollection.Oauth2ResponseCollection
	if err := json.Unmarshal(bodyBytes, &oauthResponse); err != nil {
		_this.printLog(500, err.Error())
		return entity.CacheValueEntity{
			StatusCode: 500,
			Message:    "failed to decode response body: " + err.Error(),
		}, err
	}

	// Retornar la respuesta correctamente mapeada
	return entity.CacheValueEntity{
		StatusCode:   resp.StatusCode,
		AccessToken:  oauthResponse.AccessToken,
		TokenType:    oauthResponse.TokenType,
		ExpiresIn:    oauthResponse.ExpiresIn,
		RefreshToken: oauthResponse.RefreshToken,
		Message:      "success",
	}, nil
}
