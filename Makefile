ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run:
	go run  $(PROJECT_PATH)
