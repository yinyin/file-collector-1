package filecollector1

// CollectorFunc is the collector interface
type CollectorFunc func(collectState *CollectState, setup *CollectSetup) (err error)
