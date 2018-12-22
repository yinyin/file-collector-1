package filecollector1

import "crypto/sha256"

// FileState represent state of one result file
type FileState struct {
	FilePath      string
	CheckSum      [sha256.Size]byte
	SourceFiles   []string
	ConflictFiles []string
}

// CollectState represents state of collecting
type CollectState struct {
	FileStates map[string]*FileState
}
