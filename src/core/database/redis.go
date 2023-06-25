package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"dev.bonfhir/fhirpit/src/core"
	"github.com/redis/go-redis/v9"
)

type RedisConfiguration struct {
	Host   string      `json:"Host"`
	Port   int         `json:"Port"`
	Logger *log.Logger `json:"-"`
}

type RedisDatabase struct {
	RedisConfiguration
	adapter *redis.Client
}

func (db *RedisDatabase) GetAdapter() interface{} {
	return db.adapter
}

func (db *RedisDatabase) SetAdapter(adapter interface{}) {
	db.adapter = adapter.(*redis.Client)
}

func (db *RedisDatabase) Close() {
	db.adapter.Close()
}

func InitializeRedis(config AerospikeConfiguration) (DatabaseClient, error) {
	// Simple singleton
	if client != nil {
		return nil, errors.New("database client already initialized")
	}

	client = &RedisDatabase{}

	// Create a new aerospike client
	adapter := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	client.SetAdapter(adapter)

	return client, nil
}

func (db *RedisDatabase) GetSnomedDescription(conceptId string) ([]core.SnomedDescription, error) {
	ctx := context.Background()

	val, err := db.adapter.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)

	return nil, nil
}

func (db *RedisDatabase) PutSnomedDescription(record core.SnomedDescription) error {
	ctx := context.Background()

	err := db.adapter.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		return err
	}
	return nil
}
