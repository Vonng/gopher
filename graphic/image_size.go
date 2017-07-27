package graphic

import (
	"os"
	"image"
	_ "image/png"
	_ "image/jpeg"
)

func GetImgSize(id string) (x, y int, err error) {
	path := "/Users/vonng/go/src/tbimg/img/TB1Mb7kMVXXXXb7XVXXSutbFXXX.jpg"
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return
	}
	sz := img.Bounds().Size()
	return sz.X, sz.Y, nil
}
