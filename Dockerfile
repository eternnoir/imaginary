# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM ubuntu:16.04
MAINTAINER tomas@aparicio.me

ENV LIBVIPS_VERSION 8.5.6

# Installs libvips + required libraries
RUN \

  # Install dependencies
  apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y \
  automake build-essential curl \
  gobject-introspection gtk-doc-tools libglib2.0-dev libjpeg-turbo8-dev libpng12-dev \
  libwebp-dev libtiff5-dev libgif-dev libexif-dev libxml2-dev libpoppler-glib-dev \
  swig libmagickwand-dev libpango1.0-dev libmatio-dev libopenslide-dev libcfitsio3-dev \
  libgsf-1-dev fftw3-dev liborc-0.4-dev librsvg2-dev && \

  # Build libvips
  cd /tmp && \
  curl -OL https://github.com/jcupitt/libvips/archive/v${LIBVIPS_VERSION}/vips-${LIBVIPS_VERSION}.tar.gz && \
  tar zvxf vips-${LIBVIPS_VERSION}.tar.gz && \
  cd /tmp/libvips-${LIBVIPS_VERSION} && \
  ./configure --enable-debug=no --without-python $1 && \
  make && \
  make install && \
  ldconfig && \

  # Clean up
  apt-get remove -y curl automake build-essential && \
  apt-get autoremove -y && \
  apt-get autoclean && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Server port to listen
ENV PORT 9000

# Go version to use
ENV GOLANG_VERSION 1.9.2

# gcc for cgo
RUN apt-get update && apt-get install -y \
    gcc curl git libc6-dev make ca-certificates \
    --no-install-recommends \
  && rm -rf /var/lib/apt/lists/*

ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 de874549d9a8d8d8062be05808509c09a88a248e77ec14eb77453530829ac02b

RUN curl -fsSL --insecure "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
  && echo "$GOLANG_DOWNLOAD_SHA256 golang.tar.gz" | sha256sum -c - \
  && tar -C /usr/local -xzf golang.tar.gz \
  && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

# Fetch the latest version of the package
RUN go get -u golang.org/x/net/context
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/aws/aws-sdk-go/aws

# Copy imaginary sources
COPY . $GOPATH/src/github.com/h2non/imaginary

# Compile imaginary
RUN go build -race -o bin/imaginary github.com/h2non/imaginary

# Run the entrypoint command by default when the container starts.
ENTRYPOINT ["bin/imaginary"]

# Expose the server TCP port
EXPOSE 9000
