FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o md2vim .



FROM scratch
COPY --from=builder /build/md2vim /

# Command to run
ENTRYPOINT ["/md2vim"]
