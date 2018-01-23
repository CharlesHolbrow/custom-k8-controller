#!/bin/bash

if [ "$1" == "" ]; then
    echo "No tag provided"
    exit 1
fi

set -x #echo on

env GOOS=linux GOARCH=amd64 go build -o app && \
docker build --tag "gcr.io/udacity-charles-initial/podcount:$1" . && \
gcloud docker -- push "gcr.io/udacity-charles-initial/podcount:$1" && \
kubectl delete deployment demo

# note that if the build fails, we will not "kubectl delete" so run (below) will likely fail

kubectl run --rm -i demo "--image=gcr.io/udacity-charles-initial/podcount:$1"
