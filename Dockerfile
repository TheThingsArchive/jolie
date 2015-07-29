# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:latest

#Cache dependencies
RUN go get	"github.com/streadway/amqp"

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/thethingsnetwork/jolie
WORKDIR /go/src/github.com/thethingsnetwork/jolie

RUN go build .

CMD ["./jolie"]
