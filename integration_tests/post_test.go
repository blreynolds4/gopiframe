package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// struct we get back from post photos
type postResponse struct {
	Message string   `json:"message"`
	Files   []string `json:"files"`
}

// unmarshall the response
func unmarshallPostPhotos(r *http.Response) (postResponse, error) {
	pr := postResponse{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return pr, err
	}

	err = json.Unmarshal(body, &pr)
	return pr, err
}

// PostFilesTestSuite groups together tests for posting files to the photo service.
type PostFilesTestSuite struct {
	suite.Suite
	POST_PATH  string
	FIELD_NAME string
}

func (s *PostFilesTestSuite) SetupTest() {
	s.POST_PATH = "http://localhost:8080/photos"
	s.FIELD_NAME = "uploadImages"
}

// Test Posting no files
func (s *PostFilesTestSuite) TestPostNoFiles() {
	request, err := s.newfileUploadRequest(s.POST_PATH, "wrongtag", make([][]string, 0))
	require.Nil(s.T(), err)
	require.NotNil(s.T(), request, "Couldn't create upload request")
	resp, err := s.doRequest(request)
	require.Nil(s.T(), err)

	result, err := unmarshallPostPhotos(resp)
	assert.Nil(s.T(), err)

	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	require.NotContains(s.T(), resp.Header, "Location")
	assert.Equal(s.T(), 0, len(result.Files))
}

// Test Posting file from wrong tag
func (s *PostFilesTestSuite) TestPostFromBadTag() {
	request, err := s.newfileUploadRequest(s.POST_PATH, "wrongtag", [][]string{{"photos/image000.jpg", "photos/image000.jpg"}})
	require.Nil(s.T(), err)
	require.NotNil(s.T(), request, "Couldn't create upload request")
	resp, err := s.doRequest(request)
	require.Nil(s.T(), err)

	result, err := unmarshallPostPhotos(resp)
	assert.Nil(s.T(), err)

	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	require.NotContains(s.T(), resp.Header, "Location")
	assert.Equal(s.T(), 0, len(result.Files))
}

// Test Posting one file
func (s *PostFilesTestSuite) TestPostOneFile() {
	s.testPostNFiles(1)
}

// Test Posting 10 files
func (s *PostFilesTestSuite) TestPostTenFiles() {
	s.testPostNFiles(10)
}

// Test Posting 100 files
func (s *PostFilesTestSuite) TestPostHundredFiles() {
	s.testPostNFiles(100)
}

func (s *PostFilesTestSuite) testPostNFiles(postCount int) {
	// build array of files
	toUpload := make([][]string, postCount)
	for i := 0; i < postCount; i++ {
		filename := fmt.Sprintf("photos/image%03d.jpg", i)
		fsname := fmt.Sprintf("photos/image%03d.jpg", i%5)

		toUpload[i] = []string{filename, fsname}
	}

	request, err := s.newfileUploadRequest(s.POST_PATH, s.FIELD_NAME, toUpload)
	require.Nil(s.T(), err)
	require.NotNil(s.T(), request, "Couldn't create upload request")
	resp, err := s.doRequest(request)
	require.Nil(s.T(), err)

	result, err := unmarshallPostPhotos(resp)
	assert.Nil(s.T(), err)

	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	assert.Equal(s.T(), postCount, len(result.Files))
	require.Contains(s.T(), resp.Header, "Location")
	require.Equal(s.T(), len(resp.Header.Values("Location")), postCount)
}

// Creates a new file upload http request with optional extra params
// filenames is a list of string pairs, first is the name of file for request, second is name to load for
// the data to send.
func (s *PostFilesTestSuite) newfileUploadRequest(uri string, formFieldName string, fileNames [][]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, fns := range fileNames {
		fmt.Println("File pair ", fns)
		// add file from fs
		fmt.Println("Opening ", fns[1])
		file, err := os.Open(fns[1])
		if err != nil {
			fmt.Println("Create open file error: ", err)
			return nil, err
		}
		defer file.Close()

		// create the filename for the loaded file
		fmt.Println("Creating form file", fns[0])
		part, err := writer.CreateFormFile(formFieldName, filepath.Base(fns[0]))
		if err != nil {
			fmt.Println("Create form file error: ", err)
			return nil, err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, err
		}

	}

	// all files added, close the writer and create the request
	err := writer.Close()
	if err != nil {
		fmt.Println("writer close error: ", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func (s *PostFilesTestSuite) doRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(request)
}

func (s *PostFilesTestSuite) dumpBody(b io.Reader) string {
	bodyBytes, err := ioutil.ReadAll(b)
	if err != nil {
		return err.Error()
	}

	return string(bodyBytes)
}

// TestPostFiles is the root method to run the test suite
func TestPostFiles(t *testing.T) {
	suite.Run(t, new(PostFilesTestSuite))
}
