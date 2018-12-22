package tgz

import filecollector "github.com/yinyin/file-collector-1"

// SupportedSuffix is slice of supported file name suffix.
var SupportedSuffix = []string{".tar.gz", ".tgz"}

// TypeName is default name of this collector.
var TypeName = "tgz"

// Collector implements a collector for `.tar.gz` or `.tgz` suffixed files.
func Collector(collectState *filecollector.CollectState, setup *filecollector.CollectSetup) (err error) {
	// TODO: impl
	return nil
}
