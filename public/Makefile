.PHONY: init assets generate run test dev

init:
	npm install
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/vektra/mockery/v2@latest
	go install github.com/cosmtrek/air@latest
	go mod tidy

assets:
	tailwindcss -i ./assets/tailwind.css -o ./static/styles.css

generate:
	@templ generate


test:
	go test ./...

dev:
	@air

stop:
	@pkill -f "cmd/main.go"
