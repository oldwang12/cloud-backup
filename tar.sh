#!/bin/bash

# demo: bash tar.sh test /local/data /百度网盘/数据冷备 http://localhost:5244 alist-88bb6ecc-ae41-43b4-b74b-6b104bfa3032I1xheiJw34zUSaI7FQl1JptaI584AO2cGAIIvUwtUGTRWW

# 判断参数个数是否为6个
if [ $# -ne 6 ]; then
  echo "Error: Expected 5 arguments."
  echo "Usage: $0 <filename> <src> <dst> <alist_host> <alist_token>"
  exit 1
fi

FILE_PATH=$1
# 使用awk命令分割字符串并输出最后一个部分
FILE_NAME=$(echo $FILE_PATH | awk -F'/' '{print $NF}')
ALIST_SRC=$2
ALIST_DST=$3
ALIST_HOST=$4
ALIST_DATA_DIR=$5
ALIST_TOKEN=$6

if [ ! -e $FILE_PATH ]; then
    echo "文件不存在"
    exit 1
fi

# 以下是你希望执行的操作，当参数个数为三个时执行
echo "文件名: $FILE_NAME"
echo "Alist src: $ALIST_SRC"
echo "Alist dst: $ALIST_DST"
echo "Alist Host: $ALIST_HOST"
echo "Alist Data dir: $ALIST_DATA_DIR"
echo "Alist Token: $ALIST_TOKEN"

TIME=$(date +"%Y-%m-%d-%H%M")
echo 当前时间: $TIME
FILE_NAME_TAR_GZ="${TIME}-${FILE_NAME}.tar.gz"

echo 备份文件名: $FILE_NAME_TAR_GZ

tar -czf ${ALIST_DATA_DIR}/${FILE_NAME_TAR_GZ} $FILE_PATH

curl "$ALIST_HOST/api/fs/copy" \
  -H 'Accept: application/json, text/plain, */*' \
  -H "Authorization: $ALIST_TOKEN" \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/json;charset=UTF-8' \
  -H "Origin: $ALIST_HOST" \
  -H "Referer: $ALIST_HOST/local/data" \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36' \
  -H 'sec-ch-ua: "Not/A)Brand";v="99", "Google Chrome";v="115", "Chromium";v="115"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "macOS"' \
  --data-raw '{"src_dir":"'"$ALIST_SRC"'","dst_dir":"'"$ALIST_DST"'","names":["'"$FILE_NAME_TAR_GZ"'"]}' \
  --compressed