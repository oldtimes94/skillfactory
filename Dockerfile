FROM golang:1.22 AS builder
WORKDIR /go/src/pipeline
COPY go.mod go.sum ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /pipeline

FROM alpine:3.15
LABEL version="1.0.1"
LABEL maintainer="mikhail.bespalko@yandex.ru"
WORKDIR /root/
COPY --from=builder /pipeline .
ENTRYPOINT ["./pipeline"]

#docker build -t pipeline .
#docker run -d pipeline
#docker exec -it <container_id> sh
#./pipeline