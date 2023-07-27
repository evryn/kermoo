package Time

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var InitialTime *time.Time

func MustSetInitialTime() {
	filePath := "/tmp/busybox/initiated-at"

	var err error

	if _, err = os.Stat(filePath); err == nil {
		fmt.Println("Getting existence initialization file...")
		InitialTime, err = readAndParseDatetime(filePath)

	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creating new initialization file...")
		InitialTime, err = writeCurrentDatetime(filePath)
	}

	if err != nil {
		panic(fmt.Errorf("unable to set or retreve initial time: %w", err))
	}
}

func writeCurrentDatetime(filePath string) (*time.Time, error) {
	// Get the current datetime with milliseconds
	currentDatetime := time.Now()

	os.MkdirAll(filepath.Dir(filePath), 0700)

	// Write the datetime to the file
	err := ioutil.WriteFile(filePath, []byte(currentDatetime.Format(time.RFC3339Nano)), 0644)
	if err != nil {
		return nil, err
	}

	return &currentDatetime, nil
}

func readAndParseDatetime(filePath string) (*time.Time, error) {
	// Read the datetime from the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse the datetime
	parsedDatetime, err := time.Parse(time.RFC3339Nano, string(data))
	if err != nil {
		return nil, err
	}

	return &parsedDatetime, nil
}
