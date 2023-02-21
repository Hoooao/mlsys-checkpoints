
FROM golang:1.19-alpine
WORKDIR /app
EXPOSE 9090
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build .
CMD [ "go", "run", "."]