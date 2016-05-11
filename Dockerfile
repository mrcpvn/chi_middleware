from golang:alpine

ADD . /go/src/github.com/mrcpvn/chi_middleware

RUN go install github.com/mrcpvn/chi_middleware

EXPOSE 8000

ENTRYPOINT ["chi_middleware"]
CMD ["0.0.0.0:8000"]
