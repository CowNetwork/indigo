FROM golang:1.16-alpine AS builder

WORKDIR /out

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

FROM gcr.io/distroless/base

COPY --from=builder ./out/app .

CMD ["./app"]