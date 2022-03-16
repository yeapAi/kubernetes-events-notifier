FROM golang:1.17-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY . ./
RUN go mod download && go build -v -o main

FROM scratch
COPY --from=build /app/main /
CMD ["/main"]
