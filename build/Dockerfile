FROM node:18-bookworm AS node_builder

RUN npm i -g npm@latest
RUN npm i -g pnpm@latest

WORKDIR /build

COPY website website

WORKDIR /build/website

RUN pnpm build

FROM golang:1.21-bookworm AS golang_builder

WORKDIR /build

COPY . .
COPY --from=node_builder /build/website/dist /build/website/dist

RUN go vet -v
RUN go test ./... -v --race

RUN CGO_ENABLED=0 go build -tags "production" -ldflags "-s -w" -o bin/main .

FROM gcr.io/distroless/static-debian11

ENV LISTEN_ADDR=":8080"

COPY --from=golang_builder /build/bin/main /bin/app

CMD [ "/bin/app" ]
