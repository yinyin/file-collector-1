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
	if err = state.MakeCheckSumFile(operationCfg.ChecksumFilePath); nil != err {
		log.Fatalf("failed on writing checksum file [%s]: %v", operationCfg.ChecksumFilePath, err)
		return
	}
	if !state.Check() {
		log.Fatalf("result is not in success state.")
	} else {
		log.Printf("INFO: complete successfully.")
	}
}
