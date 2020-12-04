package ipfs

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/daneroo/timewalker"
)

// consts Might go into a per reader/writer config..
const (
	defaultBasePath = "./data/ipfs"
)

// dirFor calculates the directory path
// for ipfs this is not an absolute path
func dirFor(grain timewalker.Duration) string {
	return strings.ToLower(grain.String())
}

// pathFor calculates the file path
func pathFor(grain timewalker.Duration, intvl timewalker.Interval) string {
	dir := dirFor(grain)
	file := fmt.Sprintf("%s.jsonl", intvl.Start.Format(time.RFC3339))
	return path.Join(dir, file)
}

// fileIn return a slice for full paths to the file in the appropriate directory
// TODO(daneroo): filter for any inappropriate file (or subdirs);
//   could use filePath.Walk, but that cannot perform filtering (only skip dir or rest of current)
func filesIn(basePath string, grain timewalker.Duration) ([]string, error) {
	// dir := dirFor(grain)

	//Important Note: ReadDir guarantees order by filename
	// files, err := ioutil.ReadDir(dir)
	// if err != nil {
	// 	return nil, err
	// }

	// map each os.FileInfo to a full path
	var filenames []string // == nil
	// for _, file := range files {
	// 	filename := path.Join(dir, file.Name())
	// 	filenames = append(filenames, filename)
	// }

	return filenames, nil
}
