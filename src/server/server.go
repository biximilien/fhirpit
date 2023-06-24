package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os"

	"dev.bonfhir/fhirpit/src/core"
	"dev.bonfhir/fhirpit/src/server/resources"
	"dev.bonfhir/fhirpit/src/server/resources/code_system"
)

var records []core.SnomedDescription

var help = flag.Bool("help", false, "Show help")
var importFlag = false

func main() {

	flag.BoolVar(&importFlag, "import", false, "Import the data")
	flag.Parse()

	// show usage information
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	var buf bytes.Buffer
	logger := log.New(&buf, "logger: ", log.Lshortfile)
	logger.SetOutput(os.Stdout)

	configuration := core.AerospikeConfiguration{
		Host:   "localhost",
		Port:   3000,
		Logger: logger,
	}
	configuration.DefaultPolicy.SocketTimeout = 10000
	configuration.DefaultPolicy.MaxRetries = 3
	configuration.DefaultPolicy.SleepBetweenRetries = 1000
	core.InitializeAerospike(configuration)

	// import the SNOMED data
	if importFlag {
		log.Println("Importing data...")

		records = core.ReadTextFile("./data/sct2_Description_Full-en_US1000124_20230301.txt")

		for _, record := range records {
			core.PutSnomedDescription(record)
		}

	} else {

		code_system_h := &resources.CodeSystem{}
		code_system_lookup := &code_system.CodeSystemLookupOperation{}

		// start a web server
		log.Println("Starting server...")
		mux := http.NewServeMux()
		mux.HandleFunc("/Patient", resources.PatientHandler)
		mux.HandleFunc("/CodeSystem", code_system_h.CodeSystemHandler)
		mux.HandleFunc("/CodeSystem/$lookup", code_system_lookup.CodeSystemLookupOperationHandler)
		http.ListenAndServe(":8080", mux)
	}
}
