package filecollector1

import (
	"crypto/sha256"
	"io"
)

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

// CollectWithReader collects content from given `reader` and save into file
// specified by `destination`.
func (state *CollectState) CollectWithReader(destination *CollectDest, reader io.Reader, sourceFilePath string) (err error) {
	existedState, ok := state.FileStates[destination.ToPath]
	if ok {
		// TODO: check if check-sum is the same one and update state
	} else {
		// TODO: duplicate to target path and add state
	}
	return nil
}
