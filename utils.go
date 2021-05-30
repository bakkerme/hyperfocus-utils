package hfutils

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// HFHTTPClient a HTTP client with timeout
type HFHTTPClient interface {
	Get(url string) (*http.Response, []byte, error)
	// Post(url string, body []string, timeout time.Duration)
}

// DefaultHTTPTimeout provides the default HTTP request timeout value
var DefaultHTTPTimeout = 10 * time.Second

// HTTPClient is an implementation of the HFHTTPClient interface for http requets
// with modifications specific to this project
type HTTPClient struct {
	ConnTimeout time.Duration
}

// Get is a simplified Get with a configurable timeout value.
// You can provide the default package timeout by using DefaultHTTPTimeout as the timeout
// value
func (ht *HTTPClient) Get(url string) (*http.Response, []byte, error) {
	tr := &http.Transport{
		IdleConnTimeout: ht.ConnTimeout,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, body, err
}

// OSCommandProvider provides the ability to run commands on the OS level
type OSCommandProvider interface {
	Run(string, ...string) (*[]byte, error)
}

// OSCommand implements OSCommandProvider to provide the ability to run commands on the OS
type OSCommand struct{}

// Run runs a command on the OS
func (osc *OSCommand) Run(name string, arg ...string) (*[]byte, error) {
	out, err := exec.Command(name, arg...).Output()
	return &out, err
}

// DirReaderProvider provides the ability to read file directories on disk
type DirReaderProvider interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filepath string) ([]byte, error)
	GetHomeDirPath() string
}

// DirReader implements DirReaderProvider to provide the ability to read file directories on disk
type DirReader struct{}

// ReadDir reads the contents of a directory
func (dr *DirReader) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

// ReadFile reads the contents of a file
func (dr *DirReader) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// GetHomeDirPath mocks the path to the user's home directory
func (dr *DirReader) GetHomeDirPath() string {
	return os.Getenv("HOME")
}

// EnvReader provides the ability to look up environment variables
type EnvReader interface {
	LookupEnv(key string) (string, bool)
}

// EnvRead implements EnvReader to provide the ability to look up environment variables
type EnvRead struct{}

// LookupEnv retrieves the value of the environment variable named
// by the key. If the variable is present in the environment the
// value (which may be empty) is returned and the boolean is true.
// Otherwise the returned value will be empty and the boolean will
// be false.
func (envr *EnvRead) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}
