rm -rf dist
env GOOS=linux GOARCH=arm64 go build -v -o ./dist/hfx-api ./cmd/hfx-api
rsync -avz ./dist/ rpi:hfx-api
