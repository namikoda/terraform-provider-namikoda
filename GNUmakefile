TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
COVER_TEST?=$$(go list ./... |grep -v 'vendor')

default: build

build: fmtcheck
	go install

clean:
	rm -Rf build/

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

testrace: fmtcheck
	TF_ACC= go test -race $(TEST) $(TESTARGS)

cover:
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test $(COVER_TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile: fmtcheck
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./builtin/providers/aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build crossbuild test testacc testrace cover vet fmt fmtcheck errcheck vendor-status test-compile clean cross-compile dist

dist: build/linux_amd64/terraform-provider-namikoda build/linux_386/terraform-provider-namikoda build/darwin_amd64/terraform-provider-namikoda build/windows_amd64/terraform-provider-namikoda.exe build/windows_386/terraform-provider-namikoda.exe 

build/linux_amd64/terraform-provider-namikoda:
	GOOS=linux GOARCH=amd64 go build -o build/linux_amd64/terraform-provider-namikoda

build/linux_386/terraform-provider-namikoda:
	GOOS=linux GOARCH=386 go build -o build/linux_386/terraform-provider-namikoda

build/darwin_amd64/terraform-provider-namikoda:
	GOOS=darwin GOARCH=amd64 go build -o build/darwin_amd64/terraform-provider-namikoda

build/windows_amd64/terraform-provider-namikoda.exe:
	GOOS=windows GOARCH=amd64 go build -o build/windows_amd64/terraform-provider-namikoda.exe

build/windows_386/terraform-provider-namikoda.exe:
	GOOS=windows GOARCH=386 go build -o build/windows_386/terraform-provider-namikoda.exe

