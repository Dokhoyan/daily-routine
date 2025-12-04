FROM golang:1.24.1-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/main ./cmd

FROM alpine:3.13 AS final

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /

RUN mkdir -p /logs && chmod 777 /logs

COPY --from=build /bin/main /main

EXPOSE 8000

ENTRYPOINT ["/main"]


