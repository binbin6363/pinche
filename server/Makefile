.PHONY: build build-all clean vet fmt run

GOCMD=GO111MODULE=on go
BINARY=server
GOARGS=-ldflags="-s -w"

build:
	@echo "build go binary..."
	${GOCMD} build ${GOARGS} -o ${BINARY} cmd/main.go
	@echo "build go binary done"

build-all:
	@echo "build all go files..."
	${GOCMD} build ${GOARGS} ./...
	@echo "build all go files done"

clean:
	@echo "clean go binary..."
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	@echo "clean go binary done"

vet:
	@echo "vet go files..."
	${GOCMD} vet ${GOARGS} ./...
	@echo "vet go files done"

fmt:
	@echo "fmt go files..."
	${GOCMD} fmt ${GOARGS} ./...
	@echo "fmt go files done"

run:
	@echo "run server..."
	${GOCMD} run ${GOARGS} cmd/main.go 
