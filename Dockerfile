FROM golang:bookworm

ENV CGO_ENABLED=1
ENV CGO_CFLAGS_ALLOW=-Xpreprocessor

WORKDIR /usr/local/goferbot/lib
# Install ImageMagick dependencies and build tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    pkg-config \
    libpng-dev \
    libjpeg-dev \
    libwebp-dev \
    libtiff-dev \
    libfreetype6-dev \
    liblcms2-dev \
    libopenjp2-7-dev \
    libxml2-dev \
    zlib1g-dev \
    libbz2-dev \
    autoconf \
    automake \
    libtool \ 
    libltdl-dev

# Download and extract ImageMagick 7 source code
RUN wget https://imagemagick.org/archive/ImageMagick-7.1.1-46.tar.gz && \
    tar -xzf ImageMagick-7.1.1-46.tar.gz && \
    rm ImageMagick-7.1.1-46.tar.gz 

WORKDIR ImageMagick-7.1.1-46
# Configure, build, and install ImageMagick
RUN ./configure --prefix=/usr/local --enable-shared --with-modules --with-webp=yes --with-openjp2=yes && \
    make -j$(nproc) && \
    make install && \
    ldconfig

WORKDIR /usr/local/goferbot
COPY go.mod go.sum ./
RUN go mod download 

COPY . .
RUN touch /usr/local/goferbot/sql/chats.db

RUN go build -o gofer .

CMD ["./gofer"]
