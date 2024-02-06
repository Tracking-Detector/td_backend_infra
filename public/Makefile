.PHONY: init assets generate run test

init:
	npm install
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/vektra/mockery/v2@latest
	go mod tidy

assets:
	tailwindcss -i ./assets/tailwind.css -o ./static/styles.css

generate:
	@templ generate

run: assets generate
	@go run cmd/main.go

test:
	go test ./...
