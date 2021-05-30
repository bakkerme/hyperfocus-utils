package hfutils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func loadFile(filePath string) ([]byte, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return file, fmt.Errorf("Loading Response File with path %s failed: %s", filePath, err)
	}

	return file, nil
}

// MockHTTPClient provides a mock for the HFHTTPClient interface
type MockHTTPClient struct {
	ThrowError   bool
	StatusCode   int
	Body         []byte
	BodyFilePath string
	Validate     *func(url string)
}

// Get allows for a configurable Get mock
func (ht *MockHTTPClient) Get(url string) (*http.Response, []byte, error) {
	if ht.ThrowError {
		return nil, nil, errors.New("The puppy-girl did the loudest bark")
	}

	// The validation callback allows an implementor to run a validation step
	// on the input of the function
	if ht.Validate != nil {
		validate := *ht.Validate
		validate(url)
	}

	statuscode := 200
	if ht.StatusCode > 0 {
		statuscode = ht.StatusCode
	}

	response := &http.Response{
		StatusCode: statuscode,
	}

	body := ht.Body
	if ht.BodyFilePath != "" {
		file, err := loadFile(ht.BodyFilePath)
		if err != nil {
			return nil, nil, err
		}

		body = file
	}

	return response, body, nil
}

// MockEnvRead provides a mock for the EnvReader interface
type MockEnvRead struct {
	ReturnValueForInput map[string]string
}

// LookupEnv mocks the loading of a Environment variable lookup
func (menvr *MockEnvRead) LookupEnv(key string) (string, bool) {
	value, didFind := menvr.ReturnValueForInput[key]
	return value, didFind
}
