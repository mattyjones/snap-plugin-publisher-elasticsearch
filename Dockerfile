FROM golang:latest

#ENV GOBIN=/usr/bin
#ENV GOPATH=/go

RUN mkdir -p /go/src/formational\.net && \
  mkdir /go/src/formational\.net/snap-plugin-publisher-elasticsearch && \
  curl https://glide.sh/get | sh && \
  go get github.com/Sirupsen/logrus && \
  wget http://snap.ci.snap-telemetry.io/snap/latest_build/linux/x86_64/snapteld && \
  wget http://snap.ci.snap-telemetry.io/snap/latest_build/linux/x86_64/snaptel && \
  wget https://github.com/intelsdi-x/snap-plugin-publisher-file/releases/download/2/snap-plugin-publisher-file_linux_x86_64 && \
  mv /go/bin/snap-plugin-publisher-file_linux_x86_64 /go/bin/snap-plugin-publisher-file && \
  wget https://github.com/intelsdi-x/snap-plugin-collector-psutil/releases/download/8/snap-plugin-collector-psutil_linux_x86_64 && \
  mv /go/bin/snap-plugin-collector-psutil_linux_x86_64 /go/bin/snap-plugin-collector-psutil && \
  chmod a+x /go/bin/snap* && cd /go && go get -d github.com/intelsdi-x/snap-plugin-lib-go/... && \
  sed -i s/IsVersion3/IsVersion4/g /go/src/github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin/rpc/plugin.pb.go && \
  cd /go && go get github.com/intelsdi-x/snap-plugin-lib-go/... && mkdir -p /var/log/snap

ADD . /go/src/formational.net/snap-plugin-publisher-elasticsearch/

RUN cp /go/src/formational.net/snap-plugin-publisher-elasticsearch/task.yml /tmp/ && \
  cp /go/src/formational.net/snap-plugin-publisher-elasticsearch/snap.sh /go/bin/ && \
  chmod a+x /go/bin/snap.sh
RUN cd /go/src && go build -o ../bin/snap-plugin-publisher-elasticsearch formational.net/snap-plugin-publisher-elasticsearch/main.go
