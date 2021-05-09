FROM golang:1.16-alpine AS builder

RUN mkdir /out
WORKDIR /out

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o app ./cmd/main.go

FROM gcr.io/distroless/base

COPY --from=builder ./out/app .

CMD ["./app"]