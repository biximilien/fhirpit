package database

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
)

func setupRedis() {
	var buf bytes.Buffer
	logger := log.New(&buf, "logger: ", log.Lshortfile)
	logger.SetOutput(os.Stdout)
	configuration := RedisConfiguration{
		Host:   "localhost",
		Port:   6379,
		Logger: logger,
	}
	InitializeRedis(configuration)
}

func BenchmarkRedisGetSnomedDescription(b *testing.B) {

	fmt.Println("Setup")
	setupRedis() // this require a lot of time
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
