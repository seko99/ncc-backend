package exporter

import (
	"bufio"
	s3storage "code.evixo.ru/ncc/ncc-backend/pkg/storage/s3"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/text/encoding/charmap"
	"os"
	"strings"
	"time"
)

type S3ExportWriter struct {
	Storage    *s3storage.Storage
	RemotePath string
	Quoted     bool
	Separator  string
	Encoding   *charmap.Charmap
}

func (ths S3ExportWriter) quoted(str string) string {
	return fmt.Sprintf("\"%s\"", str)
}

func (ths S3ExportWriter) asBytes(slice []string) []byte {
	return []byte(strings.Join(slice, ths.Separator) + "\n")
}

func (ths S3ExportWriter) Write(data []exporter.ExportData, withHeader ...bool) error {
	if data == nil || len(data) == 0 {
		return fmt.Errorf("empty or nil data")
	}

	tmpFileName := fmt.Sprintf("%s.tmp", uuid.NewString())
	csvFile, err := os.CreateTemp("", tmpFileName)
	defer os.Remove(csvFile.Name())
	w := bufio.NewWriter(csvFile)

	if len(withHeader) > 0 && withHeader[0] {
		_, err = w.Write(ths.asBytes(data[0].Header()))
		if err != nil {
			return fmt.Errorf("can't write CSV header: %w", err)
		}
	}

	for _, d := range data {
		var dataSlice []string
		slice := d.ToSlice()

		encoder := ths.Encoding.NewEncoder()
		for _, s := range slice {
			encodedString, err := encoder.String(s)
			if err != nil {
				encodedString = s
			}

			if ths.Quoted {
				dataSlice = append(dataSlice, fmt.Sprintf("\"%s\"", encodedString))
			} else {
				dataSlice = append(dataSlice, encodedString)
			}
		}

		_, err := w.Write(ths.asBytes(dataSlice))
		if err != nil {
			return fmt.Errorf("can't write CSV record: %w", err)
		}
	}

	w.Flush()
	err = csvFile.Close()
	if err != nil {
		return fmt.Errorf("can't close CSV file: %w", err)
	}

	{
		f, err := os.Open(csvFile.Name())
		defer f.Close()
		if err != nil {
			return fmt.Errorf("can't open tmp file: %w", err)
		}

		remoteName := fmt.Sprintf("%s-%s.csv", data[0].FileName(), time.Now().Format("2006-01-02-15-04-05"))

		err = ths.Storage.PutFile(csvFile.Name(), remoteName)
		if err != nil {
			return fmt.Errorf("can't put object: %w", err)
		}
	}
	return nil

}

func (ths S3ExportWriter) GetErrors(exportTime time.Time, path, errorFileName string, d exporter.ExportData) ([]exporter.ExportData, error) {
	return []exporter.ExportData{}, nil
}

func NewS3Writer(
	storage *s3storage.Storage,
	remotePath string,
	quoted ...bool,
) (S3ExportWriter, error) {
	w := S3ExportWriter{
		Storage:    storage,
		RemotePath: remotePath,
		Separator:  ";",
		Encoding:   charmap.Windows1251,
	}

	if len(quoted) > 0 && quoted[0] {
		w.Quoted = true
	}

	return w, nil
}
