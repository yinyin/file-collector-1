package main

import "log"

func main() {
	sourceFolderPath, destinationFolderPath, operationCfg, err := parseCommandParam()
	if nil != err {
		log.Fatalf("invalid command option: %v", err)
		return
	}
	log.Printf("Source Folder: %s\n", sourceFolderPath)
	log.Printf("Destination Folder: %s\n", destinationFolderPath)
	log.Printf("Resulted Checksum File: %s\n", operationCfg.ChecksumFilePath)
}
