package ipfs

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/daneroo/timewalker"
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
