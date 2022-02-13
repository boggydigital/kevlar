package kvas

type ReduxFabric struct {
	Aggregates  ReduxAggregates
	Transitives ReduxTransitives
	Atomics     ReduxAtomics
}

func initFabric(rf *ReduxFabric) *ReduxFabric {
	if rf == nil {
		rf = &ReduxFabric{}
	}

	if rf.Aggregates == nil {
		rf.Aggregates = make(ReduxAggregates)
	}

	if rf.Transitives == nil {
		rf.Transitives = make(ReduxTransitives)
	}

	if rf.Atomics == nil {
		rf.Atomics = make(ReduxAtomics)
	}

	return rf
}
