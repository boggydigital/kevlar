package kvas

type ReduxFabric struct {
	Aggregates  reduxAggregates
	Transitives reduxTransitives
	Atomics     reduxAtomics
}

func initFabric(rf *ReduxFabric) *ReduxFabric {
	if rf == nil {
		rf = &ReduxFabric{}
	}

	if rf.Aggregates == nil {
		rf.Aggregates = make(reduxAggregates)
	}

	if rf.Transitives == nil {
		rf.Transitives = make(reduxTransitives)
	}

	if rf.Atomics == nil {
		rf.Atomics = make(reduxAtomics)
	}

	return rf
}
