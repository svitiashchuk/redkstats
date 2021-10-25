package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"redkstats/pkg/gatherer"
	"redkstats/pkg/stats"
	"redkstats/pkg/utils"
	"sort"
	"time"
)

var (
	batchSize        = 100
	fetchCount int64 = 100
	filename         = "keys.csv"
	ctx              = context.TODO()
)

func main() {
	rdb := setupRedis()

	exp := gatherer.NewExporterMemory()
	gat := gatherer.NewGatherer(rdb, &gatherer.Options{
		FetchCount: fetchCount,
		BatchSize:  batchSize,
		Exporter:   exp,
	})

	gat.Gather(ctx)
	b := exp.GetExportedData()

	for k, v := range b {
		fmt.Printf("%d: %+v\n", k, v)
	}

	st := stats.NewStats(&stats.Options{
		PrefixLen:          5,
		PrefixNamespaceSep: ":",
	})

	p, err := st.BatchToPrefixStatsInfoMap(b)
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(stats.ByAvg(p))

	for k, v := range p {
		fmt.Printf("%d: %s = %+v\n", k, v.Prefix, v.Avg*float64(time.Second)/float64(time.Minute))
	}
}

// todo pass parameters
func setupRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func getResultFile() *os.File {
	file, err := redisflu.GetOrCreateFileToWrite(filename)
	if err != nil {
		log.Fatal(err)
	}

	return file
}
