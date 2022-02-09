package kvas

//reduxAggregates are individual values that only serve the purpose
//of aggregating other values, e.g. `text` = [`title`, `description`, ...].
//Reductions won't have values of `text`, only `title` and `description`,
//and clients can use `text` when they needs `title` and/or `description`.
//For example when matching terms to `text` - either `title` or `description`
//match would return positive result for `text`.
//Clients only need to establish this aggregation relationship and reduxList
//would take care of replacing aggregate with specific values at runtime.
type reduxAggregates map[string][]string

func (ra reduxAggregates) IsAggregate(key string) bool {
	for a, _ := range ra {
		if a == key {
			return true
		}
	}
	return false
}

func (ra reduxAggregates) Aggregates() []string {
	aggr := make([]string, 0, len(ra))
	for a, _ := range ra {
		aggr = append(aggr, a)
	}
	return aggr
}

func (ra reduxAggregates) Detail(key string) []string {
	return ra[key]
}

func (ra reduxAggregates) DetailAll(keys ...string) map[string]bool {
	details := make(map[string]bool)

	for _, k := range keys {
		if ra.IsAggregate(k) {
			for _, dk := range ra.Detail(k) {
				details[dk] = true
			}
		} else {
			details[k] = true
		}
	}

	return details
}

func (ra reduxAggregates) Aggregate(key string) string {
	for a, details := range ra {
		for _, d := range details {
			if d == key {
				return a
			}
		}
	}
	return ""
}
