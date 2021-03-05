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
ENV ASSETS=./assets
ENV ROLE_COLOR="16280460"

ENV FEEDBACK_WEBHOOK=
ENV DEBUG_WEBHOOK=
ENV DEBUG=false

CMD /app/sweetheart
