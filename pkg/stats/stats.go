package stats

import (
	"redkstats/pkg/gatherer"
	"sort"
	"strconv"
)

type Info struct {
	Prefix                string
	Values                []int
	Avg                   float64
	Median                float64
	NinetyFifthPercentile int
}

type PrefixStatsInfoMap []Info

func BatchToPrefixStatsInfoMap(batch gatherer.Batch, prefixLen int) (PrefixStatsInfoMap, error) {
	m, err := GroupBatchByPrefix(batch, prefixLen)
	if err != nil {
		return nil, err
	}

	p := PrefixStatsInfoMap{}
	for k, v := range m {
		s := getStatsInfoFromSlice(v)
		s.Prefix = k

		p = append(p, s)
	}

	return p, nil
}

func getStatsInfoFromSlice(data []int) Info {
	sort.Ints(data)

	return Info{
		Values:                data,
		Avg:                   getAvg(data),
		Median:                getMedian(data),
		NinetyFifthPercentile: getPercentileRank(data, 95),
	}
}

func getPercentileRank(sortedData []int, percent int) int {
	l := len(sortedData)
	rank := percent * l / 100

	return sortedData[rank]
}

func getMedian(sortedData []int) float64 {
	l := len(sortedData)

	if l%2 == 0 {
		return getAvg(sortedData[l/2-1 : l/2+1])
	} else {
		return float64(sortedData[l/2])
	}
}

func getAvg(sortedData []int) float64 {
	var sum int

	for _, v := range sortedData {
		sum += v
	}

	return float64(sum) / float64(len(sortedData))
}

func GroupBatchByPrefix(batch gatherer.Batch, prefixLen int) (map[string][]int, error) {
	var prefix string
	var m = make(map[string][]int)

	for k, v := range batch {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		if prefixLen > len(k) {
			prefix = k[:len(k)-1]
		} else {
			prefix = k[:prefixLen-1]
		}

		m[prefix] = append(m[prefix], i)
	}

	return m, nil
}
