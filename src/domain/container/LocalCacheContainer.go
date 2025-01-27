package container

import (
	"auth-server-proxy/src/domain/entity"
	"errors"
	"sync"
	"time"
)

type LocalCacheContainer struct {
	objects      map[string]entity.LocalCacheEntity
	locks        map[string]*sync.Mutex // Bloqueo por clave
	mu           sync.RWMutex           // Para proteger el acceso al mapa
	expireSecond int
}

// Singleton
var (
	localCacheInstance *LocalCacheContainer
	localCacheOnce     sync.Once
)

func NewLocalCacheContainer() *LocalCacheContainer {
	localCacheOnce.Do(func() {
		localCacheInstance = &LocalCacheContainer{
			objects:      make(map[string]entity.LocalCacheEntity),
			locks:        make(map[string]*sync.Mutex),
			expireSecond: 5, // TTL predeterminado de 5 segundos
		}
	})
	return localCacheInstance
}

func (sc *LocalCacheContainer) StartLock(key string) *sync.Mutex {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Si no existe un bloqueo para esta clave, crear uno
	if _, exists := sc.locks[key]; !exists {
		sc.locks[key] = &sync.Mutex{}
	}
	return sc.locks[key]
}

func (sc *LocalCacheContainer) Register(name string, obj entity.CacheValueEntity) error {
	newObj := entity.LocalCacheEntity{
		IsBlocked:       true,
		CacheExpireTime: time.Now().Add(time.Duration(30) * time.Second),
		StoreValue:      obj,
	}

	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Bloquear el objeto y establecer su TTL
	sc.objects[name] = newObj

	return nil
}

func (sc *LocalCacheContainer) Get(name string) (entity.CacheValueEntity, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	obj, exists := sc.objects[name]
	if !exists {
		return entity.CacheValueEntity{}, errors.New("object does not exist")
	}

	// Verificar si el objeto ha expirado
	if time.Now().After(obj.CacheExpireTime) {
		return entity.CacheValueEntity{}, errors.New("object expired")
	}

	return obj.StoreValue, nil
}

func (sc *LocalCacheContainer) Refresh(name string, newObj entity.CacheValueEntity) error {
	currentObj := entity.LocalCacheEntity{
		IsBlocked:       false,
		CacheExpireTime: time.Now().Add(time.Duration(sc.expireSecond) * time.Second),
		StoreValue:      newObj,
	}

	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.objects[name] = currentObj

	return nil
}

func (sc *LocalCacheContainer) Delete(name string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if _, exists := sc.objects[name]; !exists {
		return errors.New("object does not exist")
	}

	delete(sc.objects, name)
	return nil
}

func (sc *LocalCacheContainer) GetBlockedStatus(name string) (bool, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	obj, exists := sc.objects[name]
	if !exists {
		return true, errors.New("object does not exist")
	}

	return obj.IsBlocked, nil
}
