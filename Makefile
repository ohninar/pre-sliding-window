run:
	go run main.go -path="./in_imagens" -resX=30 -resY=30 -normal=true -label=0

gerar-gray:
	sh gray.sh

gerar-gray-normalizado:
	sh gray-normalizado.sh

gerar-black-and-white:
	sh black-and-white.sh

gerar-black-and-white-normalizado:
	sh black-and-white-normalizado.sh
