package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

func showUsage() {
	println("[USAGE] tocsv (input dir) (out.csv)")
}

func main() {
	//
	if len(os.Args) < 2 {
		showUsage()
		return
	}
	indir := os.Args[1]
	outfile := "out.csv"
	if len(os.Args) >= 3 {
		outfile = os.Args[2]
	}
	//
	files, err := filepath.Glob(indir + "/*.*")

	fp, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fp.Close()

	for _, file := range files {
		ext := filepath.Ext(file)
		ext = strings.ToLower(ext)
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
			// ok
		} else {
			continue
		}
		s, err := toCSV(file)
		if err != nil {
			fmt.Println("[error]", file, ":", err)
			continue
		}
		fp.WriteString(s)
	}
}

func toCSV(path string) (string, error) {
	reader, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	img, _, err := image.Decode(reader)
	if err != nil {
		return "", err
	}
	bounds := img.Bounds()
	csv := ""
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r = r & 0xFF
			g = g & 0xFF
			b = b & 0xFF
			csv += fmt.Sprintf("%d,%d,%d,", r, g, b)
			// println(x, y, "=", r, g, b)
		}
	}
	if len(csv) > 0 {
		csv = csv[0:len(csv)-1] + "\r\n"
	}
	return csv, nil
}
