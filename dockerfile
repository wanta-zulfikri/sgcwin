FROM golang:alpine AS builder

WORKDIR /app/

COPY . .

RUN echo 

RUN go mod tidy

RUN go build -o /app/bin /app/main.go

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /app/ .

CMD [ "/app/bin" ]