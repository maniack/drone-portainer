FROM golang:alpine AS builder
ENV CGO_ENABLED=0
RUN apk --update --no-cache add git build-base upx
COPY . $GOPATH/src/github.com/maniack/drone-portainer
WORKDIR $GOPATH/src/github.com/maniack/drone-portainer
RUN go get && go build -ldflags "-s -w -X main.version=$(git describe --tags --long --always)" -o drone-portainer
RUN upx -q --best --ultra-brute -o /bin/drone-portainer drone-portainer

FROM alpine:latest
RUN apk --update --no-cache add ca-certificates
COPY --from=builder /bin/drone-portainer /bin/drone-portainer
ENTRYPOINT /bin/drone-portainer
