package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/nfnt/resize"
)

const usage = `resizepng resizes png file using nfnt/resize's Lanczos3. A resized file is created a name of "resize_" + png file name.

usage

    resizepng filename
    resizepng filename rate

args

    fileName    png file
    rate        image size rate(%) default 100
`

func main() {
	len := len(os.Args)
	if len != 2 && len != 3 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	filename := os.Args[1]
	if _, err := os.Stat(filename); err != nil {
		log.Fatal(err)
	}

	var rate = 100
	if len == 3 {
		rt, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		rate = rt
	}

	img, err := decode(filename)
	if err != nil {
		log.Fatal(err)
	}

	width, height, err := widthHeight(filename)
	if err != nil {
		log.Fatal(err)
	}

	rimg := resize.Resize(uint(width*rate/100), uint(height*rate/100), img, resize.Lanczos3)

	if err := output(rimg, filename); err != nil {
		log.Fatal(err)
	}
}

func decode(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func widthHeight(filename string) (int, int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	cfg, err := png.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}

	return cfg.Width, cfg.Height, nil
}

func output(img image.Image, filename string) error {
	dir, file := filepath.Split(filename)
	f, err := os.Create(filepath.Join(dir, "resize_"+file))
	if err != nil {
		return err
	}
	defer f.Close()

	png.Encode(f, img)

	return nil
}
