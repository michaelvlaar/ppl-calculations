package adapters

import (
	"bytes"
	"io"
	"os/exec"
)

type ImageService struct{}

func MustNewImageService() *ImageService {
	return &ImageService{}
}

func (s *ImageService) SvgToPng(svg io.Reader) (io.Reader, error) {
	cmd := exec.Command("convert", "svg:-", "png:-")
	cmd.Stdin = svg

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return &out, nil
}
