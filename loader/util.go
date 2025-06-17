package loader

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

// openCSV returns a csv.Reader given a file path
func openCSV(path string) (*csv.Reader, io.Closer, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, nil, err
    }
    r := csv.NewReader(f)
    r.TrimLeadingSpace = true
    return r, f, nil
}

func parseInt(s string) (int, error) {
    if s == "" {
        return 0, nil
    }
    return strconv.Atoi(s)
} 