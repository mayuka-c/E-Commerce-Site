FROM golang:1.20.2-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./

RUN go build -o e-commerce .

FROM alpine:3.17.2

COPY --from=builder /app/e-commerce .
EXPOSE 8181
CMD [ "./e-commerce" ]