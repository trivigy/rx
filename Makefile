all: clean docs

list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs | tr -s ' '  '\n'
.PHONY: list

clean:
	@ rm -rf ./bin
.PHONY: clean

docs:
	@ mkdir -p /tmp/go/doc
	@ rm -rf /tmp/go/src/github.com/syncaide/rx
	@ mkdir -p /tmp/go/src/github.com/syncaide/rx
	@ tar -c --exclude='.git' . | tar -x -C /tmp/go/src/github.com/syncaide/rx
	@ echo "http://localhost:6060/pkg/github.com/syncaide/rx"
	@ GOROOT=/tmp/go/ GOPATH=/tmp/go/ godoc -http=localhost:6060
.PHONY: docs
