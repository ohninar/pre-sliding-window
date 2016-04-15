run:
	go run main.go -path="./in_imagens" -resX=30 -resY=30 -normal=false -label=0

gerar-gray:
	sh gray.sh

gerar-gray-normalizado:
	sh gray-normalizado.sh
