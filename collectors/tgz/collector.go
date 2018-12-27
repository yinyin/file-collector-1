package tgz

import (
	"compress/gzip"
	"os"
	"path/filepath"

	filecollector "github.com/yinyin/file-collector-1"
	tarcollector "github.com/yinyin/file-collector-1/collectors/tar"
)

// SupportedSuffix is slice of supported file name suffix.
var SupportedSuffix = []string{".tar.gz", ".tgz"}

// TypeName is default name of this collector.
var TypeName = "tgz"

// Collector implements a collector for `.tar.gz` or `.tgz` suffixed files.
func Collector(collectState *filecollector.CollectState, setup *filecollector.CollectSetup, sourceFolderPath string) (err error) {
	path := filepath.Join(sourceFolderPath, setup.FilePath)
	fp, err := os.Open(path)
	if nil != err {
		return err
	}
	defer fp.Close()
	gzfp, err := gzip.NewReader(fp)
	if nil != err {
		return err
	}
	return tarcollector.CollectViaReader(collectState, setup, gzfp, path)
}

var collectorDiscoverInstance filecollector.CollectorDiscover

// DefaultTGZCollectorDiscover returns default instance of collector discovery routine
func DefaultTGZCollectorDiscover() (collectorDiscover filecollector.CollectorDiscover) {
	if nil == collectorDiscoverInstance {
		collectorDiscoverInstance = filecollector.NewSimpleCollectorImplementation(SupportedSuffix, TypeName, Collector)
	}
	return collectorDiscoverInstance
}
