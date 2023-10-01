run:
	cd website && pnpm i && pnpm build
	go run .

build:
	cd website && pnpm i && pnpm build
	go build -tags "production" -ldflags "-s -w" -o bin/main .
