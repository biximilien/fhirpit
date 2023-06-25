package database

import (
	"errors"
	"log"
	"strconv"
	"time"

	"dev.bonfhir/fhirpit/src/core"
	aero "github.com/aerospike/aerospike-client-go/v6"
	asl "github.com/aerospike/aerospike-client-go/v6/logger"
)

type AerospikeConfiguration struct {
	Host          string `json:"Host"`
	Port          int    `json:"Port"`
	DefaultPolicy struct {
		SocketTimeout       int `json:"SocketTimeout"`
		MaxRetries          int `json:"MaxRetries"`
		SleepBetweenRetries int `json:"SleepBetweenRetries"`
	} `json:"DefaultPolicy"`
	Logger *log.Logger `json:"-"`
}

type AerospikeDatabase struct {
	AerospikeConfiguration
	adapter *aero.Client
}

var (
	// Base policy for aeropike
	basePolicy *aero.BasePolicy
	// Batch operation policy
	batchPolicy *aero.BatchPolicy
	// Write policy we use for writing
	writePolicy *aero.WritePolicy

	clientPolicy *aero.ClientPolicy
)

func (db *AerospikeDatabase) GetAdapter() interface{} {
	return db.adapter
}

func (db *AerospikeDatabase) SetAdapter(adapter interface{}) {
	db.adapter = adapter.(*aero.Client)
}

func InitializeAerospike(config AerospikeConfiguration) (DatabaseClient, error) {
	// Simple singleton
	if client != nil {
		return nil, errors.New("database client already initialized")
	}
	client = &AerospikeDatabase{}
	// Tell aerospike database to use json tags for serialization
	aero.SetAerospikeTag("json")
	// Create a new aerospike client
	adapter, err := aero.NewClient(config.Host, config.Port)
	if err != nil {
		return nil, err
	}
	// Store our client in the singleton instance
	client.SetAdapter(adapter)
	// Create our policies
	basePolicy = aero.NewPolicy()
	basePolicy.SocketTimeout = time.Duration(config.DefaultPolicy.SocketTimeout) * time.Millisecond
	basePolicy.MaxRetries = config.DefaultPolicy.MaxRetries
	basePolicy.SleepBetweenRetries = time.Duration(config.DefaultPolicy.SleepBetweenRetries) * time.Millisecond
	writePolicy = aero.NewWritePolicy(0, 0)
	writePolicy.BasePolicy = *basePolicy
	writePolicy.SendKey = true
	batchPolicy = aero.NewBatchPolicy()
	batchPolicy.BasePolicy = *basePolicy
	clientPolicy = aero.NewClientPolicy()
	clientPolicy.ConnectionQueueSize = 100
	clientPolicy.LimitConnectionsToQueueSize = true

	adapter.DefaultPolicy = basePolicy
	adapter.DefaultBatchPolicy = batchPolicy
	adapter.DefaultWritePolicy = writePolicy
	adapter.DefaultBatchPolicy.GetBasePolicy().ExitFastOnExhaustedConnectionPool = true
	adapter.DefaultPolicy.GetBasePolicy().ExitFastOnExhaustedConnectionPool = true
	adapter.DefaultWritePolicy.GetBasePolicy().ExitFastOnExhaustedConnectionPool = true

	asl.Logger.SetLogger(config.Logger)
	asl.Logger.SetLevel(asl.DEBUG)

	return client, nil
}

func (db AerospikeDatabase) PutSnomedDescription(record core.SnomedDescription) error {
	key, err := aero.NewKey("terminology", "snomed_description", record.Id)
	if err != nil {
		return err
	}

	// efftime, active, modid, cid, lcode, typeid, term, csid
	effectiveTime := aero.NewBin("efftime", record.EffectiveTime)
	active := aero.NewBin("active", record.Active)
	moduleId := aero.NewBin("modid", int(record.ModuleId))
	conceptId := aero.NewBin("cid", record.ConceptId)
	languageCode := aero.NewBin("lcode", record.LanguageCode)
	typeId := aero.NewBin("typeid", int(record.TypeId))
	term := aero.NewBin("term", record.Term)
	caseSignificanceId := aero.NewBin("csid", int(record.CaseSignificanceId))

	// INSERT INTO terminology.snomed_description (PK, effectiveTime, active, moduleId, conceptId, languageCode, typeId, term, caseSignificanceId) VALUES (1, "20201212", "1", "asd", "asd", "asd", "asd", "asd", "asd")

	return db.adapter.PutBins(writePolicy, key, effectiveTime, active, moduleId, conceptId, languageCode, typeId, term, caseSignificanceId)
}

func (db AerospikeDatabase) GetSnomedDescription(conceptId string) ([]core.SnomedDescription, error) {

	statement := aero.NewStatement("terminology", "snomed_description", "cid", "term")
	statement.Filter = aero.NewEqualFilter("cid", conceptId)

	queryPolicy := aero.NewQueryPolicy()
	queryPolicy.MaxRecords = 20

	descriptions := []core.SnomedDescription{}

	recset, err := db.adapter.Query(queryPolicy, statement)
	if err != nil {
		return nil, err
	}

	for records := range recset.Results() {
		if records.Err != nil {
			return nil, records.Err
		}
		if records != nil {
			id, err := strconv.Atoi(records.Record.Key.Value().String())
			if err != nil {
				return nil, err
			}
			efftime, err := strconv.Atoi(records.Record.Bins["efftime"].(string))
			if err != nil {
				return nil, err
			}
			active := records.Record.Bins["active"].(string) == "1"
			modid, err := strconv.ParseUint(records.Record.Bins["modid"].(string), 10, 64)
			if err != nil {
				return nil, err
			}
			cid, err := strconv.Atoi(records.Record.Bins["cid"].(string))
			if err != nil {
				return nil, err
			}
			lcode := records.Record.Bins["lcode"].(string)
			typeid, err := strconv.ParseUint(records.Record.Bins["typeid"].(string), 10, 64)
			if err != nil {
				return nil, err
			}
			term := records.Record.Bins["term"].(string)
			csid, err := strconv.ParseUint(records.Record.Bins["csid"].(string), 10, 64)
			if err != nil {
				return nil, err
			}

			descriptions = append(descriptions, core.SnomedDescription{
				Id:                 id,
				EffectiveTime:      efftime,
				Active:             active,
				ModuleId:           modid,
				ConceptId:          cid,
				LanguageCode:       lcode,
				TypeId:             typeid,
				Term:               term,
				CaseSignificanceId: csid,
			})
		}
	}

	return descriptions, nil
}

func (db AerospikeDatabase) Close() {
	// Close the connection
	db.adapter.Close()
}
