package format

import "fmt"

type Format string

const (
	JSONFormat  = "json"
	TableFormat = "table"
)

func (f *Format) String() string {
	return string(*f)
}

func (f *Format) Set(s string) error {
	switch s {
	case JSONFormat, TableFormat:
		*f = Format(s)
		return nil
	default:
		return fmt.Errorf("unknown format: %s", s)
	}
}

func (f *Format) Type() string {
	return "Format"
}
