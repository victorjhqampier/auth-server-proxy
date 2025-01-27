package entity

import (
	"time"
)

type LocalCacheEntity struct {
	IsBlocked       bool
	CacheExpireTime time.Time
	StoreValue      CacheValueEntity
}

type CacheValueEntity struct {
	Message      string
	StatusCode   int
	AccessToken  string
	TokenType    string
	ExpiresIn    int
	RefreshToken string
}
