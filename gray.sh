#!/bin/bash

diretorio1=$(pwd)/English/Img/GoodImg/Bmp/Sample00
diretorio2=$(pwd)/English/Img/GoodImg/Bmp/Sample0

for i in $(seq 62);
do
  if [ $i -le 9 ]; then
    go run main.go -path=$diretorio1$i -resX=30 -resY=30 -normal=false -label=$i > out_imagens/gray/$i.txt
  else
    go run main.go -path=$diretorio2$i -resX=30 -resY=30 -normal=false -label=$i > out_imagens/gray/$i.txt
  fi
done
