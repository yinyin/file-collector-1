package filecollector1

import (
	"fmt"
	"strings"
)

// CollectorFunc is the collector interface
type CollectorFunc func(collectState *CollectState, setup *CollectSetup) (err error)

// CollectorCallable specify callable interface of collector.
type CollectorCallable interface {
	// RunCollect perform collector operation.
	RunCollect(collectState *CollectState, setup *CollectSetup) (err error)
}

// CollectorDiscover discover collector with given setup.
type CollectorDiscover interface {
	// GetCollectorByType returns a collector callable base on given collector
	// type string. Returns `nil` if given condition cannot reach usable one.
	GetCollectorByType(collectorType string) (collector CollectorCallable)

	// GetCollectorByFilePath returns a collector callable base on given file
	// path string. Returns `nil` if given condition cannot reach usable one.
	GetCollectorByFilePath(filePath string) (collector CollectorCallable)
}

// SimpleCollectorImplementation implements CollectorCallable
// and CollectorDiscover interface.
type SimpleCollectorImplementation struct {
	supportedSuffix   []string
	typeName          string
	collectorCallable CollectorFunc
}

// NewSimpleCollectorImplementation returns a simple collector implementation.
// All string given should in lower-case.
func NewSimpleCollectorImplementation(supportedSuffix []string, typeName string, collectorCallable CollectorFunc) (c *SimpleCollectorImplementation) {
	return &SimpleCollectorImplementation{
		supportedSuffix:   supportedSuffix,
		typeName:          typeName,
		collectorCallable: collectorCallable,
	}
}

// RunCollect perform collector operation.
func (c *SimpleCollectorImplementation) RunCollect(collectState *CollectState, setup *CollectSetup) (err error) {
	return c.collectorCallable(collectState, setup)
}

// GetCollectorByType returns a collector callable base on given collector type
// string. Returns `nil` if given condition cannot reach usable one.
func (c *SimpleCollectorImplementation) GetCollectorByType(collectorType string) (collector CollectorCallable) {
	t := strings.ToLower(collectorType)
	if t == c.typeName {
		return c
	}
	return nil
}

// GetCollectorByFilePath returns a collector callable base on given file path
// string. Returns `nil` if given condition cannot reach usable one.
func (c *SimpleCollectorImplementation) GetCollectorByFilePath(filePath string) (collector CollectorCallable) {
	n := strings.ToLower(filePath)
	for _, s := range c.supportedSuffix {
		if strings.HasSuffix(n, s) {
			return c
		}
	}
	return nil
}

var registeredCollectorDiscover = make([]CollectorDiscover, 0)

// RegisterCollectorDiscover adds given instance of collector discover for use with
func RegisterCollectorDiscover(collectorDiscover CollectorDiscover) {
	registeredCollectorDiscover = append(registeredCollectorDiscover, collectorDiscover)
}

// FindCollectorCallable lookup collector by given type and file path.
func FindCollectorCallable(collectorType string, filePath string) (collector CollectorCallable, err error) {
	for _, d := range registeredCollectorDiscover {
		if collectorType != "" {
			if collector = d.GetCollectorByType(collectorType); nil != collector {
				return collector, nil
			}
		}
		if collector = d.GetCollectorByFilePath(filePath); nil != collector {
			return collector, nil
		}
	}
	return nil, fmt.Errorf("cannot reach collector callable: type=[%s], path=[%s]", collectorType, filePath)
}
