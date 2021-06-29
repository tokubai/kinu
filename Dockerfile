FROM debian:buster-slim
MAINTAINER Takatoshi Maeda <me@tmd.tw>

ENV PATH $PATH:/usr/local/go/bin:/usr/local/go/vendor/bin

WORKDIR /tmp
RUN env DEBIAN_FRONTEND=noninteractive apt update && \
    apt install -y libwebp-dev libjpeg-dev libpng-dev pkg-config \
                   git wget build-essential && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

ENV IMAGE_MAGICK_VERSION=6.9.12-17
RUN wget https://download.imagemagick.org/ImageMagick/download/releases/ImageMagick-${IMAGE_MAGICK_VERSION}.tar.gz && \
    tar xvzf ImageMagick-${IMAGE_MAGICK_VERSION}.tar.gz && \
    cd ImageMagick-${IMAGE_MAGICK_VERSION} && ./configure &&  make  && make install && ldconfig && \
    rm -rf /tmp/*

ENV GOLANG_VERSION 1.16.5
ENV GOROOT /usr/local/go
ENV GOPATH /usr/local/go/vendor

ENV KINU_VERSION 1.0.0.alpha13
ENV KINU_BIND 0.0.0.0:80
ENV KINU_LOG_LEVEL info
ENV KINU_LOG_FORMAT ltsv
ENV KINU_RESIZE_ENGINE ImageMagick
ENV KINU_STORAGE_TYPE File
ENV KINU_FILE_DIRECTORY /var/local/kinu

WORKDIR /kinu-build
COPY . .
RUN mkdir -p /tmp/golang && \
    wget https://storage.googleapis.com/golang/go${GOLANG_VERSION}.linux-amd64.tar.gz -q -P /tmp/golang && \
    cd /tmp/golang && \
    tar zxf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    mv ./go /usr/local/go && \
    rm -rf /tmp/* && \
    cd /kinu-build && \
    go build -o /usr/local/bin/kinu . && \
    mkdir -p /var/local/kinu && \
    rm -rf /usr/local/go /root/.cache

CMD ["kinu"]
