FROM golang:1.19-alpine AS build
RUN apk add build-base

WORKDIR /build
ENV CGO_ENABLED=1
COPY . .
RUN go mod tidy
RUN go test -v -race -cover ./...
RUN GOOS=linux GOARCH=amd64 go build -o api .


FROM alpine:latest
WORKDIR /app
COPY --from=build /build/api /build/resources ./

ENV EMAIL_FROM jobs@example.com
ENV EMAIL_SUBJECT A new job alert has arrived!
ENV EMAIL_TEMPLATE ./resources/body.tmpl
ENV NOTIFICATION_WORKERS 10
ENV HTTP_ADDRESS :8080
ENV DATABASE_NAME production.db
ENV JOBBERWOCKY_URL http://localhost:8090

CMD ["./api"]