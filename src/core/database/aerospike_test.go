package database

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
)

func setup() {
	var buf bytes.Buffer
	logger := log.New(&buf, "logger: ", log.Lshortfile)
	logger.SetOutput(os.Stdout)
	configuration := AerospikeConfiguration{
		Host:   "localhost",
		Port:   3000,
		Logger: logger,
	}
	configuration.DefaultPolicy.SocketTimeout = 10000
	configuration.DefaultPolicy.MaxRetries = 3
	configuration.DefaultPolicy.SleepBetweenRetries = 1000
	InitializeAerospike(configuration)
}

func BenchmarkGetSnomedDescription(b *testing.B) {

	fmt.Println("Setup")
	setup() // this require a lot of time
	fmt.Println("Start Test")
	b.Run("mytest", func(b *testing.B) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			recs, err := client.GetSnomedDescription("100922009")
			if err != nil {
				log.Fatal(err)
			}
			for _, rec := range recs {
				fmt.Println(rec)
			}
		}
	})
}
