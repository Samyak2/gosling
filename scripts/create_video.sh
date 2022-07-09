#!/bin/bash

set -euxo pipefail

imagepath="$1"
inputaudiopath="$2"
outputvideopath="$3"

ffmpeg -loop 1 -i "$imagepath" -i "$inputaudiopath" -c:a copy -c:v libx264 -preset ultrafast -tune stillimage -pix_fmt yuv420p -shortest -vf "pad=ceil(iw/2)*2:ceil(ih/2)*2" "$outputvideopath"

