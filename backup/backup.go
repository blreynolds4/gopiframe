package backup

import (
	"fmt"

	"github.com/blreynolds4/photopi-api/stager"
)

// PhotoBackup is a service that backs up photos
// in a source directory
// if return is ni, the file has been backed up
type PhotoBackup interface {
	BackupPhoto(source string) error
	Stop()
}

type awsBackup struct {
	stager   stager.PhotoStager
	saveChan chan string
}

func NewAWSBackup(stager stager.PhotoStager, bufferSize int) PhotoBackup {
	saver := awsBackup{
		stager:   stager,
		saveChan: make(chan string, bufferSize),
	}

	go func() {
		for {
			// more will be true if the
			source, more := <-saver.saveChan
			if more {
				go saver.backupAndStage(source)
			}
		}
	}()

	return &saver
}

func (a *awsBackup) backupAndStage(source string) {
	// save the photo to aws
	a.awsBackup((source))

	// add the photo to staging
	a.stager.StagePhoto(source)
}

func (a *awsBackup) awsBackup(source string) {
	// save the photo to aws
	fmt.Println("backing up", source)
	fmt.Println("DONE backing up", source)
}

func (a *awsBackup) BackupPhoto(source string) error {
	a.saveChan <- source
	return nil
}

func (a *awsBackup) Stop() {
	close(a.saveChan)
}
