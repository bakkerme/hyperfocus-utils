package hfutils

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

// MockFileInfo provides a mock for the FileInfo interface
type MockFileInfo struct {
	IName  string
	ISize  int64
	IIsDir bool
}

// Name provides the name
func (f MockFileInfo) Name() string {
	return f.IName
}

// Size provides the size
func (f MockFileInfo) Size() int64 {
	return f.ISize
}

// Mode provides the unix file permissions
func (f MockFileInfo) Mode() os.FileMode {
	return 0
}

// ModTime provides the time of modification
func (f MockFileInfo) ModTime() time.Time {
	return time.Now()
}

// IsDir is a bool that specifies whether or not it's a directory
func (f MockFileInfo) IsDir() bool {
	return f.IIsDir
}

// Sys represents the underlying data source
func (f MockFileInfo) Sys() interface{} {
	return nil
}

// MockDirReader provides a mock DirReaderProvider with
// a number of properties available to tailor the output of
// filesystem operations supported through the DirReader
type MockDirReader struct {
	ExpectedDirname            *string
	ExpectedFilename           *string
	T                          *testing.T
	ShouldErrorReadDir         bool
	ShouldErrorReadFile        bool
	ReturnReadDirValue         *[]os.FileInfo
	ReturnReadFileValue        *[]byte
	ReturnReadFileValueForPath map[string][]byte
	ReturnHomeDirPath          *string
}

// ReadDir mocks the DirReader ReadDir function
func (mdr *MockDirReader) ReadDir(dirname string) ([]os.FileInfo, error) {
	if mdr.ShouldErrorReadDir {
		return nil, errors.New("oh the humanity")
	}

	if mdr.ExpectedDirname != nil {
		if mdr.T == nil {
			panic("Please provide MockDirReader a T")
		}

		if dirname != *mdr.ExpectedDirname {
			mdr.T.Errorf("dirname provided to ReadDir was not the expected value. Expected %s, got %s", *mdr.ExpectedDirname, dirname)
		}
	}

	if mdr.ReturnReadDirValue != nil {
		return *mdr.ReturnReadDirValue, nil
	}

	return nil, nil
}

// ReadFile mocks the DirReader ReadFile function
func (mdr *MockDirReader) ReadFile(path string) ([]byte, error) {
	if mdr.ShouldErrorReadFile {
		return nil, errors.New("oh the humanity")
	}

	if mdr.ExpectedFilename != nil {
		if mdr.T == nil {
			panic("Please provide MockDirReader a T")
		}

		if path != *mdr.ExpectedFilename {
			mdr.T.Errorf("path provided to ReadFile was not the expected value. Expected %s, got %s", *mdr.ExpectedFilename, path)
		}
	}

	if mdr.ReturnReadFileValue != nil {
		return *mdr.ReturnReadFileValue, nil
	}

	if mdr.ReturnReadFileValueForPath[path] != nil {
		return mdr.ReturnReadFileValueForPath[path], nil
	}

	return nil, nil
}

// GetHomeDirPath mocks the path to the user's home directory
func (mdr *MockDirReader) GetHomeDirPath() string {
	if mdr.ReturnHomeDirPath == nil {
		mdr.T.Error("Attempted to run GetHomeDirPath but no ReturnHomeDirPath was provided")
	}

	return *mdr.ReturnHomeDirPath
}

// MismatchError returns a readable error message that can be used when two values do not match
func MismatchError(functionName string, expected interface{}, got interface{}) string {
	return fmt.Sprintf("%s did not return expected results.\nExpected\n%+v\ngot\n%+v", functionName, expected, got)
}

// ExpectedError returns an error message that can be used when a function should have returned an error
func ExpectedError(functionName string) string {
	return fmt.Sprintf("%s should have returned an error", functionName)
}

// UnexpectedError returns an error message that can be used when a function should have returned an error
func UnexpectedError(functionName string, err error) string {
	return fmt.Sprintf("%s returned unexpected error %s", functionName, err)
}
