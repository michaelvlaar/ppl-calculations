package models

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/sirupsen/logrus"
	"ppl-calculations/domain/export"
)

type ExportData struct {
	ID          string
	Name        string
	ViewUrl     string
	DownloadUrl string
	CreatedAt   string
}

func OverviewFromExports(ex []export.Export) []ExportData {
	exs := make([]ExportData, 0)

	for _, e := range ex {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(e); err != nil {
			logrus.WithError(err).Error("encoding export for view url")
			continue
		}

		var gzipBuffer bytes.Buffer
		gz := gzip.NewWriter(&gzipBuffer)

		_, err := gz.Write(buf.Bytes())
		if err != nil {
			logrus.WithError(err).Error("encoding export for view url")
			continue
		}
		err = gz.Close()
		if err != nil {
			logrus.WithError(err).Error("encoding export for view url")
			continue
		}

		exs = append(exs, ExportData{
			ID:          e.ID.String(),
			Name:        e.Name.String(),
			CreatedAt:   e.CreatedAt.Format("15:04:05 02-01-2006"),
			ViewUrl:     fmt.Sprintf("/view?d=%s", base64.URLEncoding.EncodeToString(gzipBuffer.Bytes())),
			DownloadUrl: fmt.Sprintf("/download?d=%s", base64.URLEncoding.EncodeToString(gzipBuffer.Bytes())),
		})
	}

	return exs
}
