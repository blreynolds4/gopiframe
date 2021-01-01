package stager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/blreynolds4/photopi-api/naming"
)

// Stager is a service that moves a file
// from the provided path to the slide show dir
type PhotoStager interface {
	StagePhoto(source string) error
	Stop()
}

type directoryStager struct {
	stageDir  string
	stageChan chan string
}

func NewDirectoryStager(stageDir string, bufferSize int) PhotoStager {
	stager := directoryStager{
		stageDir:  stageDir,
		stageChan: make(chan string, bufferSize),
	}

	go func() {
		for {
			// more will be true if the
			source, more := <-stager.stageChan
			if more {
				ext := filepath.Ext(source)
				filename := filepath.Base(source)
				// remove the extension
				filename = strings.ReplaceAll(filename, ext, "")
				destination := naming.UniqueFileName(stager.stageDir, filename, ext)
				err := os.Rename(source, destination)
				if err != nil {
					// requeue the file for staging
					fmt.Println("Failed to move ", source, "because", err.Error())
					// failed staging will get requeued by another process
					// doing it here is likely to create an infinite loop until
					// the error is corrected
				}
				fmt.Println("Staged", destination)
			}
		}
	}()

	return &stager
}

func (d *directoryStager) StagePhoto(source string) error {
	d.stageChan <- source
	return nil
}

func (d *directoryStager) Stop() {
	close(d.stageChan)
}
