package fileFunc

import (
	"bytes"
	"encoding/csv"
	"os"
)

type CsvWriter struct {
	buf    *bytes.Buffer
	writer *csv.Writer
}

func NewCsvFileWriter(out string) (*CsvWriter, error) {
	f, err := os.Create(out)
	if err != nil {
		return nil, err
	}
	writer := csv.NewWriter(f)
	return &CsvWriter{writer: writer}, nil
}
func NewCsvWriter() *CsvWriter {
	b := &bytes.Buffer{}
	f := csv.NewWriter(b)
	return &CsvWriter{writer: f, buf: b}
}
func (cw *CsvWriter) Write(value []string) (err error) {
	return cw.writer.Write(value)
}
func (cw *CsvWriter) Output() string {
	cw.writer.Flush()
	return cw.buf.String()
}
func (cw *CsvWriter) Flush() {
	cw.writer.Flush()
}

func WriteCsv(file string, rows [][]string) error {
	csv := NewCsvWriter()
	err := csv.writer.WriteAll(rows)
	if err != nil {
		return err
	}
	return WriteFileTrunc(file, csv.Output())
}
func ReadCsv(file string) ([][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	return reader.ReadAll()
}
