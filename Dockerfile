FROM golang

COPY boulderbar-bot /go/bin/boulderbar-bot
ENV TOKEN=

ENTRYPOINT /go/bin/boulderbar-bot
