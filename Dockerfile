FROM golang:1.22-alpine As builder

WORKDIR /phoo/

RUN apk update && apk add git upx

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /usr/bin/phoo ./cmd/

RUN upx -9 /usr/bin/phoo

FROM alpine

WORKDIR /phoo/

COPY --from=builder /usr/bin/phoo /usr/bin/phoo