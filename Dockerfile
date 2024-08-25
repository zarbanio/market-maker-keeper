FROM --platform=linux/amd64 registry.hamdocker.ir/bitex/go-builder:v0.4.4 as builder
WORKDIR /build
COPY . . 
RUN  go env -w GO111MODULE=on && \
     go env -w GOPROXY=https://goproxy.io,direct

RUN go mod download &&  make code-gen

RUN CGO_ENABLED=0
RUN GOOS=linux
RUN go build main.go

FROM --platform=linux/amd64 gcr.io/distroless/base-debian10
COPY --from=builder /build/main .

ADD store/migrations /migrations

CMD ["./main", "run", "--config", "/etc/app/config.yaml"]