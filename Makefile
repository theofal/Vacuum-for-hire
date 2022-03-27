ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run:
	go run  $(PROJECT_PATH)
update-dependencies:
	go get -u $(PROJECT_PATH)/...
verify-dependencies:
	go test all