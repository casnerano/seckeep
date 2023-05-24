project-build:
	cp .env.dist .env
	go mod download && go mod verify
	docker compose build
project-init: project-build
project-run:
	docker compose up -d
project-test:
	go test -coverprofile=coverage.out ./... > /dev/null && go tool cover -func=coverage.out
server-build:
	CGO_ENABLED=0 go build -o ./seckeep-server ./cmd/server/main.go
client-build:
	CGO_ENABLED=0 go build -o ./seckeep ./cmd/client/main.go
load-example:
	./seckeep data create credential --login="javascript" --password="null-is-object???" --meta="For e-mail account" --meta="Work account" > /dev/null
	./seckeep data create card --number="4012888888881881" --month-year="06.28" --owner="Ivan Ivanov" --cvv="732" --meta="My debit visa card" > /dev/null
	./seckeep data create text --value="My secret plan to develop a new JS-library." --meta="Secret plan" > /dev/null
	./seckeep data create text --value="Not a bug, but a feature." --meta="My list of aphorisms" > /dev/null
	./seckeep data create document --file="./Makefile" --meta="Example doc" > /dev/null
	./seckeep data list
