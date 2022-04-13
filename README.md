# FileScanner library

This is simple file scanner library that allows to return as a stream (asynchronous) search of files paths from filesystem.

## Usage

```go
import (
	    "github.com/mateuszmierzwinski/filescanner"
		
		"filepath"
		"log"
		"sync"
	)

func main() {
    wg := sync.WaitGroup{}

    f := filescanner.Scanner{}
	resStream,errStream := f.Search("/mnt/drive/", ".png", &wg)

	// Found files handling
	go func() {
	    for {
		    foundFile := <- resStream
			log.Printf("Found file: %s", filepath.Join(foundFile.Path, foundFile.Name))
        }   	
    }()

	// Errors handling
    go func() {
        for {
            errorOccured := <- errStream
            log.Printf("Processing error at %s: %s", errorOccured.Path, errorOccured.Err.Error())
        }   	
    }()

	// this wait group waits for FileScanner to finish searching
    wg.Wait()
}
```