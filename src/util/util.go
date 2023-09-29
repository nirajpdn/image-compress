package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/h2non/bimg"
)

var AllowedExtensions = map[string]string{
	"jpeg": "jpeg",
	"svg":  "svg",
	"png":  "png",
	"webp": "webp",
	"heif": "heif",
	"avif": "avif",
	"gif":  "gif",
}

var extensionToImageType = map[string]bimg.ImageType{
	"jpeg": bimg.JPEG,
	"png":  bimg.PNG,
	"svg":  bimg.SVG,
	"webp": bimg.WEBP,
	"heif": bimg.HEIF,
	"avif": bimg.AVIF,
	"gif":  bimg.GIF,
}

func ImageProcessing(buffer []byte, quality int, dirname string, extension string) (string, error) {

	filename := strings.Replace(uuid.New().String(), "-", "", -1) + "." + extension

	converted, err := bimg.NewImage(buffer).Convert(extensionToImageType[extension])
	if err != nil {
		return filename, err
	}

	compressed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: quality})
	if err != nil {
		return filename, err
	}

	writeError := bimg.Write(fmt.Sprintf("./"+dirname+"/%s", filename), compressed)
	if writeError != nil {
		return filename, writeError
	}

	return filename, nil
}

func CreateFolder(dirname string) error {
	_, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirname, 0755)
		if errDir != nil {
			return errDir
		}
	}
	return nil
}
