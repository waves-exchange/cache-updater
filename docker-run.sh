# !/bin/bash

docker build -t cache-updater .
docker run -itd --name cache-updater -p 5000:5000 cache-updater