# build stage
FROM golang:alpine AS build-stage
WORKDIR /build
ADD go.mod .
COPY . .
RUN go build -o main cmd/main.go

# production stage
FROM alpine
WORKDIR /build
COPY --from=build-stage /build/main /build/main
CMD ["./main"]