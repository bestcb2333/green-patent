#!/bin/bash
go build -o app . && \
PORT=8702 \
JWT_KEY=114514 \
DB_USER=patent \
DB_NAME=patent \
DB_PASS=patent \
DB_HOST=127.0.0.1 \
DB_PORT=3306 \
PATH_PATENT=/srv/patent/patent/ \
PATH_REPORT=/srv/patent/report/ \
./app
