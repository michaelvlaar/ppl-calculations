package adapters

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
)

type ImageService struct{}

func MustNewImageService() *ImageService {
	return &ImageService{}
}

func (s *ImageService) SvgToPng(svg io.Reader) (io.Reader, error) {
	cmd := exec.Command("rsvg-convert", "-f", "png")

	cmd.Stdin = svg

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"stdout": out.String(),
			"stderr": stderr.String(),
		}).Error("executing rsvg-convert")
		return nil, err
	}

	// Geef de PNG-output terug als io.Reader
	return &out, nil
}
