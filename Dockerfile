FROM golang:1.20-buster AS build

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -o /bin/app

FROM ubuntu:22.04


COPY --from=build /bin/app /bin

EXPOSE 8000

CMD ["/bin/app"]
