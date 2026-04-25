package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/trackpoint/internal/diff"
)

// Format represents a supported export format.
type Format string

const (
	FormatCSV  Format = "csv"
	FormatJSON Format = "json"
	FormatTSV  Format = "tsv"
)

// ParseFormat parses a format string into a Format value.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "csv":
		return FormatCSV, nil
	case "json":
		return FormatJSON, nil
	case "tsv":
		return FormatTSV, nil
	default:
		return "", fmt.Errorf("unsupported export format: %q (supported: csv, json, tsv)", s)
	}
}

// Write serialises the diff result in the requested format to w.
func Write(w io.Writer, result diff.Result, format Format) error {
	switch format {
	case FormatCSV:
		return writeDelimited(w, result, ',')
	case FormatTSV:
		return writeDelimited(w, result, '\t')
	case FormatJSON:
		return writeJSON(w, result)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func writeJSON(w io.Writer, result diff.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(result.Changes)
}

func writeDelimited(w io.Writer, result diff.Result, comma rune) error {
	cw := csv.NewWriter(w)
	cw.Comma = comma
	if err := cw.Write([]string{"key", "type", "old_value", "new_value"}); err != nil {
		return err
	}
	for _, c := range result.Changes {
		row := []string{
			c.Key,
			string(c.Type),
			fmt.Sprintf("%v", c.OldValue),
			fmt.Sprintf("%v", c.NewValue),
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
