FROM golang:1.16-alpine
ENV CGO_ENABLED=0

WORKDIR /usr/src/app
COPY . .
RUN go build -o tfc-badge *.go

FROM alpine
COPY --from=0 /usr/src/app/tfc-badge /

CMD [ "./tfc-badge" ]
