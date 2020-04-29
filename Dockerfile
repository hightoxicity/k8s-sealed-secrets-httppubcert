FROM golang:1.13.8-alpine
RUN mkdir -p /go/src/github.com/hightoxicity/k8s-sealed-secrets-httppubcert
WORKDIR /go/src/github.com/hightoxicity/k8s-sealed-secrets-httppubcert
COPY . ./
RUN ls -al /go/src/github.com/hightoxicity/k8s-sealed-secrets-httppubcert
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s -v -extldflags -static" -a main.go

FROM scratch
COPY --from=0 /go/src/github.com/hightoxicity/k8s-sealed-secrets-httppubcert/main /k8s-sealed-secrets-httppubcert
CMD ["/k8s-sealed-secrets-httppubcert"]
