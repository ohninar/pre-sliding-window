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
		fmt.Println("erro:", dir)
		os.Exit(1)
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		img := loadImage(path)
		saveFile("1-"+info.Name(), img)
		img = resize.Resize(resX, resY, img, resize.Lanczos3)
		saveFile("2-"+info.Name(), img)
		img = escalaCinza(img)
		saveFile("3-"+info.Name(), img)
		img = escalaPretoBranco(img)
		saveFile("4-"+info.Name(), img)
		pixels := getPixels(img)
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

func getPixels(img image.Image) []pixel {

	bounds := img.Bounds()
	pixels := make([]pixel, bounds.Dx()*bounds.Dy())

	i := 0
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[i].r = uint8(r)
			pixels[i].g = uint8(g)
			pixels[i].b = uint8(b)
			pixels[i].a = uint8(a)
			i++
		}
	}

	return pixels
}

func escalaCinza(img image.Image) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	imgRect := image.Rect(0, 0, w, h)
	gray := image.NewGray(imgRect)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(oldColor)
			gray.Set(x, y, grayColor)
		}
	}
	return gray
}

func escalaPretoBranco(img image.Image) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	imgRect := image.Rect(0, 0, w, h)
	gray := image.NewGray(imgRect)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, _, _, _ := img.At(x, y).RGBA()
			if r > 32767 {
				r = 255
			} else {
				r = 0
			}
			//gray.Set(x, y, color.Gray{uint8((r*299 + g*587 + b*114) / 1000)})
			gray.Set(x, y, color.Gray{uint8(r)})

		}
	}
	return gray
}

func normalization(value uint8) float64 {
	return float64(value) / 255.0

}

func saveFile(path string, file image.Image) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("erro:", err)
	}
	defer f.Close()

	err = png.Encode(f, file)
	if err != nil {
		fmt.Println("erro:", err)
	}

}
