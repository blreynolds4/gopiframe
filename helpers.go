package main

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/blreynolds4/photopi-api/backup"
	"github.com/palantir/stacktrace"
	"github.com/unrolled/render"
)

const DEFAULT_UPLOAD_TAG_NAME string = "uploadImages"
const DEFAULT_PHOTO_PATH string = "./piphotos"
const DEFAULT_UI_PATH string = "./ui/build"
const DEFAULT_SLIDESHOW_DIR string = "./slideshow"

// AppContext holds application configuration data
type AppContext struct {
	Render    *render.Render
	Version   string
	Env       string
	Port      string
	TagName   string
	PhotoPath string
	UIPath    string
	PhotoSave backup.PhotoBackup
}

// Healthcheck will store information about its name and version
type Healthcheck struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
}

// Status is a custom response object we pass around the system and send back to the customer
// 404: Not found
// 500: Internal Server Error
type Status struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// CreateContextForTestSetup initialises an application context struct
// for testing purposes
func CreateContextForTestSetup() AppContext {
	testVersion := "0.0.0"
	ctx := AppContext{
		Render:    render.New(),
		Version:   testVersion,
		Env:       local,
		Port:      "3001",
		TagName:   DEFAULT_UPLOAD_TAG_NAME,
		PhotoPath: DEFAULT_PHOTO_PATH,
		UIPath:    DEFAULT_UI_PATH,
	}
	return ctx
}

// ParseVersionFile returns the version as a string, parsing and validating a file given the path
func ParseVersionFile(versionPath string) (string, error) {
	dat, err := ioutil.ReadFile(versionPath)
	if err != nil {
		return "", stacktrace.Propagate(err, "error reading version file")
	}
	version := string(dat)
	version = strings.Trim(strings.Trim(version, "\n"), " ")
	// regex pulled from official https://github.com/sindresorhus/semver-regex
	semverRegex := `^v?(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)(?:-[\da-z\-]+(?:\.[\da-z\-]+)*)?(?:\+[\da-z\-]+(?:\.[\da-z\-]+)*)?$`
	match, err := regexp.MatchString(semverRegex, version)
	if err != nil {
		return "", stacktrace.Propagate(err, "error executing regex match")
	}
	if !match {
		return "", stacktrace.NewError("string in VERSION is not a valid version number")
	}
	return version, nil
}
