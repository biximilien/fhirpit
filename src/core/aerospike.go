package core

import (
	"log"
	"strconv"
	"time"

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
}

var (
	// Aerospike client instance
	client *aero.Client
	// Base policy for aeropike
	basePolicy *aero.BasePolicy
	// Batch operation policy
	batchPolicy *aero.BatchPolicy
	// Write policy we use for writing
	writePolicy *aero.WritePolicy

	clientPolicy *aero.ClientPolicy
)

func InitializeAerospike(config AerospikeConfiguration) {
	// Simple singleton
	if client != nil {
		log.Fatal("Aerospike client already initialized")
	}
	// Tell aerospike database to use json tags for serialization
	aero.SetAerospikeTag("json")
	// Create a new aerospike client
	lclient, err := aero.NewClient(config.Host, config.Port)
	if err != nil {
		log.Fatal(err)
	}
	// Store our client in the singleton instance
	client = lclient
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
	clientPolicy.ConnectionQueueSize = 2
	clientPolicy.LimitConnectionsToQueueSize = true

	client.DefaultPolicy = basePolicy
	client.DefaultBatchPolicy = batchPolicy
	client.DefaultWritePolicy = writePolicy
	client.DefaultBatchPolicy.GetBasePolicy().ExitFastOnExhaustedConnectionPool = true
	client.DefaultPolicy.GetBasePolicy().ExitFastOnExhaustedConnectionPool = true
	client.DefaultWritePolicy.GetBasePolicy().ExitFastOnExhaustedConnectionPool = true

	asl.Logger.SetLogger(config.Logger)
	asl.Logger.SetLevel(asl.DEBUG)
}

func PutSnomedDescription(record SnomedDescription) error {
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

	return client.PutBins(writePolicy, key, effectiveTime, active, moduleId, conceptId, languageCode, typeId, term, caseSignificanceId)
}

func GetSnomedDescription(conceptId string) []SnomedDescription {

	statement := aero.NewStatement("terminology", "snomed_description", "cid", "term")
	statement.Filter = aero.NewEqualFilter("cid", conceptId)

	queryPolicy := aero.NewQueryPolicy()
	queryPolicy.MaxRecords = 20

	descriptions := []SnomedDescription{}

	recset, err := client.Query(queryPolicy, statement)
	if err != nil {
		log.Fatal(err)
	}

	for records := range recset.Results() {
		if records.Err != nil {
			log.Fatal(records.Err.Error())
		}
		if records != nil {
			id, err := strconv.Atoi(records.Record.Key.Value().String())
			if err != nil {
				log.Fatal(err)
			}
			efftime, err := strconv.Atoi(records.Record.Bins["efftime"].(string))
			if err != nil {
				log.Fatal(err)
			}
			active := records.Record.Bins["active"].(string) == "1"
			modid, err := strconv.ParseUint(records.Record.Bins["modid"].(string), 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			cid, err := strconv.Atoi(records.Record.Bins["cid"].(string))
			if err != nil {
				log.Fatal(err)
			}
			lcode := records.Record.Bins["lcode"].(string)
			typeid, err := strconv.ParseUint(records.Record.Bins["typeid"].(string), 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			term := records.Record.Bins["term"].(string)
			csid, err := strconv.ParseUint(records.Record.Bins["csid"].(string), 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			descriptions = append(descriptions, SnomedDescription{
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

	return descriptions
}
