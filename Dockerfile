FROM ubuntu:xenial
MAINTAINER Takatoshi Maeda <me@tmd.tw>

ENV PATH $PATH:/usr/local/go/bin:/usr/local/go/vendor/bin

RUN env DEBIAN_FRONTEND=noninteractive apt-get update && \
    apt-get build-dep -y imagemagick && \
    apt-get install -y libwebp-dev git wget && \
    rm -rf /var/lib/apt/lists/*

ENV IMAGEMAGICK_VERSION 6.9.4-4
ENV LD_LIBRARY_PATH /usr/local/lib

RUN mkdir -p /tmp/imagemagick && \
    wget http://www.imagemagick.org/download/ImageMagick-${IMAGEMAGICK_VERSION}.tar.gz -q -P /tmp/imagemagick/ && \
    cd /tmp/imagemagick && \
    tar zxf ImageMagick-${IMAGEMAGICK_VERSION}.tar.gz && \
    cd ImageMagick-${IMAGEMAGICK_VERSION} && \
    ./configure \
      --prefix=/usr \
      --libdir=/usr/lib/x86_64-linux-gnu \
      --with-modules \
      --disable-openmp \
      --with-jemalloc && \
    make && \
    make install && \
    rm -rf /tmp/imagemagick

ENV GOLANG_VERSION 1.6.2
ENV GOROOT /usr/local/go
ENV GOPATH /usr/local/go/vendor

RUN mkdir -p /tmp/golang && \
    wget https://storage.googleapis.com/golang/go${GOLANG_VERSION}.linux-amd64.tar.gz -q -P /tmp/golang && \
    cd /tmp/golang && \
    tar zxf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    mv ./go /usr/local/go && \
    mkdir /usr/local/go/vendor && \
    rm -rf /tmp/golang

ENV KINU_VERSION 1.0.0.alpha1
ENV KINU_BIND 127.0.0.1:80
ENV KINU_LOG_LEVEL info
ENV KINU_LOG_FORMAT ltsv
ENV KINU_RESIZE_ENGINE ImageMagick
ENV KINU_STORAGE_TYPE File
ENV KINU_FILE_DIRECTORY /var/local/kinu

RUN go get -d github.com/TakatoshiMaeda/kinu && \
    cd /usr/local/go/vendor/src/github.com/TakatoshiMaeda/kinu/ && \
    git fetch && git checkout refs/tags/${KINU_VERSION} && \
    go build -o /usr/local/go/vendor/bin/kinu . && \
    mkdir -p /var/local/kinu

CMD ["kinu"]
