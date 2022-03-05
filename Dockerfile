FROM golang:1.17 as builder
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
ENV CGO_ENABLED=0
COPY . /opt/code
RUN cd /opt/code && go build -o service-unit

FROM alpine:3.15.0
COPY --from=builder /opt/code/service-unit /opt/code/service-unit
WORKDIR /opt/code
CMD [ "./service-unit" ]