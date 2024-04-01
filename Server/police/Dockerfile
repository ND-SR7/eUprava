FROM golang:alpine AS build_container
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o police

FROM alpine:3.19
COPY --from=build_container /app/police /usr/bin
EXPOSE 8082
ENTRYPOINT ["police"]
