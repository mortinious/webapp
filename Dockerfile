# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN mkdir -p $GOPATH/src/webapp
ADD . /go/src/webapp


# Build the app command inside the container.
#RUN godep save
#RUN go install webapp

# Run the app command by default when the container starts.
WORKDIR $GOPATH/src/webapp
#ENTRYPOINT /go/bin/webapp

RUN go build -o webapp . 
CMD ["/go/src/app/webapp"]

# Document that the service listens on port 8080.
EXPOSE 8080