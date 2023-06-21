ROOT_DIR			:= $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
TARGET_FOLDER		:= $(ROOT_DIR)/build
CONSUMER_MAIN_FILE	:= $(ROOT_DIR)/src/consumer-srv/main.go
CONSUMER_TARGET		:= $(TARGET_FOLDER)/consumer
SSO_MAIN_FILE		:= $(ROOT_DIR)/src/sso-srv/main.go
SSO_TARGET			:= $(TARGET_FOLDER)/sso


.PHONY: build consumer sso startconsumer startsso test

build:
	@rm -f $(target_name) && go build -o $(target_name) $(target_file)

consumer:
	@make build target_name=$(CONSUMER_TARGET) target_file=$(CONSUMER_MAIN_FILE)

sso:
	@make build target_name=$(SSO_TARGET) target_file=$(SSO_MAIN_FILE)

startconsumer:
	@$(CONSUMER_TARGET)

startsso:
	@$(SSO_TARGET)

test:
	@go vet ./... && go test ./...
