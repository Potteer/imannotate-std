#!/bin/bash
mkdir -m 700 -p /tmp/s3data /tmp/s3metadata
docker stop -t 0 s3server ; \
	docker run -d --rm --name s3server -p 9000:8000 \
	-v /tmp/s3data/:/s3data \
       	-v /tmp/s3metadata/:/s3metadata \
	--env S3DATAPATH=/s3data \
	--env S3METADATAPATH=/s3metadata \
	scality/s3server

#docker stop -t 0 s3server ; docker run -d --rm --name s3server -p 9000:8000 -v /tmp/s3data/:/s3data -v /tmp/s3metadata/:/s3metadata --env S3DATAPATH=/s3data --env S3METADATAPATH=/s3metadata -e SCALITY_ACCESS_KEY_ID=accessKey1 -e SCALITY_SECRET_ACCESS_KEY=verySecretKey1 zenko/cloudserver

# zenko/cloudserver
echo """
   exec:
   mkdir -p /tmp/scality
   s3fs scality /tmp/scality -o allow_other,use_path_request_style,nonempty,url=http://localhost:8000,passwd_file=/home/fox/passwd-scality
"""

cat /home/fox/passwd-scality
