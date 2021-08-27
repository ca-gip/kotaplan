#!/bin/bash

checksum () {
   curl -s https://api.github.com/repos/ca-gip/kotaplan/releases/latest \
  | grep browser_download_url \
  | grep checksums \
  | cut -d '"' -f 4 \
  | xargs curl -sL
}

if [ "$(uname)" == "Darwin" ]; then
  echo "Downloading Darwin Release"
  mkdir -p /var/tmp/kotaplan
  curl -s https://api.github.com/repos/ca-gip/kotaplan/releases/latest \
    | grep browser_download_url \
    | grep darwin_amd64 \
    | cut -d '"' -f 4 \
    | xargs curl -sL \
    | tar xf - -C /var/tmp/kotaplan/
    sudo sh -c 'mv /var/tmp/kotaplan/kotaplan /usr/local/bin/ && chmod +x /usr/local/bin/kotaplan'
    rm -rf /var/tmp/kotaplan
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
  echo "Downloading Linux Release"
  mkdir -p /tmp/kotaplan
  curl -s  https://api.github.com/repos/ca-gip/kotaplan/releases/latest \
    | grep browser_download_url \
    | grep linux_amd64 \
    | cut -d '"' -f 4 \
    | xargs curl -sL \
    | tar xzf - -C /tmp/kotaplan
    sudo sh -c 'mv /tmp/kotaplan/kotaplan /usr/local/bin/ && chmod +x /usr/local/bin/kotaplan'
    rm -rf /tmp/kotaplan
else echo "Unsupported OS" && exit 1
fi

echo "Install done !"