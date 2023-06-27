package main

import (
	"bytes"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"

	"dev.bonfhir/fhirpit/src/core"
	"dev.bonfhir/fhirpit/src/core/database"
	"dev.bonfhir/fhirpit/src/server/resources"
	"dev.bonfhir/fhirpit/src/server/resources/code_system"
)

var records []core.SnomedDescription

var help = flag.Bool("help", false, "Show help")
var hostFlag string
var portFlag int
var userFlag string
var passwordFlag string
var dbTypeFlag string
var importFlag = false

var buf bytes.Buffer
var logger *log.Logger

func main() {

	flag.BoolVar(&importFlag, "import", false, "Import the data")
	flag.StringVar(&hostFlag, "host", "localhost", "Database host")
	flag.IntVar(&portFlag, "port", 3000, "Database port")
	flag.StringVar(&userFlag, "username", "", "Database username")
	flag.StringVar(&passwordFlag, "password", "", "Database password")
	flag.StringVar(&dbTypeFlag, "database-type", "aerospike", "Database adapter type")
	flag.Parse()

	// show usage information
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	logger = log.New(&buf, "logger: ", log.Lshortfile)
	logger.SetOutput(os.Stdout)

	client, err := initializeDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// import the SNOMED data
	if importFlag {
		log.Println("Importing data...")

		// TODO: import other files
		records = core.ReadTextFile("./data/sct2_Description_Full-en_US1000124_20230301.txt")

		for _, record := range records {
			err := client.PutSnomedDescription(record)
			if err != nil {
				log.Fatal(err)
			}
		}

	} else {

		code_system_lookup := code_system.NewCodeSystemLookupOperation(client)
		code_system_find_matches := code_system.NewCodeSystemFindMatchesOperation(client)

		// start a web server
		log.Println("Starting server...")
		mux := http.NewServeMux()
		mux.HandleFunc("/Patient", resources.PatientHandler)
		mux.HandleFunc("/CodeSystem", resources.CodeSystemHandler)
		mux.HandleFunc("/CodeSystem/$lookup", code_system_lookup.CodeSystemLookupOperationHandler)
		mux.HandleFunc("/CodeSystem/$find-matches", code_system_find_matches.CodeSystemFindMatchesOperationHandler)
		http.ListenAndServe(":8080", mux)
	}
}

func initializeDatabase() (database.DatabaseClient, error) {
	switch dbTypeFlag {
	case "aerospike":
		return initializeAerospike()
	case "redis":
		return initializeRedis()
	default:
		return nil, errors.New("invalid database type")
	}
}

func initializeAerospike() (database.DatabaseClient, error) {
	configuration := database.AerospikeConfiguration{
		Host:   hostFlag,
		Port:   portFlag,
		Logger: logger,
	}
	configuration.DefaultPolicy.SocketTimeout = 10000
	configuration.DefaultPolicy.MaxRetries = 3
	configuration.DefaultPolicy.SleepBetweenRetries = 1000
	return database.InitializeAerospike(configuration)
}

func initializeRedis() (database.DatabaseClient, error) {
	configuration := database.RedisConfiguration{
		Host:   hostFlag,
		Port:   portFlag,
		Logger: logger,
	}
	return database.InitializeRedis(configuration)
}
