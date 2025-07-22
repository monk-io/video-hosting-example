package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient(uri string) *RedisClient {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		// Fallback to simple connection
		opt = &redis.Options{
			Addr: "localhost:6379",
		}
	}

	client := redis.NewClient(opt)

	return &RedisClient{
		client: client,
	}
}

// GetClient returns the Redis client
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// HealthCheck performs a health check on Redis
func (r *RedisClient) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Ping(ctx).Err()
}

// Enqueue adds a job to the queue
func (r *RedisClient) Enqueue(queueName string, job interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if job is already a byte slice (pre-marshaled JSON)
	var jobData []byte
	var err error

	if data, ok := job.([]byte); ok {
		// Already marshaled, use as-is
		jobData = data
	} else {
		// Marshal the object to JSON
		jobData, err = json.Marshal(job)
		if err != nil {
			return err
		}
	}

	return r.client.LPush(ctx, queueName, jobData).Err()
}

// Dequeue removes and returns a job from the queue
func (r *RedisClient) Dequeue(queueName string, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := r.client.BRPop(ctx, timeout, queueName).Result()
	if err != nil {
		return nil, err
	}

	if len(result) != 2 {
		return nil, redis.Nil
	}

	return []byte(result[1]), nil
}

// Set stores a key-value pair
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (r *RedisClient) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Get(ctx, key).Result()
}

// Delete removes a key
func (r *RedisClient) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Del(ctx, key).Err()
}
