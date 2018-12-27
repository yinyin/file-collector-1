package tar

import (
	tarfile "archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"

	filecollector "github.com/yinyin/file-collector-1"
)

// SupportedSuffix is slice of supported file name suffix.
var SupportedSuffix = []string{".tar"}

// TypeName is default name of this collector.
var TypeName = "tar"

// CollectViaReader implements a collector for opened tar stream.
func CollectViaReader(collectState *filecollector.CollectState, setup *filecollector.CollectSetup, reader io.Reader, filePath string) (err error) {
	fp := tarfile.NewReader(reader)
	for {
		header, err := fp.Next()
		if nil != err {
			if err == io.EOF {
				break
			}
			return err
		}
		if header.Typeflag != tarfile.TypeReg {
			log.Printf("tar: skip unregular file: %s (%s)", header.Name, filePath)
			continue
		}
		dest := setup.FindDest(header.Name)
		if nil == dest {
			log.Printf("tar: skip file: %s (%s)", header.Name, filePath)
			continue
		}
		log.Printf("tar: extracting %v", header.Name)
		var modeBits os.FileMode
		if 0 != (header.Mode & 0111) {
			modeBits = 0755
		} else {
			modeBits = 0644
		}
		fileRefPath := filePath + ":" + header.Name
		collectState.CollectWithReader(dest, fp, modeBits, header.ModTime, fileRefPath)
	}
	return nil
}

// Collector implements a collector for `.tar` suffixed files.
func Collector(collectState *filecollector.CollectState, setup *filecollector.CollectSetup, sourceFolderPath string) (err error) {
	path := filepath.Join(sourceFolderPath, setup.FilePath)
	log.Printf("INFO: tar: collecting from [%s]", path)
	fp, err := os.Open(path)
	if nil != err {
		return err
	}
	defer fp.Close()
	return CollectViaReader(collectState, setup, fp, path)
}

var collectorDiscoverInstance filecollector.CollectorDiscover

// DefaultTarCollectorDiscover returns default instance of collector discovery routine
func DefaultTarCollectorDiscover() (collectorDiscover filecollector.CollectorDiscover) {
	if nil == collectorDiscoverInstance {
		collectorDiscoverInstance = filecollector.NewSimpleCollectorImplementation(SupportedSuffix, TypeName, Collector)
	}
	return collectorDiscoverInstance
}
