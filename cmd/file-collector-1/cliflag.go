package main

import (
	"errors"
	"flag"
	"path/filepath"

	filecollector "github.com/yinyin/file-collector-1"
)

// ErrRequireSourceFolderPath indicates path of source folder is not given
var ErrRequireSourceFolderPath = errors.New("source folder is required")

// ErrRequireDestFolderPath indicates path of destination folder is not given
var ErrRequireDestFolderPath = errors.New("destination folder is required")

// ErrRequireCollectOperationConfig indicates path of collecting operation configuration is not given
var ErrRequireCollectOperationConfig = errors.New("operation configuration is required")

func parseCommandParam() (sourceFolderPath, destinationFolderPath string, operationCfg *filecollector.CollectOperation, err error) {
	var configFilePath string
	flag.StringVar(&configFilePath, "conf", "", "path to collect configuration")
	flag.StringVar(&sourceFolderPath, "src", "", "path to source package folder")
	flag.StringVar(&destinationFolderPath, "dest", "", "path to collected file destination folder")
	flag.Parse()
	if "" == sourceFolderPath {
		err = ErrRequireSourceFolderPath
		return
	}
	if sourceFolderPath, err = filepath.Abs(sourceFolderPath); nil != err {
		return
	}
	if "" == destinationFolderPath {
		err = ErrRequireDestFolderPath
		return
	}
	if destinationFolderPath, err = filepath.Abs(destinationFolderPath); nil != err {
		return
	}
	if "" == configFilePath {
		err = ErrRequireCollectOperationConfig
		return
	}
	if operationCfg, err = filecollector.LoadCollectOperationConfiguration(configFilePath); nil != err {
		return
	}
	err = nil
	return
}
