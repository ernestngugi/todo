package controller

import (
	"encoding/json"

	"github.com/ernestngugi/todo/internal/providers"
)

type (
	CacheController interface {
		CacheValue(key string, value any) error
		Exists(key string) (bool, error)
		GetCachedValue(key string, result any) error
		RemoveFromCache(key string) error
	}

	cacheController struct {
		redisProvider providers.Redis
	}
)

func NewCacheController(
	redisProvider providers.Redis,
) CacheController {
	return &cacheController{
		redisProvider: redisProvider,
	}
}

func NewTestCacheController(
	redisProvider providers.Redis,
) *cacheController {
	return &cacheController{
		redisProvider: redisProvider,
	}
}

func (s *cacheController) CacheValue(
	key string,
	value any,
) error {

	cacheData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = s.redisProvider.Set(key, cacheData)
	if err != nil {
		return err
	}

	return nil
}

func (s *cacheController) Exists(
	key string,
) (bool, error) {

	exists, err := s.redisProvider.Exists(key)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *cacheController) GetCachedValue(
	key string,
	result any,
) error {

	payload, err := s.redisProvider.Get(key)
	if err != nil {
		return err
	}

	data, ok := payload.([]byte)
	if !ok {
		return err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	return nil
}

func (s *cacheController) RemoveFromCache(
	key string,
) error {

	err := s.redisProvider.Del(key)
	if err != nil {
		return err
	}

	return nil
}
