package jsonl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/daneroo/timewalker"
)

// consts Might go into a per reader/writer config..
const (
	defaultBasePath = "./data/jsonl"
)

// dirFor calculates the directory path
func dirFor(basePath string, grain timewalker.Duration) string {
	return path.Join(basePath, strings.ToLower(grain.String()))
}

// pathFor calculates the file path (and also make any required directories)
func pathFor(basePath string, grain timewalker.Duration, intvl timewalker.Interval) (string, error) {
	dir := dirFor(basePath, grain)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	file := fmt.Sprintf("%s.jsonl", intvl.Start.Format(time.RFC3339))
	return path.Join(dir, file), nil
}

// fileIn return a slice for full paths to the file in the appropriate directory
// TODO(daneroo): filter for any inappropriate file (or subdirs);
//   could use filePath.Walk, but that cannot perform filtering (only skip dir or rest of current)
func filesIn(basePath string, grain timewalker.Duration) ([]string, error) {
	dir := dirFor(basePath, grain)

	//Important Note: ReadDir guarantees order by filename
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// map each os.FileInfo to a full path
	var filenames []string // == nil
	for _, file := range files {
		filename := path.Join(dir, file.Name())
		filenames = append(filenames, filename)
	}

	return filenames, nil
}
