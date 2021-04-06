FROM golang:1.16.3-alpine3.12
RUN mkdir /app
ADD . /app
WORKDIR /app
## pull in any dependencies
RUN go mod download
## build with the necessary go libraries included.
RUN go build -o main .

EXPOSE 8000

CMD ["/app/main"]