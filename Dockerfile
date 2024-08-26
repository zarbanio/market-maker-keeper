FROM --platform=linux/amd64 golang:1.22.1 as builder
WORKDIR /build
COPY . . 
RUN  go env -w GO111MODULE=on && \
     go env -w GOPROXY=https://goproxy.io,direct

RUN go mod download

RUN CGO_ENABLED=0
RUN GOOS=linux
RUN go build main.go

FROM --platform=linux/amd64 gcr.io/distroless/base-debian10
COPY --from=builder /build/main .

ADD store/migrations /migrations

CMD ["./main", "run", "--config", "/etc/app/config.yaml"]