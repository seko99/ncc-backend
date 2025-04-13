package exporter

import (
	"bufio"
	"bytes"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
	"github.com/google/uuid"
	"github.com/secsy/goftp"
	"golang.org/x/text/encoding/charmap"
	"os"
	"strings"
	"time"
)

type FtpExportWriter struct {
	URL          string
	BadsURL      string
	Username     string
	Password     string
	BadsUsername string
	BadsPassword string
	RemotePath   string
	Quoted       bool
	Separator    string
	Encoding     *charmap.Charmap
}

func (ths FtpExportWriter) quoted(str string) string {
	return fmt.Sprintf("\"%s\"", str)
}

func (ths FtpExportWriter) asBytes(slice []string) []byte {
	return []byte(strings.Join(slice, ths.Separator) + "\n")
}

func (ths FtpExportWriter) Write(data []exporter.ExportData, withHeader ...bool) error {
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

		remoteName := fmt.Sprintf("%s.csv", data[0].FileName())
		_ = remoteName

		clientConfig := goftp.Config{
			User:               ths.Username,
			Password:           ths.Password,
			ConnectionsPerHost: 10,
			Timeout:            10 * time.Second,
			Logger:             os.Stderr,
		}

		client, err := goftp.DialConfig(clientConfig, ths.URL)
		defer client.Close()
		if err != nil {
			return fmt.Errorf("can't dial ftp: %w", err)
		}

		err = client.Store(remoteName, f)
		if err != nil {
			return fmt.Errorf("can't upload file %s: %w", remoteName, err)
		}
	}
	return nil

}

func (ths FtpExportWriter) GetErrors(exportTime time.Time, path, errorFileName string, d exporter.ExportData) ([]exporter.ExportData, error) {
	exportErrors := []exporter.ExportData{}

	clientConfig := goftp.Config{
		User:               ths.BadsUsername,
		Password:           ths.BadsPassword,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             os.Stderr,
	}

	client, err := goftp.DialConfig(clientConfig, ths.BadsURL)
	defer client.Close()
	if err != nil {
		return nil, fmt.Errorf("can't dial ftp: %w", err)
	}

	files, err := client.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("can't read dir: %w", err)
	}

	newestTime := exportTime.Add(-10 * time.Minute)
	var lastBad string

	for _, f := range files {
		if strings.HasPrefix(f.Name(), errorFileName) && strings.HasSuffix(f.Name(), ".bad") {
			if f.ModTime().After(newestTime) {
				newestTime = f.ModTime()
				lastBad = f.Name()
			}
		}
	}

	if len(lastBad) > 0 {
		b := new(bytes.Buffer)

		err = client.Retrieve(path+"/"+lastBad, b)
		if err != nil {
			return nil, fmt.Errorf("can't retrieve file: %w", err)
		}

		scanner := bufio.NewScanner(b)
		for scanner.Scan() {
			str := scanner.Text()

			var dataSlice []string
			slice := strings.Split(str, ths.Separator)

			decoder := ths.Encoding.NewDecoder()

			for _, s := range slice {
				decodedString, err := decoder.String(s)
				if err != nil {
					decodedString = s
				}
				if ths.Quoted {
					trim := strings.TrimPrefix(decodedString, "\"")
					trim = strings.TrimSuffix(trim, "\"")
					dataSlice = append(dataSlice, trim)
				} else {
					dataSlice = append(dataSlice, decodedString)
				}
			}

			exportError, err := d.FromSlice(dataSlice)
			if err == nil {
				exportErrors = append(exportErrors, exportError)
			}
		}
	}

	return exportErrors, nil
}

func NewFtpWriter(
	url string,
	username string,
	password string,
	badsURL string,
	badsUsername string,
	badsPassword string,
	remotePath string,
	quoted ...bool,
) (FtpExportWriter, error) {
	w := FtpExportWriter{
		URL:          url,
		Username:     username,
		Password:     password,
		BadsURL:      badsURL,
		BadsUsername: badsUsername,
		BadsPassword: badsPassword,
		RemotePath:   remotePath,
		Separator:    ";",
		Encoding:     charmap.Windows1251,
	}

	if len(quoted) > 0 && quoted[0] {
		w.Quoted = true
	}

	return w, nil
}
