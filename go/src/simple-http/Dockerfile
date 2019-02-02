FROM golang:1.9 AS builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir /go/src/simple-http
WORKDIR /go/src/simple-http

# Fetch dependencies
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

# Copy the code
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /simple-http ./src/*

FROM scratch
COPY --from=builder /simple-http ./
ENTRYPOINT ["./simple-http"]

