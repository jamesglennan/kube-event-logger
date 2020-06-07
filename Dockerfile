FROM golang:1.13.12-alpine3.12
 
WORKDIR /
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kube-event-logger cmd/kube-event-logger/main.go

FROM alpine:latest  
WORKDIR /
COPY --from=0 /kube-event-logger .
CMD ["./kube-event-logger"]