#!/bin/bash

# PKG_CONFIG_PATH=/usr/local/lib/pkgconfig CGO_CFLAGS_ALLOW='-Xpreprocessor' go run main.go config.json
PKG_CONFIG_PATH="/usr/local/Cellar/imagemagick/7.1.1-8_1/lib/pkgconfig" \
CGO_CFLAGS_ALLOW='-Xpreprocessor' \
CGO_CFLAGS="-I/usr/local/Cellar/imagemagick/7.1.1-8_1/include/ImageMagick-7" \
CGO_LDFLAGS="-L/usr/local/Cellar/imagemagick/7.1.1-8_1/lib" \
go run main.go config.json


# export CGO_CFLAGS="-I/usr/local/Cellar/imagemagick/7.1.1-8_1/include/ImageMagick-7"
# export CGO_LDFLAGS="-L/usr/local/Cellar/imagemagick/7.1.1-8_1/lib"

# pkg-config --cflags --libs MagickWand
# -Xpreprocessor -fopenmp -DMAGICKCORE_HDRI_ENABLE=1 -DMAGICKCORE_QUANTUM_DEPTH=16 -Xpreprocessor -fopenmp -DMAGICKCORE_HDRI_ENABLE=1 -DMAGICKCORE_QUANTUM_DEPTH=16   -lMagickWand-7.Q16HDRI -lMagickCore-7.Q16HDRI
