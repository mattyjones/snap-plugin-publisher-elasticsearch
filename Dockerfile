FROM gliderlabs/alpine:latest

ENV GOBIN=/usr/bin
ENV GOPATH=/go
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/go/bin

RUN apk add --update \
  alpine-sdk \
  bash \
  go \
  && rm -rf /var/cache/apk/*

RUN mkdir /go && mkdir /go/bin && mkdir /go/src
RUN mkdir /go/src/formational\.net && mkdir /go/src/formational\.net/snap-plugin-publisher-elasticsearch
RUN curl https://glide.sh/get | sh
RUN cd /go/bin && go get github.com/Sirupsen/logrus
RUN cd /go/bin && wget http://snap.ci.snap-telemetry.io/snap/latest_build/linux/x86_64/snapteld
RUN cd /go/bin && wget http://snap.ci.snap-telemetry.io/snap/latest_build/linux/x86_64/snaptel

RUN cd /go/bin && wget https://github.com/intelsdi-x/snap-plugin-publisher-file/releases/download/2/snap-plugin-publisher-file_linux_x86_64
RUN mv /go/bin/snap-plugin-publisher-file_linux_x86_64 /go/bin/snap-plugin-publisher-file
RUN cd /go/bin && wget https://github.com/intelsdi-x/snap-plugin-collector-psutil/releases/download/8/snap-plugin-collector-psutil_linux_x86_64
RUN mv /go/bin/snap-plugin-collector-psutil_linux_x86_64 /go/bin/snap-plugin-collector-psutil
RUN chmod a+x /go/bin/snap*

RUN cd /go && go get -d github.com/intelsdi-x/snap-plugin-lib-go/...
RUN sed -i s/IsVersion3/IsVersion4/g /go/src/github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin/rpc/plugin.pb.go
RUN cd /go && go get github.com/intelsdi-x/snap-plugin-lib-go/...
RUN mkdir -p /var/log/snap
ADD . /go/src/formational.net/snap-plugin-publisher-elasticsearch/
RUN mv /go/src/formational.net/snap-plugin-publisher-elasticsearch/task.yml /tmp/
RUN mv /go/src/formational.net/snap-plugin-publisher-elasticsearch/snap.sh /go/bin/
RUN chmod a+x /go/bin/snap.sh
RUN cd /go/src && go build -o ../bin/snap-plugin-publisher-elasticsearch formational.net/snap-plugin-publisher-elasticsearch/main.go
