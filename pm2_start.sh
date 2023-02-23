#/bin/bash
pm2 delete IUS || echo ""
pm2 start ./bin/linux/amd64/ImageUploadServer --name IUS -- config.json
pm2 save
pm2 log IUS