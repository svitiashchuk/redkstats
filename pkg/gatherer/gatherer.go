package gatherer

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

type Batch map[string]int64

type Gatherer struct {
	rdb *redis.Client
	opt *Options
}

func NewGatherer(rdb *redis.Client, opt *Options) *Gatherer {
	return &Gatherer{rdb: rdb, opt: opt}
}

type Options struct {
	Cursor     uint64
	Match      string
	FetchCount int64
	BatchSize  int
	Exporter   Exporter
}

func (g *Gatherer) Gather(ctx context.Context) {
	batch := make(Batch)

	cmd := g.rdb.Scan(ctx, g.opt.Cursor, g.opt.Match, g.opt.FetchCount)
	iter := cmd.Iterator()
	cnt := 0

	for iter.Next(ctx) {
		cnt += 1
		key := iter.Val()
		idleTimeCmd := g.rdb.ObjectIdleTime(ctx, key)
		batch[key] = int64(idleTimeCmd.Val().Seconds())

		if cnt >= g.opt.BatchSize {
			err := g.opt.Exporter.Export(batch)
			if err != nil {
				log.Fatal(err)
			}

			batch = make(Batch)
			cnt = 0
		}
	}

	if cnt != 0 {
		err := g.opt.Exporter.Export(batch)
		if err != nil {
			log.Fatal(err)
		}
	}
}
