package stats

import (
	"redkstats/pkg/gatherer"
	"sort"
	"strings"
)

const DefaultPrefixLen = 8

type Info struct {
	Prefix                string
	Values                []int64
	Avg                   float64
	Median                float64
	NinetyFifthPercentile int64
}

type PrefixStatsInfo []Info

type Stats struct {
	opt *Options
}

func NewStats(opt *Options) *Stats {
	return &Stats{opt: opt}
}

type Options struct {
	PrefixLen          int
	PrefixNamespaceSep string
}

func (st *Stats) BatchToPrefixStatsInfoMap(batch gatherer.Batch) (PrefixStatsInfo, error) {
	m, err := st.GroupBatchByPrefix(batch)
	if err != nil {
		return nil, err
	}

	p := PrefixStatsInfo{}
	for k, v := range m {
		s := getStatsInfoFromSlice(v)
		s.Prefix = k

		p = append(p, s)
	}

	return p, nil
}

func getStatsInfoFromSlice(data []int64) Info {
	sort.Slice(data, func(i, j int) bool { return data[i] < data[j] })

	return Info{
		Values:                data,
		Avg:                   getAvg(data),
		Median:                getMedian(data),
		NinetyFifthPercentile: getPercentileRank(data, 95),
	}
}

func getPercentileRank(sortedData []int64, percent int) int64 {
	l := len(sortedData)
	rank := percent * l / 100

	return sortedData[rank]
}

func getMedian(sortedData []int64) float64 {
	l := len(sortedData)

	if l%2 == 0 {
		return getAvg(sortedData[l/2-1 : l/2+1])
	} else {
		return float64(sortedData[l/2])
	}
}

func getAvg(sortedData []int64) float64 {
	var sum int64

	for _, v := range sortedData {
		sum += v
	}

	return float64(sum) / float64(len(sortedData))
}

func (st *Stats) GroupBatchByPrefix(batch gatherer.Batch) (map[string][]int64, error) {
	var prefix string
	var m = make(map[string][]int64)

	for k, v := range batch {
		prefixLen := st.getPrefixLen(k)
		prefix = getPrefix(k, prefixLen)
		m[prefix] = append(m[prefix], v)
	}

	return m, nil
}

func (st *Stats) getPrefixLen(k string) (prefixLen int) {
	if st.opt.PrefixNamespaceSep != "" {
		prefixLen = strings.Index(k, st.opt.PrefixNamespaceSep)

		if prefixLen != -1 {
			return prefixLen
		}
	}

	if st.opt.PrefixLen > 0 {
		return st.opt.PrefixLen
	}

	return DefaultPrefixLen
}

func getPrefix(k string, prefixLen int) string {
	if prefixLen > len(k) {
		return k
	} else {
		return k[:prefixLen]
	}
}
