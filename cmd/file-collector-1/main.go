package main

import (
	"log"

	filecollector "github.com/yinyin/file-collector-1"
	tgzimpl "github.com/yinyin/file-collector-1/collectors/tgz"
)

func registerCollectors() {
	filecollector.RegisterCollectorDiscover(tgzimpl.DefaultTGZCollectorDiscover())
}

func main() {
	registerCollectors()
	sourceFolderPath, destinationFolderPath, operationCfg, err := parseCommandParam()
	if nil != err {
		log.Fatalf("invalid command option: %v", err)
		return
	}
	log.Printf("Source Folder: %s\n", sourceFolderPath)
	log.Printf("Destination Folder: %s\n", destinationFolderPath)
	log.Printf("Resulted Checksum File: %s\n", operationCfg.ChecksumFilePath)
	state := filecollector.NewCollectState(destinationFolderPath)
	for setupIndex, collectSetup := range operationCfg.CollectSetups {
		if err := collectSetup.RunCollect(state, sourceFolderPath); nil != err {
			log.Fatalf("collect failed (%d: %s): %v", setupIndex+1, collectSetup.FilePath, err)
			return
		}
	}
	// TODO: save checksum
}
