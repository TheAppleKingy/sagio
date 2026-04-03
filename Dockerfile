FROM golang:1.26.1-alpine AS builder

WORKDIR /sagio
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
RUN chmod +x server

FROM scratch
WORKDIR /etc/sagio

COPY --from=builder /sagio/server .

ENTRYPOINT [ "/etc/sagio/server" ]