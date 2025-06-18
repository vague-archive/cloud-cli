alias r := run
alias b := build
alias t := test
alias c := cover

set dotenv-load := true

coverfile := ".coverage"

# list all tasks
[group('other')]
@list:
  just --list

# run <command>
[group('common')]
@run *ARGS:
  go run cmd/main.go {{ARGS}}

# build the executable
[group('common')]
build:
  go build -o void-cloud cmd/main.go

# run all tests
[group('common')]
test:
  go test -coverprofile {{coverfile}} "./..."

# run linter(s)
[group('common')]
lint:
  go vet ./...
  staticcheck -checks=all,-ST1000 ./...

# run formatter
[group('common')]
format:
  go fmt ./...

# run code coverage (cli)
[group('coverage')]
cover: test
  go tool cover -func={{coverfile}}

# run code coverage (html)
[group('coverage')]
cover-html: test
  go tool cover -html={{coverfile}}
