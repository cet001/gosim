all : install

clean :
	@echo ">>> Cleaning and initializing gosim project <<<"
	@go clean
	@gofmt -w .
	@go get github.com/stretchr/testify

test : clean
	@echo ">>> Running unit tests <<<"
	@go test

test-coverage : clean
	@echo ">>> Running unit tests and calculating code coverage <<<"
	@go test ./ -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out
	@echo ">>> Coverage report saved to ./coverage.html"

install : test
	@echo ">>> Building and installing gosim <<<"
	@go install
	@echo ">>> gosim installed successfully! <<<"
	@echo ""
