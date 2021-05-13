#!/bin/bash
IMAGE_NAME="ifps-andon-datasource"
DOCKER_REPO="iiicondor/$IMAGE_NAME"
VERSION="1.0.8"

# docker build -t $DOCKER_REPO:$VERSION .
# docker push $DOCKER_REPO:$VERSION
# docker tag $DOCKER_REPO:$VERSION $DOCKER_REPO:dev
# docker push $DOCKER_REPO:dev
# MESSAGE="change Panel8,Panel9"
# echo "[`date "+%Y-%m-%d %H:%M:%S"`] $VERSION => dev {$MESSAGE}" >> ImageInfo.txt

docker pull $DOCKER_REPO:$VERSION
docker tag $DOCKER_REPO:$VERSION $DOCKER_REPO:demo
docker push $DOCKER_REPO:demo
echo "[`date "+%Y-%m-%d %H:%M:%S"`] $VERSION => demo" >> ImageInfo.txt

docker rmi -f $(docker images | grep $DOCKER_REPO | awk '{print $3}')
docker image prune -f
