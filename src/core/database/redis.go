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

func InitializeRedis(config RedisConfiguration) (DatabaseClient, error) {
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

	records := []core.SnomedDescription{}

	// Get the concepts by its index
	ids, err := db.adapter.SInter(ctx, "snomed_description_by_cid:"+conceptId).Result()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		fmt.Println(id)
		// get fields
		recordId, err := db.adapter.HGet(ctx, id, "id").Result()
		if err != nil {
			return nil, err
		}
		int_recordId, err := strconv.Atoi(recordId)
		if err != nil {
			return nil, err
		}
		conceptId, err := db.adapter.HGet(ctx, id, "concept_id").Result()
		if err != nil {
			return nil, err
		}
		int_conceptId, err := strconv.Atoi(conceptId)
		if err != nil {
			return nil, err
		}
		term, err := db.adapter.HGet(ctx, id, "term").Result()
		if err != nil {
			return nil, err
		}
		records = append(records, core.SnomedDescription{
			Id:        int_recordId,
			ConceptId: int_conceptId,
			Term:      term,
		})
	}

	return records, nil
}

func (db *RedisDatabase) PutSnomedDescription(record core.SnomedDescription) error {
	ctx := context.Background()
	log.Printf("PutSnomedDescription: %v\n", record)

	// err := db.adapter.HSet(ctx, strconv.Itoa(record.Id), record.GetMap()).Err()
	err := db.adapter.HSet(ctx, strconv.Itoa(record.Id), record).Err()
	if err != nil {
		return err
	}

	log.Printf("PutSnomedDescription snomed_description_by_cid: %d\n", record.ConceptId)
	err = db.adapter.SAdd(ctx, "snomed_description_by_cid:"+strconv.Itoa(record.ConceptId), strconv.Itoa(record.Id)).Err()
	if err != nil {
		return err
	}

	return nil
}
