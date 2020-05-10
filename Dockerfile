FROM shufo/go-graphql-base:1.14 as build

WORKDIR /go/app

COPY . /go/app

RUN go mod download

RUN packr build -o app

FROM alpine

WORKDIR /app

COPY --from=build /go/app/app .

RUN addgroup go \
  && adduser -D -G go go \
  && chown -R go:go /app/app

CMD ["./app"]
