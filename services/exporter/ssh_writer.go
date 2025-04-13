package exporter

import (
	"bufio"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"context"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
	"time"
)

type SshExportWriter struct {
	URL          string
	Username     string
	RemotePath   string
	clientConfig ssh.ClientConfig
	signer       ssh.Signer
	Quoted       bool
	Separator    string
}

func (s SshExportWriter) quoted(str string) string {
	return fmt.Sprintf("\"%s\"", str)
}

func (s SshExportWriter) asBytes(slice []string) []byte {
	return []byte(strings.Join(slice, s.Separator) + "\n")
}

func (ths SshExportWriter) Write(data []exporter.ExportData, withHeader ...bool) error {
	if data == nil || len(data) == 0 {
		return fmt.Errorf("empty or nil data")
	}

	tmpFileName := fmt.Sprintf("%s.tmp", uuid.NewString())
	csvFile, err := os.CreateTemp("", tmpFileName)
	defer os.Remove(csvFile.Name())
	w := bufio.NewWriter(csvFile)

	_, err = w.Write(ths.asBytes(data[0].Header()))
	if err != nil {
		return fmt.Errorf("can't write CSV header: %w", err)
	}

	for _, d := range data {
		_, err := w.Write(ths.asBytes(d.ToSlice()))
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
		client := scp.NewClient(ths.URL, &ths.clientConfig)
		err = client.Connect()
		defer client.Close()
		if err != nil {
			return fmt.Errorf("can't connect to host: %w", err)
		}
		remoteName := fmt.Sprintf("%s/customers-%d.csv", ths.RemotePath, time.Now().Unix())
		err = client.CopyFile(context.Background(), f, remoteName, "0655")
		if err != nil {
			return fmt.Errorf("can't copy file: %e", err)
		}
	}
	return nil

}

func (ths SshExportWriter) GetErrors(exportTime time.Time, path, errorFileName string, d exporter.ExportData) ([]exporter.ExportData, error) {
	return nil, nil
}

func NewSshWriter(
	url string,
	username string,
	keyFile string,
	remotePath string,
	quoted ...bool,
) (SshExportWriter, error) {
	key, err := os.ReadFile(keyFile)
	if err != nil {
		return SshExportWriter{}, fmt.Errorf("can't read key file: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return SshExportWriter{}, fmt.Errorf("can't parse key: %w", err)
	}

	w := SshExportWriter{
		URL:        url,
		Username:   username,
		RemotePath: remotePath,
		signer:     signer,
		clientConfig: ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
		Separator: ";",
	}

	if len(quoted) > 0 && quoted[0] {
		w.Quoted = true
	}

	return w, nil
}
