package repository

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type CheckerRepository struct {
	db *redis.Client
}

func NewCheckerRepository(db *redis.Client) *CheckerRepository {
	return &CheckerRepository{db}
}

func (r *CheckerRepository) getKey(id string, destination map[string]string) (string, error) {
	bytes, err := json.Marshal(destination)
	if err != nil {
		return "", err
	}
	return "checker:" + id + ":" + string(bytes), nil
}

func (r *CheckerRepository) Get(id string, destination map[string]string) ([][]string, error) {
	key, err := r.getKey(id, destination)
	if err != nil {
		return nil, err
	}
	result := r.db.Get(context.TODO(), key)
	if result.Err() == redis.Nil {
		return nil, nil
	}
	bytes, err := result.Bytes()
	if err != nil {
		return nil, err
	}
	var ret [][]string
	if err := json.Unmarshal(bytes, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *CheckerRepository) Save(id string, destination map[string]string, val [][]string) error {
	return r.SaveWithExpiration(id, destination, val, time.Hour*3)
}

func (r *CheckerRepository) SaveWithExpiration(id string, destination map[string]string, val [][]string, expirationTime time.Duration) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	key, err := r.getKey(id, destination)
	if err != nil {
		return err
	}
	status := r.db.Set(context.TODO(), key, bytes, expirationTime)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}
