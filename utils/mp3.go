package utils

import (
	"os"

	"github.com/gabriel-vasile/mimetype"
	"github.com/hajimehoshi/go-mp3"
)

func MP3Duration(file string) (float64, error) {
	audiofile, _ := os.Open(file)
	defer audiofile.Close()

	d, err := mp3.NewDecoder(audiofile)
	if err != nil {
		return 0, err
	}
	const sampleSize = 4
	samples := d.Length() / sampleSize
	audioLength := samples / int64(d.SampleRate())
	duration := float64(audioLength)
	return duration, nil
}

func Mime(file string) (mime string, err error) {
	audiofile, _ := os.Open(file)
	defer audiofile.Close()

	mtype, err := mimetype.DetectReader(audiofile)
	if err != nil {
		return "", err
	}

	return mtype.String(), nil
}
