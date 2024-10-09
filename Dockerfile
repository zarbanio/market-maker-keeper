FROM golang:1.23.0 as builder
WORKDIR /build
COPY . . 
RUN  go env -w GO111MODULE=on && \
     go env -w GOPROXY=https://goproxy.io,direct

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags=-checklinkname=0 -o main main.go

FROM gcr.io/distroless/base-debian12

COPY --from=builder /build/main .
ADD store/migrations /migrations
EXPOSE 8080

ENTRYPOINT [ "./main","run","--config","/etc/app/config.yaml"]