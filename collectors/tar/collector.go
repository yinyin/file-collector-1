package tar

import (
	tarfile "archive/tar"
	"io"
	"log"
	"os"

	filecollector "github.com/yinyin/file-collector-1"
)

// SupportedSuffix is slice of supported file name suffix.
var SupportedSuffix = []string{".tar"}

// TypeName is default name of this collector.
var TypeName = "tar"

// CollectViaReader implements a collector for opened tar stream.
func CollectViaReader(collectState *filecollector.CollectState, setup *filecollector.CollectSetup, reader io.Reader) (err error) {
	fp := tarfile.NewReader(reader)
	for {
		header, err := fp.Next()
		if nil != err {
			if err == io.EOF {
				break
			}
			return err
		}
		log.Printf("file in tar: %v", header.Name)
	}
	return nil
}

// Collector implements a collector for `.tar` suffixed files.
func Collector(collectState *filecollector.CollectState, setup *filecollector.CollectSetup) (err error) {
	fp, err := os.Open(setup.FilePath)
	if nil != err {
		return err
	}
	defer fp.Close()
	return CollectViaReader(collectState, setup, fp)
}

var collectorDiscoverInstance filecollector.CollectorDiscover

// DefaultTarCollectorDiscover returns default instance of collector discovery routine
func DefaultTarCollectorDiscover() (collectorDiscover filecollector.CollectorDiscover) {
	if nil == collectorDiscoverInstance {
		collectorDiscoverInstance = filecollector.NewSimpleCollectorImplementation(SupportedSuffix, TypeName, Collector)
	}
	return collectorDiscoverInstance
}
