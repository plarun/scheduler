CMD_CLIENT = ./cmd/client
BIN_CLIENT = ./bin/client
SVC_CLIENT = ./service/client
CMD_EVENTSERVER = ./cmd/eventserver
BIN_EVENTSERVER = ./bin/eventserver
SVC_EVENTSERVER = ./service/eventserver
CMD_VALIDATOR = ./cmd/validator
BIN_VALIDATOR = ./bin/validator
SVC_VALIDATOR = ./service/validator

all: build

# build all services
build: bin_dir build_client build_eventserver build_validator

# create build dir
bin_dir:
	mkdir -p bin

# build client service
build_client:
	go build -o ${BIN_CLIENT} ${CMD_CLIENT}

# build event server service
build_eventserver:
	go build -o ${BIN_EVENTSERVER} ${CMD_EVENTSERVER}

# build validator service
build_validator:
	go build -o ${BIN_VALIDATOR} ${CMD_VALIDATOR}

# start all services
run: build
	./bin/eventserver
	./bin/validator

test: test_client test_eventserver test_validator

# test all tests in client service
test_client:
	go test ${SVC_CLIENT}/...

# test all tests in event server service
test_eventserver:
	go test ${SVC_EVENTSERVER}/...

# test all tests in validator service
test_validator:
	go test ${SVC_VALIDATOR}/...
