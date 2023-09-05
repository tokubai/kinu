FROM debian:buster-slim
MAINTAINER Takatoshi Maeda <me@tmd.tw>

ENV PATH $PATH:/usr/local/go/bin:/usr/local/go/vendor/bin

WORKDIR /tmp
RUN env DEBIAN_FRONTEND=noninteractive apt update && \
    apt install -y libwebp-dev libpng-dev ghostscript pkg-config \
                   git wget build-essential && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

ENV LIBJPEG_VERSION=2.1.1
ENV LIBJPEG_DPKG_URL=https://sourceforge.net/projects/libjpeg-turbo/files/${LIBJPEG_VERSION}/libjpeg-turbo-official_${LIBJPEG_VERSION}_amd64.deb/download
RUN wget $LIBJPEG_DPKG_URL -O libjpeg-turbo-official_${LIBJPEG_VERSION}_amd64.deb && \
    dpkg -i libjpeg-turbo-official_${LIBJPEG_VERSION}_amd64.deb && \
    ln -fs /opt/libjpeg-turbo/include/*.h /usr/include/ && \
    ln -fs /opt/libjpeg-turbo/lib64/lib* /usr/lib/x86_64-linux-gnu/ && \
    ldconfig && \
    rm -rf /tmp/*

ENV IMAGE_MAGICK_VERSION=6.9.12-30
RUN wget https://download.imagemagick.org/ImageMagick/download/releases/ImageMagick-${IMAGE_MAGICK_VERSION}.tar.xz && \
    tar xvf ImageMagick-${IMAGE_MAGICK_VERSION}.tar.xz && \
    cd ImageMagick-${IMAGE_MAGICK_VERSION} && ./configure &&  make  && make install && ldconfig && \
    rm -rf /tmp/*

ENV GOLANG_VERSION 1.16.5
ENV GOROOT /usr/local/go
ENV GOPATH /usr/local/go/vendor

ENV KINU_VERSION 1.0.0
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
    cd /kinu-build && \
    go build -o /usr/local/bin/kinu . && \
    mkdir -p /var/local/kinu && \
    rm -rf /usr/local/go /root/.cache /tmp/*

CMD ["kinu"]
