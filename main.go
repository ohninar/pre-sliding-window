package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/nfnt/resize"
)

type pixel struct {
	r, g, b, a uint8
	normal     float32
}

func main() {
	/* O que preciso fazer:
		   1) (OK) Percorrer o diretório por imagens;
		   2) (OK) Abrir imagens;
		   3) (OK) Transformar em 30x30;
		   4) (OK) Transformar em preto e branco;
		   5) (OK) Obter valor decimal de cada pixel;
		   6) (OK) Organizar a saída como uma linha de matriz;
		   7) Subistituir os valores de pixel entre 0 e 1;
		   8) Incluir na última coluna com o label que identifique a imagem entre alpha-numeric

	     Matriz 0 255 cinza
	     Matriz normalizada cinza
	     Matriz preto e branco
	*/

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("    -path=\"./in_imagens\" -resX=30 -resY=30 -normal=false -label=0\n")
		flag.PrintDefaults()
	}

	pathImages := flag.String("path", "/tmp", "a string")
	resX := flag.Int("resX", 30, "an int")
	resY := flag.Int("resY", 30, "an int")
	normal := flag.Bool("normal", false, "a bool")
	label := flag.Int("label", 0, "an int")

	flag.Parse()

	if len(os.Args) != 6 {
		flag.Usage()
		fmt.Println(flag.Args())
		fmt.Println(flag.NArg())
		fmt.Println(os.Args)
		fmt.Println(len(os.Args))
		os.Exit(1)
	}

	images := getImages(*pathImages, uint(*resX), uint(*resY), *normal)

	for _, img := range images {
		matrix := ""
		for _, pixel := range img {
			if *normal {
				matrix += strconv.FormatFloat((normalization(pixel.r)), 'f', 2, 64) + ", "
			} else {
				matrix += strconv.Itoa(int(pixel.r)) + ", "
			}
		}
		matrix += strconv.Itoa(*label) + ""
		fmt.Println(matrix)
	}

}

func getImages(dir string, resX uint, resY uint, normal bool) [][]pixel {

	var images [][]pixel

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Println("erro", dir)
		os.Exit(1)
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		img := loadImage(path)
		img = resize.Resize(resX, resY, img, resize.Lanczos3)
		pixels := getPixels(img, normal)
		images = append(images, pixels)
		return nil

	})

	return images
}

func loadImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func getPixels(img image.Image, normal bool) []pixel {

	bounds := img.Bounds()
	pixels := make([]pixel, bounds.Dx()*bounds.Dy())

	for i := 0; i < bounds.Dx()*bounds.Dy(); i++ {
		x := i % bounds.Dx()
		y := i / bounds.Dx()

		r, g, b, a := color.GrayModel.Convert(img.At(x, y)).RGBA()

		pixels[i].r = uint8(r)
		pixels[i].g = uint8(g)
		pixels[i].b = uint8(b)
		pixels[i].a = uint8(a)
	}
	return pixels
}

func normalization(value uint8) float64 {
	return float64(value) / 255.0

}
