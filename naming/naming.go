package naming

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dsoprea/go-exif"
	log "github.com/dsoprea/go-logging"
)

const DATE_TIME_TAG = "DateTime"

// ImageNamer generates an image filename based on
// metadata in the image
type ImageNamer interface {
	NameImage(image []byte) (string, error)
}

type exifDataNamer struct {
}

func NewExifImageNamer() ImageNamer {
	result := exifDataNamer{}

	return &result
}

func (edn *exifDataNamer) NameImage(data []byte) (string, error) {
	rawExif, err := exif.SearchAndExtractExif(data)
	if err != nil {
		return "", err
	}

	// Run the parse for the exif data
	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	photoTimeStamp := ""
	visitor := func(fqIfdPath string, ifdIndex int, tagId uint16, tagType exif.TagType, valueContext exif.ValueContext) (err error) {
		// We're just looking for the Date Time to use for the filename, so this visitor will
		// exit when it finds tag DateTime and saves it's valueString
		defer func() {
			if state := recover(); state != nil {
				err = log.Wrap(state.(error))
				log.Panic(err)
			}
		}()

		ifdPath, err := im.StripPathPhraseIndices(fqIfdPath)
		log.PanicIf(err)

		it, err := ti.Get(ifdPath, tagId)
		if err != nil {
			if log.Is(err, exif.ErrTagNotFound) {
				fmt.Printf("WARNING: Unknown tag: [%s] (%04x)\n", ifdPath, tagId)
				return nil
			} else {
				log.Panic(err)
			}
		}

		//fmt.Println("Checking tag", it.Name)
		if it.Name == DATE_TIME_TAG {
			valueString := ""
			var value interface{}
			if tagType.Type() == exif.TypeUndefined {
				var err error
				value, err = valueContext.Undefined()
				if err != nil {
					if err == exif.ErrUnhandledUnknownTypedTag {
						value = nil
					} else {
						log.Panic(err)
					}
				}

				photoTimeStamp = fmt.Sprintf("%v", value)
			} else {
				valueString, err = valueContext.FormatFirst()
				log.PanicIf(err)

				photoTimeStamp = valueString
			}

			// found datatime stamp remove :'s and spaces for filename
			photoTimeStamp = strings.ReplaceAll(photoTimeStamp, ":", "-")
			photoTimeStamp = strings.ReplaceAll(photoTimeStamp, " ", "-")
			return nil
		}

		return nil
	}

	_, err = exif.Visit(exif.IfdStandard, im, ti, rawExif, visitor)
	return photoTimeStamp, err
}

func UniqueFileName(root, filename, extension string) string {
	fqFilename := filepath.Join(root, filename+extension)

	check := fqFilename
	uniqExt := 1
	for {

		if _, err := os.Stat(check); os.IsNotExist(err) {
			// found a unique name
			fqFilename = check
			break
		}

		// add an extension to the name and keep trying
		check = filepath.Join(root, fmt.Sprintf("%s_%d", filename, uniqExt)+extension)
		uniqExt = uniqExt + 1
	}

	return fqFilename
}
