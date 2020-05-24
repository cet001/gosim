all : install

clean :
	@echo ">>> Cleaning and initializing gosim project <<<"
	@go clean
	@gofmt -w .
	@go get github.com/cet001/gometric
	@go get github.com/cet001/mathext
	@go get github.com/stretchr/testify

test : clean
	@echo ">>> Running unit tests <<<"
	@go test ./ ./models/tfidf

test-coverage : clean
	@echo ">>> Running unit tests and calculating code coverage <<<"
	@go test ./ ./models/tfidf -cover

install : test
	@echo ">>> Building and installing gosim <<<"
	@go install
	@echo ">>> gosim installed successfully! <<<"
	@echo ""
