#!/bin/bash
APP_NAME="public-image-upload-server"
docker rm -f $APP_NAME || echo ""
docker run -it \
--name $APP_NAME \
-v $(pwd)/IMAGES:/home/morphs/IMAGES:rw \
-v $(pwd)/IMAGES-ONE-HOT:/home/morphs/IMAGES-ONE-HOT:rw \
--mount type=bind,source="$(pwd)"/config.json,target=/home/morphs/ImageUploadServerGo/config.json \
-p 7391:7391 \
$APP_NAME config.json
