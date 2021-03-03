FROM golang:1.15-alpine3.12

LABEL image="Sweetheart"
LABEL maintainer="github.com/meir"
LABEL madew="love"

WORKDIR /app
COPY . .

RUN go build ./cmd/sweetheart

ARG VERSION=???

ENV VERSION=$VERSION
ENV PREFIX="-"
ENV TOKEN=

CMD /app/sweetheart
