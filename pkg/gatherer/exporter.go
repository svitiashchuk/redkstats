package gatherer

import (
	"encoding/csv"
	"strconv"
)

type Exporter interface {
	Export(batch Batch) error
}

type ExporterCSV struct {
	Writer *csv.Writer
}

type ExporterMemory struct {
	data Batch
}

func NewExporterCSV(writer *csv.Writer) *ExporterCSV {
	return &ExporterCSV{Writer: writer}
}

func (exp *ExporterCSV) Export(batch Batch) error {
	defer exp.Writer.Flush()

	for k, v := range batch {
		vs := strconv.FormatInt(v, 10)
		err := exp.Writer.Write([]string{k, vs})
		if err != nil {
			return err
		}
	}

	return nil
}

func NewExporterMemory() *ExporterMemory {
	return &ExporterMemory{data: make(Batch)}
}

func (exp *ExporterMemory) Export(batch Batch) error {
	for k, v := range batch {
		exp.data[k] = v
	}

	return nil
}

func (exp *ExporterMemory) GetExportedData() Batch {
	return exp.data
}
