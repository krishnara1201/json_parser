FROM golang:1.24-alpine AS build
WORKDIR /app

COPY go.mod ./
COPY ast ./ast
COPY lexer ./lexer
COPY parser ./parser
COPY main.go ./

RUN go build -o /bin/json-parser .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /bin/json-parser /usr/local/bin/json-parser

ENV PORT=8080
EXPOSE 8080

CMD ["json-parser"]
