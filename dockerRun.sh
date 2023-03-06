#!/bin/bash
APP_NAME="public-image-upload-server"
sudo docker rm $APP_NAME || echo ""
#sudo docker run -it $APP_NAME
id=$(sudo docker run -dit \
--name $APP_NAME \
-v $(pwd)/IMAGES:/home/morphs/IMAGES:rw \
-v $(pwd)/IMAGES-ONE-HOT:/home/morphs/IMAGES-ONE-HOT:rw \
--mount type=bind,source="$(pwd)"/config.json,target=/home/morphs/ImageUploadServerGo/config.json \
-p 7391:7391 \
$APP_NAME config.json)
sudo docker logs -f $id
