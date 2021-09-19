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
)

var (
	batchSize        = 100
	fetchCount int64 = 100
	filename         = "keys.csv"
	ctx              = context.TODO()
)

func main() {
	rdb := setupRedis()
	gat := gatherer.NewGatherer(rdb, &gatherer.Options{
		FetchCount: fetchCount,
		BatchSize:  batchSize,
	})

	exp := gatherer.NewExporterMemory()
	gat.Gather(ctx, exp)
	b := exp.GetExportedData()

	p, err := stats.BatchToPrefixStatsInfoMap(b, 5)
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(stats.ByAvg(p))

	for k, v := range p {
		fmt.Printf("%d: %+v\n", k, v)
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
