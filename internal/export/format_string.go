package export

import "fmt"

// String implements the Stringer interface for Format.
func (f Format) String() string {
	return string(f)
}

// Formats returns all supported export formats.
func Formats() []Format {
	return []Format{FormatCSV, FormatTSV, FormatJSON}
}

// FormatNames returns the string names of all supported formats.
func FormatNames() []string {
	formats := Formats()
	names := make([]string, len(formats))
	for i, f := range formats {
		names[i] = f.String()
	}
	return names
}

// MarshalText implements encoding.TextMarshaler.
func (f Format) MarshalText() ([]byte, error) {
	if _, err := ParseFormat(string(f)); err != nil {
		return nil, fmt.Errorf("export.Format.MarshalText: %w", err)
	}
	return []byte(f), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (f *Format) UnmarshalText(data []byte) error {
	parsed, err := ParseFormat(string(data))
	if err != nil {
		return fmt.Errorf("export.Format.UnmarshalText: %w", err)
	}
	*f = parsed
	return nil
}
