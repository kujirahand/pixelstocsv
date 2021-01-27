package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

var resizeX uint = 0
var resizeY uint = 0
var infile string = ""
var outfile string = ""

func showUsage() {
	println(
		"[USAGE]\n" +
			"pixelstocsv (input dir) (out.csv) [opttions]\n" +
			"[options]\n" +
			"--resize=width,height]")
}

func strToUIntDef(s string, def uint) uint {
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	if n < 0 {
		return def
	}
	return uint(n)
}

func main() {
	//
	if len(os.Args) < 2 {
		showUsage()
		return
	}
	// parse Args
	for i, v := range os.Args {
		if i == 0 {
			continue
		}
		if v == "" {
			continue
		}
		ch := v[0]
		if ch == '-' {
			a := strings.SplitN(v, "=", 2)
			key := a[0]
			val := ""
			if len(a) == 2 {
				val = a[1]
			}
			if key == "--resize" {
				axy := strings.Split(val, ",")
				if len(axy) >= 2 {
					resizeX = strToUIntDef(axy[0], 0)
					resizeY = strToUIntDef(axy[1], 0)
				} else {
					resizeX = strToUIntDef(val, 0)
					resizeY = resizeX
				}
				continue
			}
			continue
		}
		if infile == "" {
			infile = v
			continue
		}
		if outfile == "" {
			outfile = v
			continue
		}
	}
	if outfile == "" {
		outfile = "out.csv"
	}
	// open output file
	fp, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fp.Close()

	// file or dir
	fi, err := os.Stat(infile)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !(fi.IsDir()) {
		s, err := toCSV(infile)
		if err != nil {
			fmt.Println(err)
			return
		}
		fp.WriteString(s)
		fmt.Println("ok.")
		return
	}

	files, err := filepath.Glob(infile + "/*.*")
	for _, file := range files {
		s, err := toCSV(file)
		if err != nil {
			fmt.Println("[error]", file, ":", err)
			continue
		}
		fp.WriteString(s)
	}
	fmt.Println("ok.")
}

func toCSV(path string) (string, error) {
	ext := filepath.Ext(path)
	ext = strings.ToLower(ext)
	println("- ", path)
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
		// ok
	} else {
		return "", fmt.Errorf("Invalid File Format")
	}

	// read image
	reader, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	img, _, err := image.Decode(reader)
	if err != nil {
		return "", err
	}
	if resizeX > 0 {
		img = resize.Resize(resizeX, resizeY, img, resize.Lanczos3)
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
