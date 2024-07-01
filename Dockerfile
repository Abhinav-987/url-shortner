FROM golang:alpine

WORKDIR /project/go-docker/

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o /project/go-docker/build/myapp .

EXPOSE 8000
ENTRYPOINT [ "/project/go-docker/build/myapp" ]
