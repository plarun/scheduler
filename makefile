CMD_CLIENT = ./cmd/client
BIN_CLIENT = ./bin/client
SVC_CLIENT = ./service/client
CMD_EVENTSERVER = ./cmd/eventserver
BIN_EVENTSERVER = ./bin/eventserver
SVC_EVENTSERVER = ./service/eventserver
CMD_VALIDATOR = ./cmd/validator
BIN_VALIDATOR = ./bin/validator
SVC_VALIDATOR = ./service/validator
CMD_ALLOCATOR = ./cmd/allocator
BIN_ALLOCATOR = ./bin/allocator
SVC_ALLOCATOR = ./service/allocator
CMD_WORKER = ./cmd/worker
BIN_WORKER = ./bin/worker
SVC_WORKER = ./service/worker

all: build

# build all services
build: bin_dir build_client build_eventserver build_validator build_allocator build_worker

# create build dir
bin_dir:
	mkdir -p bin
	mkdir -p cmd ${CMD_CLIENT} ${CMD_EVENTSERVER} ${CMD_VALIDATOR} ${CMD_ALLOCATOR} ${CMD_WORKER}

# build client service
build_client:
	go build -o ${BIN_CLIENT} ${CMD_CLIENT}

# build event server service
build_eventserver:
	go build -o ${BIN_EVENTSERVER} ${CMD_EVENTSERVER}

# build validator service
build_validator:
	go build -o ${BIN_VALIDATOR} ${CMD_VALIDATOR}

# build allocator service
build_allocator:
	go build -o ${BIN_ALLOCATOR} ${CMD_ALLOCATOR}

# build worker service
build_worker:
	go build -o ${BIN_WORKER} ${CMD_WORKER}

# start all services
run: build start_eventserver start_validator start_allocator start_worker

# start eventserver service in background
eventserver:
	./bin/eventserver &

# start validator service in background
validator:
	./bin/validator &

# start allocator service in background
allocator:
	./bin/allocator &

# start worker service in background
worker:
	./bin/worker &

test: test_client test_eventserver test_validator test_allocator test_worker

# test all tests in client service
test_client:
	go test ${SVC_CLIENT}/...

# test all tests in event server service
test_eventserver:
	go test ${SVC_EVENTSERVER}/...

# test all tests in validator service
test_validator:
	go test ${SVC_VALIDATOR}/...

# test all tests in allocator service
test_allocator:
	go test ${SVC_ALLOCATOR}/...

# test all tests in worker service
test_worker:
	go test ${SVC_WORKER}/...
