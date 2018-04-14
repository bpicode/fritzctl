FROM ubuntu:xenial

ARG go_version=1.9.2
ARG fritzctl_revision=master
ARG fritzctl_version=unknown

RUN apt-get update \
  && apt-get install -y --no-install-recommends make wget git \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*
RUN wget --quiet https://storage.googleapis.com/golang/go${go_version}.linux-amd64.tar.gz
RUN tar -xf go${go_version}.linux-amd64.tar.gz
RUN mv go /usr/local
RUN mkdir -p /root/go/src/github.com/bpicode

WORKDIR /root/go/src/github.com/bpicode
RUN git clone https://github.com/bpicode/fritzctl.git
WORKDIR /root/go/src/github.com/bpicode/fritzctl
RUN git checkout ${fritzctl_revision}
RUN mkdir build

ENV GOPATH=/root/go
ENV GOROOT=/usr/local/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
ENV FRITZCTL_VERSION=$fritzctl_version

ENTRYPOINT [ "make", "sysinfo", "deps", "dist_all"]