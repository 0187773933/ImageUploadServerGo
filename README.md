# Image Upload Server

## ImageMagick Linux Install

1. `sudo apt-get install libmagickwand-dev`
2. `go get gopkg.in/gographics/imagick.v3/imagick`
3. `sudo nano /etc/environment`
4. `PKG_CONFIG_PATH="/usr/local/lib/pkgconfig"`
5. `source /etc/environment`

## ImageMagick OSX Install

1. `brew install imagemagick`
2. `go get gopkg.in/gographics/imagick.v3/imagick`
3. `export PKG_CONFIG_PATH=/usr/local/lib/pkgconfig`
4. `export CGO_CFLAGS_ALLOW='-Xpreprocessor'`
5. `pkg-config --cflags --libs MagickWand`

## TODO

1. Fix Docker