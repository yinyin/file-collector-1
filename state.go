package filecollector1

import (
	"bytes"
	"crypto/sha256"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func saveFileContent(destFilePath string, reader io.Reader) (digest []byte, err error) {
	d := filepath.Dir(destFilePath)
	if err = os.MkdirAll(d, 0755); nil != err {
		return nil, err
	}
	fp, err := os.OpenFile(destFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if nil != err {
		return nil, err
	}
	defer fp.Close()
	h := sha256.New()
	writer := io.MultiWriter(fp, h)
	if _, err := io.Copy(writer, reader); nil != err {
		return nil, err
	}
	digest = h.Sum(nil)
	return digest, nil
}

// FileState represent state of one result file
type FileState struct {
	FilePath      string
	CheckSum      []byte
	SourceFiles   []string
	ConflictFiles []string
}

func newFileState(filePath string, checkSum []byte, sourceFilePath string) (state *FileState) {
	return &FileState{
		FilePath: filePath,
		CheckSum: checkSum,
		SourceFiles: []string{
			sourceFilePath,
		},
		ConflictFiles: nil,
	}
}

func (state *FileState) addFileRecord(reader io.Reader, sourceFilePath string) (err error) {
	h := sha256.New()
	if _, err := io.Copy(h, reader); nil != err {
		return err
	}
	digest := h.Sum(nil)
	if bytes.Equal(state.CheckSum, digest) {
		state.SourceFiles = append(state.SourceFiles, sourceFilePath)
	} else {
		state.ConflictFiles = append(state.ConflictFiles, sourceFilePath)
	}
	return nil
}

// CollectState represents state of collecting
type CollectState struct {
	DestinationFolderPath string
	FileStates            map[string]*FileState
}

// NewCollectState creates a new instance of CollectState with given
// destionation folder path.
func NewCollectState(destinationFolderPath string) (state *CollectState) {
	return &CollectState{
		DestinationFolderPath: destinationFolderPath,
		FileStates:            make(map[string]*FileState),
	}
}

// CollectWithReader collects content from given `reader` and save into file
// specified by `destination`.
func (state *CollectState) CollectWithReader(destination *CollectDest, reader io.Reader, modeBits os.FileMode, modifyTime time.Time, sourceFilePath string) (err error) {
	existedState, ok := state.FileStates[destination.ToPath]
	if ok {
		return existedState.addFileRecord(reader, sourceFilePath)
	} else {
		destFilePath := filepath.Join(state.DestinationFolderPath, destination.ToPath)
		digest, err := saveFileContent(destFilePath, reader)
		if nil != err {
			return err
		}
		if err = os.Chmod(destFilePath, modeBits); nil != err {
			log.Printf("ERROR: cannot set mode bits of given file [%s; 0%o]: %v", destFilePath, modeBits, err)
		}
		if err = os.Chtimes(destFilePath, time.Now(), modifyTime); nil != err {
			log.Printf("ERROR: cannot set modify time of given file [%s; %v]: %v", destFilePath, modifyTime, err)
		}
		fileState := newFileState(destFilePath, digest, sourceFilePath)
		state.FileStates[destination.ToPath] = fileState
	}
	return nil
}
