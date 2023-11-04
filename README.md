# Image Upload Server

## ImageMagick Linux Install

1. `sudo apt-get install imagemagick libmagickwand-dev`
2. `go get gopkg.in/gographics/imagick.v2/imagick`
3. `sudo nano /etc/environment`
4. `PKG_CONFIG_PATH="/usr/local/lib/pkgconfig"`
5. `CGO_CFLAGS_ALLOW='-Xpreprocessor'`
5. `source /etc/environment`
6. `pkg-config --cflags --libs MagickWand`

## ImageMagick OSX Install
- https://github.com/gographics/imagick/issues/286

1. `brew install pkg-config imagemagick@6`
2. `go get gopkg.in/gographics/imagick.2/imagick`
3. go run main.go config.json

## Old - ImageMagick OSX Install
3. `echo $PKG_CONFIG_PATH`
3. `find /usr/local/Cellar -name "ImageMagick*.pc"
4. https://trac.macports.org/wiki/Migration
4. `export PKG_CONFIG_PATH="/path/to/imagemagick/pkgconfig:$PKG_CONFIG_PATH"`
3. `export PKG_CONFIG_PATH=/usr/local/lib/pkgconfig`
4. `export CGO_CFLAGS_ALLOW='-Xpreprocessor'`
5. `pkg-config --cflags --libs MagickWand`

## TODO

1. Fix Import of config.json file ??
