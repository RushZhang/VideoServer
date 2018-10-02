#! /bin/bash

# Build web UI
# 可以实现热部署，比如改变了templates里的html文件，不用重启ui服务器，直接build就可以了

cd ~/Go_WS/bin
if [ ! -d video_server_web_ui  ];then
  mkdir testgrid
else
  echo dir video_server_web_ui exist
fi

cd ~/Go_WS/src/video_server/web
echo install..
go install
cp ~/Go_WS/bin/web ~/Go_WS/bin/video_server_web_ui/web
cp -R ~/Go_WS/src/video_server/templates ~/Go_WS/bin/video_server_web_ui
