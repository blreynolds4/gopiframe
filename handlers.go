package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/blreynolds4/photopi-api/naming"
)

// HandlerFunc is a custom implementation of the http.HandlerFunc
type HandlerFunc func(http.ResponseWriter, *http.Request, AppContext)

// makeHandler allows us to pass an environment struct to our handlers, without resorting to global
// variables. It accepts an environment (Env) struct and our own handler function. It returns
// a function of the type http.HandlerFunc so can be passed on to the HandlerFunc in main.go.
func makeHandler(ctx AppContext, fn func(http.ResponseWriter, *http.Request, AppContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, ctx)
	}
}

// HealthcheckHandler returns useful info about the app
func HealthcheckHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	check := Healthcheck{
		AppName: "photopi-api",
		Version: ctx.Version,
	}
	ctx.Render.JSON(w, http.StatusOK, check)
}

type postResponse struct {
	Message string   `json:"message"`
	Files   []string `json:"files"`
}

// AddPhotosHandler accepts one or more photos to add to the slideshows
func AddPhotosHandler(w http.ResponseWriter, req *http.Request, ctx AppContext) {
	fmt.Printf("Handling Photos POST request: %+v\n", req)
	result := postResponse{}

	// FormFile returns the first file for the given key in ctx.TagName
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	mpReader, err := req.MultipartReader()
	if err != nil {
		ctx.Render.Text(w, http.StatusInternalServerError, fmt.Sprintf("Error getting mp reader %s", err.Error()))
		return
	}

	for {
		p, err := mpReader.NextPart()
		if err == io.EOF {
			// no more files to save
			break
		}
		if err != nil {
			ctx.Render.Text(w, http.StatusInternalServerError, fmt.Sprintf("Error reading part %s", err.Error()))
			return
		}

		// make sure this part gets closed
		defer p.Close()

		// only save images from the expected form field, skip over the rest
		if p.FormName() == ctx.TagName {
			// read current photo
			fmt.Printf("Uploaded File: %+v from form %s\n", p.FileName(), p.FormName())
			result.Files = append(result.Files, p.FileName())
			data, err := ioutil.ReadAll(p)
			if err != nil {
				result.Message = fmt.Sprintf("Error reading photo %s: %s", p.FileName(), err.Error())
				ctx.Render.JSON(w, http.StatusInternalServerError, result)
				return
			}

			createdPath, err := addFileToPath(ctx.PhotoPath, p.FileName(), data)
			if err != nil {
				result.Message = fmt.Sprintf("Error saving photo %s: %s", p.FileName(), err.Error())
				ctx.Render.JSON(w, http.StatusInternalServerError, result)
				return
			}

			// backup and stage the photo
			ctx.PhotoSave.BackupPhoto(createdPath)

			// create a location header for the added file with the unique filename
			w.Header().Add("Location", newURL(createdPath, req))
		}
	}

	// all good
	fmt.Println("Returning success")
	result.Message = "Successfully uploaded files"
	ctx.Render.JSON(w, http.StatusOK, result)
}

func addFileToPath(rootDir, filename string, data []byte) (string, error) {
	// need to create a unique filename for our new file, starting with what we
	// have and adding numeric extentions until it doesn't exist
	namer := naming.NewExifImageNamer()
	exifName, err := namer.NameImage(data)
	if err != nil {
		return "", err
	}

	if "" == exifName {
		exifName = useUploadTime()
	}

	fqFilename := naming.UniqueFileName(rootDir, exifName, filepath.Ext(filename))

	// write the new file
	err = ioutil.WriteFile(fqFilename, data, 0644)
	if err != nil {
		return "", err
	}

	// return the name we used
	return fqFilename, nil
}

func newURL(file string, req *http.Request) string {
	baseFileame := filepath.Base(file)
	newUrlPath := filepath.Join(req.URL.Path, baseFileame)
	newUrl := &url.URL{
		Scheme: "http",
		Host:   req.Host,
		Path:   newUrlPath,
	}

	fmt.Println("New URL", newUrl.String())
	return newUrl.String()
}

func useUploadTime() string {
	now := time.Now()
	return now.Format("2006-01-02-15-04-05")
}
