FROM golang:1.16 AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GOPROXY=direct
env http_proxy=http://child-prc.intel.com:913
env HTTP_PROXY=http://child-prc.intel.com:913
env https_proxy=http://child-prc.intel.com:913
env HTTPS_PROXY=http://child-prc.intel.com:913
env NO_PROXY=intel.com,.intel.com,localhost,127.0.0.1,10.0.0.0/8,192.168.0.0/16,172.18.0.0/12
env no_proxy=intel.com,.intel.com,localhost,127.0.0.1,10.0.0.0/8,192.168.0.0/16,172.18.0.0/12

WORKDIR /work
COPY . /work

# Build admission-webhook
RUN  \
  go build -o bin/admission-webhook .

# ---
FROM ubuntu:22.04 AS run

COPY --from=build /work/bin/admission-webhook /usr/local/bin/
COPY bin/syft /usr/local/bin/
COPY syft.yaml /root/.syft.yaml

CMD ["admission-webhook"]
