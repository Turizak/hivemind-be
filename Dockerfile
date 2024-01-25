FROM golang:1.21-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /app/main

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/main .

CMD ["./main"]