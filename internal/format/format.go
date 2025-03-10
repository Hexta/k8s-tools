package format

import "fmt"

type Format string

type Options struct {
	NoHeaders bool
}

const (
	JSONFormat     = "json"
	TableFormat    = "table"
	VerticalFormat = "vertical"
)

func (f *Format) String() string {
	return string(*f)
}

func (f *Format) Set(s string) error {
	switch s {
	case JSONFormat, TableFormat, VerticalFormat:
		*f = Format(s)
		return nil
	default:
		return fmt.Errorf("unknown format: %s", s)
	}
}

func (f *Format) Type() string {
	return "Format"
}

func Apply(format Format, options Options, input *Data) (string, error) {
	switch format {
	case JSONFormat:
		return JSON(input), nil
	case TableFormat:
		return Table(input, options), nil
	case VerticalFormat:
		return Vertical(input), nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}
