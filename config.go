package filecollector1

import (
	"encoding/json"
	"io/ioutil"
)

// CollectDest represents collect destination
type CollectDest struct {
	FromPath string `json:"from"`
	ToPath   string `json:"to"`
}

// CollectSetup represents source of file collection
type CollectSetup struct {
	FilePath          string
	Destinations      []*CollectDest
	collectorCallable CollectorCallable
}

// RunCollect invoke collector callable to do collecting operation on given
// setup.
func (x *CollectSetup) RunCollect(collectState *CollectState, sourceFolderPath string) (err error) {
	return x.collectorCallable.RunCollect(collectState, x, sourceFolderPath)
}

// FindDest searchs for given file path in destination configurations.
func (x *CollectSetup) FindDest(filePath string) (dest *CollectDest) {
	filePath = CleanupPathPrefix(filePath)
	if "" == filePath {
		return nil
	}
	for _, d := range x.Destinations {
		if filePath == d.FromPath {
			return d
		}
	}
	return nil
}

// UnmarshalJSON implements Unmarshaler interface
func (x *CollectSetup) UnmarshalJSON(b []byte) (err error) {
	var d struct {
		FilePath      string         `json:"source"`
		Destinations  []*CollectDest `json:"dests"`
		CollectorType string         `json:"type"`
	}
	if err = json.Unmarshal(b, &d); nil != err {
		return err
	}
	x.FilePath = d.FilePath
	x.Destinations = d.Destinations
	if x.collectorCallable, err = FindCollectorCallable(d.CollectorType, d.FilePath); nil != err {
		return err
	}
	return nil
}

// CollectOperation shows setup of file collecting operation
type CollectOperation struct {
	ChecksumFilePath string          `json:"checksum"`
	CollectSetups    []*CollectSetup `json:"sources"`
}

// LoadCollectOperationConfiguration loads collecting operation configuration from given file path.
func LoadCollectOperationConfiguration(filePath string) (operation *CollectOperation, err error) {
	content, err := ioutil.ReadFile(filePath)
	if nil != err {
		return nil, err
	}
	var t CollectOperation
	if err = json.Unmarshal(content, &t); nil != err {
		return nil, err
	}
	return &t, nil
}
