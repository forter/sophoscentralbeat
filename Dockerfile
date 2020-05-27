FROM golang:stretch AS builder
MAINTAINER logrhythm RnD

WORKDIR /go/src/github.com/logrhythm/sophoscentralbeat
RUN mkdir -p /config
RUN apt-get update && \
    apt-get install -y \
    git gcc g++ binutils make
RUN mkdir -p ${GOPATH}/src/github.com/elastic && git clone https://github.com/elastic/beats ${GOPATH}/src/github.com/elastic/beats
COPY . /go/src/github.com/logrhythm/sophoscentralbeat/
RUN make
RUN chmod +x sophoscentralbeat
# ---

FROM ubuntu:latest
COPY --from=builder /go/src/github.com/logrhythm/sophoscentralbeat/sophoscentralbeat /bin/sophoscentralbeat
RUN apt-get -y update \
 && apt-get -y install ca-certificates dumb-init curl \
 && update-ca-certificates
VOLUME  /config/beat.yml
ENTRYPOINT [ "/bin/sophoscentralbeat" ]
CMD [ "-e -c /config/beat.yml" ]
