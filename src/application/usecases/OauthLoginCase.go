package usecases

import (
	"auth-server-proxy/src/application/adapters"
	"auth-server-proxy/src/domain/container"
	"auth-server-proxy/src/domain/entity"
	apiManagerinfraquery "auth-server-proxy/src/infrastructure/apiManagerInfrastructure/query"
	"errors"
	"time"
)

type OauthLoginCase struct {
	oauthApi   apiManagerinfraquery.GetHttp
	localCache *container.LocalCacheContainer
	path       string
}

func NewOauthLoginCase(path string) OauthLoginCase {
	return OauthLoginCase{
		path:       path,
		oauthApi:   apiManagerinfraquery.NewGetHttp(path),
		localCache: container.NewLocalCacheContainer(),
	}
}
func (_this *OauthLoginCase) Start(params map[string]string, headers map[string]string) (adapters.Response, error) {
	const maxRetries = 5
	cacheKey := headers["Authorization"]

	// Obtener el bloqueo asociado a esta clave
	cacheLock := _this.localCache.StartLock(cacheKey)

	for i := 0; i < maxRetries; i++ {
		cache, err := _this.localCache.Get(cacheKey)

		if err != nil {
			// Bloquear el acceso exclusivo para esta clave. + Doble verificación dentro del bloqueo
			cacheLock.Lock()
			cache, err = _this.localCache.Get(cacheKey)
			if err == nil {
				cacheLock.Unlock()
				continue
			}

			_this.localCache.Register(cacheKey, entity.CacheValueEntity{})

			// Solicitar un nuevo token al servidor principal
			res, err := _this.oauthApi.GetJwtClientCredentials(params, headers)
			if err != nil {
				// _this.localCache.Delete(cacheKey)
				_this.localCache.Refresh(cacheKey, entity.CacheValueEntity{StatusCode: 500})
				cacheLock.Unlock()
				return adapters.Response{}, err
			}

			_this.localCache.Refresh(cacheKey, res)
			cacheLock.Unlock()

			return adapters.Response{
				StatusCode:  res.StatusCode,
				AccessToken: res.AccessToken,
				TokenType:   res.TokenType,
				ExpiresIn:   res.ExpiresIn,
			}, nil
		}

		// Si el caché está bloqueado, aplicar política de reintentos
		if val, _ := _this.localCache.GetBlockedStatus(cacheKey); val {
			// fmt.Println("reintento ....", i)
			time.Sleep(time.Duration(1000*(i+1)) * time.Millisecond)
			continue
		}

		return adapters.Response{
			StatusCode:  cache.StatusCode,
			AccessToken: cache.AccessToken,
			TokenType:   cache.TokenType,
			ExpiresIn:   cache.ExpiresIn,
		}, nil
	}

	return adapters.Response{}, errors.New("failed to obtain token after max retries")
}
