check:
	goimports -w $$(find . -type f -name '*.go' ! -path "./vendor/*" ! -path "./.git/*")
	go vet ./...
	golint $$(find . -type f -name '*.go' ! -path "./vendor/*" ! -path "./.git/*")

test: check
	go test -v ./...
