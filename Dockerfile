FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY /src .

RUN pwd
RUN tree

RUN go mod init myapp
RUN go mod tidy

RUN go build -o backend .

FROM alpine:latest

COPY --from=builder app/backend .
COPY /data /data

EXPOSE 8000

CMD ["./backend"]