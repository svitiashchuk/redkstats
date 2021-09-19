package gatherer

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

// TDOO
type Batch map[string]string

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
}

func (g *Gatherer) Gather(ctx context.Context, exp Exporter) {
	batch := make(Batch)

	cmd := g.rdb.Scan(ctx, g.opt.Cursor, g.opt.Match, g.opt.FetchCount)
	iter := cmd.Iterator()
	cnt := 0

	for iter.Next(ctx) {
		cnt += 1
		key := iter.Val()
		idleTimeCmd := g.rdb.ObjectIdleTime(ctx, key)
		batch[key] = getIdleTimeStr(idleTimeCmd)

		if cnt >= g.opt.BatchSize {
			err := exp.Export(batch)
			if err != nil {
				log.Fatal(err)
			}

			batch = make(Batch)
			cnt = 0
		}
	}

	if cnt != 0 {
		err := exp.Export(batch)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// todo choose by optional parameter duration type
func getIdleTimeStr(cmd *redis.DurationCmd) string {
	return strconv.FormatInt(cmd.Val().Nanoseconds(), 10)
}
