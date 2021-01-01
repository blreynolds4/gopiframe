package main

import (
	"fmt"
	"log"
	"os"

	"github.com/blreynolds4/photopi-api/backup"
	"github.com/blreynolds4/photopi-api/stager"
	"github.com/unrolled/render"
)

const local string = "LOCAL"

func main() {
	var (
		// environment variables
		env        = os.Getenv("ENV")         // LOCAL, DEV, STG, PRD
		port       = os.Getenv("PORT")        // server traffic on this port
		version    = os.Getenv("VERSION")     // path to VERSION file
		tagName    = os.Getenv("UPLOAD_TAG")  // tag files are uploaded in
		photosPath = os.Getenv("PHOTOS_PATH") // get the location to save files
		uiPath     = os.Getenv("UI_PATH")     // get the location of the ui app
		showPath   = os.Getenv("SHOW_PATH")
	)

	if env == "" || env == local {
		// running from localhost, so set some default values
		env = local
		port = "8080"
		version = "VERSION"
		tagName = DEFAULT_UPLOAD_TAG_NAME
		photosPath = DEFAULT_PHOTO_PATH
		uiPath = DEFAULT_UI_PATH
		showPath = DEFAULT_SLIDESHOW_DIR
	}

	// create the photo path if needed
	if _, err := os.Stat(photosPath); os.IsNotExist(err) {
		err := os.Mkdir(photosPath, 0744)
		if err != nil {
			fmt.Println("Unable to create phtoto path ", photosPath)
			os.Exit(1)
		}
	}

	// create the slideshow path if needed
	if _, err := os.Stat(showPath); os.IsNotExist(err) {
		err := os.Mkdir(showPath, 0744)
		if err != nil {
			fmt.Println("Unable to create slideshow path ", showPath)
			os.Exit(1)
		}
	}

	// reading version from file
	version, err := ParseVersionFile(version)
	if err != nil {
		log.Fatal(err)
	}

	// create staging and backup
	stage := stager.NewDirectoryStager(showPath, 25)
	saver := backup.NewAWSBackup(stage, 25)

	// initialse application context
	ctx := AppContext{
		Render:    render.New(),
		Version:   version,
		Env:       env,
		Port:      port,
		TagName:   tagName,
		PhotoPath: photosPath,
		UIPath:    uiPath,
		PhotoSave: saver,
	}

	defer func() {
		stage.Stop()
		fmt.Println("Staged stopped")
		saver.Stop()
		fmt.Println("Saver stopped")
	}()

	// start application
	StartServer(ctx)
}
