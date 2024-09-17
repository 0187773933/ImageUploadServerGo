#!/bin/bash

# Dynamically find the version of ImageMagick installed via Homebrew
IMAGEMAGICK_VERSION=$(brew info imagemagick --json | jq -r '.[0].installed[0].version')
IMAGEMAGICK_PATH="/usr/local/Cellar/imagemagick/$IMAGEMAGICK_VERSION"

# Use pkg-config to set the environment variables automatically
PKG_CONFIG_PATH="/usr/local/Cellar/imagemagick/7.1.1-36/lib/pkgconfig" \
CGO_CFLAGS_ALLOW='-Xpreprocessor' \
CGO_CFLAGS=$(pkg-config --cflags MagickWand) \
CGO_LDFLAGS=$(pkg-config --libs MagickWand) \
go run main.go config.json