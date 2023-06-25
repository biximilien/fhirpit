package main

import (
	"bytes"
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
var importFlag = false

var buf bytes.Buffer
var logger *log.Logger

func main() {

	flag.BoolVar(&importFlag, "import", false, "Import the data")
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
			client.PutSnomedDescription(record)
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
	configuration := database.AerospikeConfiguration{
		Host:   "localhost",
		Port:   3000,
		Logger: logger,
	}
	configuration.DefaultPolicy.SocketTimeout = 10000
	configuration.DefaultPolicy.MaxRetries = 3
	configuration.DefaultPolicy.SleepBetweenRetries = 1000
	return database.InitializeAerospike(configuration)
}
