package filecollector1

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
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

// FileStatesByPath is slice of FileState with sort interface implemented.
type FileStatesByPath []*FileState

func (a FileStatesByPath) Len() int      { return len(a) }
func (a FileStatesByPath) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FileStatesByPath) Less(i, j int) bool {
	return a[i].FilePath < a[j].FilePath
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

// Check check if errors in collecting operation.
func (state *CollectState) Check() (success bool) {
	success = true
	for _, fileState := range state.FileStates {
		if 0 != len(fileState.ConflictFiles) {
			success = false
			log.Printf("WARN: have conflict file: %v - %v", fileState.FilePath, fileState.ConflictFiles)
		}
	}
	return success
}

// MakeCheckSumFile generates checksum file via current file state records.
func (state *CollectState) MakeCheckSumFile(filePath string) (err error) {
	records := make([]*FileState, 0, len(state.FileStates))
	for _, s := range state.FileStates {
		records = append(records, s)
	}
	sort.Sort(FileStatesByPath(records))
	fp, err := os.OpenFile(filepath.Join(state.DestinationFolderPath, filePath), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if nil != err {
		return err
	}
	defer fp.Close()
	for _, s := range records {
		hexChksum := make([]byte, hex.EncodedLen(len(s.CheckSum)))
		hex.Encode(hexChksum, s.CheckSum)
		relPath, err := filepath.Rel(state.DestinationFolderPath, s.FilePath)
		if nil != err {
			return err
		}
		fp.Write(hexChksum)
		fp.WriteString(" *")
		fp.WriteString(relPath)
		fp.WriteString("\n")
	}
	return nil
}
