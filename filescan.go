package filescaner

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

/*
FileEntry node describing file path, name and size - applicable in most use cases
*/
type FileEntry struct {
	Path string
	Name string
	Size int64
}

/*
ProcessingError describes error in stream processing including path, that was processed
*/
type ProcessingError struct {
	Path string
	Err  error
}

/*
Scanner returns scanner that seeks trough filesystem and returns files by provided substring
*/
type Scanner struct {
	fileStream chan *FileEntry
	errStream  chan *ProcessingError
}

/*
Search starts file scanning (in separate thread) and returns two channels initialized - FileEntry channel and ProcessingError channel

Inputs are searchPath describing path on the disk to be searched, fileNameSubstr that is a part of name searched (for example .extension) and
wg - WaitGroup used to finish app processing after search processing is done (recursively searched whole tree).
*/
func (f *Scanner) Search(searchPath string, fileNameSubstr string, wg *sync.WaitGroup) (chan *FileEntry, chan *ProcessingError) {
	if wg != nil {
		wg.Add(1)
	}

	// add buffered stream so it can fill channels before pulling (no deadlock)
	f.fileStream, f.errStream = make(chan *FileEntry, 10), make(chan *ProcessingError, 10)

	go f.run(wg, searchPath, fileNameSubstr)

	return f.fileStream, f.errStream
}

/*
run function is internal runner method that works as separate goroutine
*/
func (f *Scanner) run(wg *sync.WaitGroup, path string, substr string) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	// start processing
	scanEntriesStream(path, strings.ToLower(substr), f.fileStream, f.errStream)

	return
}

/*
scanEntriesStream recursively steps through directories and processes
*/
func scanEntriesStream(path, extension string, stream chan *FileEntry, errStream chan *ProcessingError) {
	dir, err := os.ReadDir(path)
	if err != nil {
		errStream <- &ProcessingError{
			Path: path,
			Err:  err,
		}
		return
	}

	for _, dirEntry := range dir {
		name := dirEntry.Name()
		if dirEntry.IsDir() {
			scanEntriesStream(filepath.Join(path, name), extension, stream, errStream)
			continue
		}

		if strings.HasSuffix(strings.ToLower(name), extension) {
			dirEntryInfo, _ := dirEntry.Info()
			stream <- &FileEntry{
				Path: path,
				Name: name,
				Size: dirEntryInfo.Size(),
			}
		}
	}
}
