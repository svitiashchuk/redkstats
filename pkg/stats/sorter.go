package stats

type ByAvg PrefixStatsInfo

func (a ByAvg) Len() int           { return len(a) }
func (a ByAvg) Less(i, j int) bool { return a[i].Avg < a[j].Avg }
func (a ByAvg) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByMedian PrefixStatsInfo

func (a ByMedian) Len() int           { return len(a) }
func (a ByMedian) Less(i, j int) bool { return a[i].Median < a[j].Median }
func (a ByMedian) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByLen PrefixStatsInfo

func (a ByLen) Len() int           { return len(a) }
func (a ByLen) Less(i, j int) bool { return len(a[i].Values) < len(a[j].Values) }
func (a ByLen) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByNinetyFifthPercentile PrefixStatsInfo

func (a ByNinetyFifthPercentile) Len() int { return len(a) }
func (a ByNinetyFifthPercentile) Less(i, j int) bool {
	return a[i].NinetyFifthPercentile < a[j].NinetyFifthPercentile
}
func (a ByNinetyFifthPercentile) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type ByPrefix PrefixStatsInfo

func (a ByPrefix) Len() int           { return len(a) }
func (a ByPrefix) Less(i, j int) bool { return a[i].Prefix < a[j].Prefix }
func (a ByPrefix) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
