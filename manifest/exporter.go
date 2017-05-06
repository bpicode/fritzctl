package manifest

import (
	"io"

	"gopkg.in/yaml.v2"
)

// Exporter allows exporting of a Plan.
type Exporter struct {
	writer io.Writer
}

type marshaller func(in interface{}) ([]byte, error)

// ExporterTo creates an Exporter, which writes to the given io.Writer.
func ExporterTo(w io.Writer) *Exporter {
	return &Exporter{writer: w}
}

// Export performs the export.
func (e *Exporter) Export(plan *Plan) error {
	return e.export(plan, yaml.Marshal)
}

func (e *Exporter) export(p *Plan, m marshaller) error {
	bs, err := m(p)
	if err != nil {
		return err
	}
	_, err = e.writer.Write(bs)
	return err
}
