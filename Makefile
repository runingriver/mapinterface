.PHONY: lint


lint:
	find . -name "*.go"  | grep -v mocks | xargs goimports -w
	find . -name "*.go"  | grep -v mocks | xargs gofmt -w
	go vet . ./api/... ./itferr/... ./mapitf/... ./pkg/...

test:
	go test -gcflags=-l -v -count=1 .