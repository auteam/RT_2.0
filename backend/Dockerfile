FROM golang AS builder
COPY . .
ENV GOPATH=/go
ENV GOBIN=/bin
RUN go get && go build -o backend
CMD ./backend
