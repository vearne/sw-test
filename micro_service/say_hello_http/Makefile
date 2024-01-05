COMMAND_NAME = sw-sayhello

.PHONY: clean
clean:
	rm -rf $(COMMAND_NAME)

.PHONY: build
build:
	go build -mod=vendor -toolexec="/var/tmp/skywalking-go-agent" -a -o $(COMMAND_NAME) .

.PHONY: run
run:
	env SW_AGENT_NAME=$(COMMAND_NAME) ./$(COMMAND_NAME)