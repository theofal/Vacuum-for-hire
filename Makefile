ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# LOCAL
run:
	go run  $(PROJECT_PATH) WEBDRIVER_PATH=/Users/theofalgayrettes/GolandProjects/Vacuum-for-hire/Utils/macos/chromedriver
update-dependencies:
	go get -u $(PROJECT_PATH)/...
verify-dependencies:
	go test all

# DOCKER
build-docker:
	docker build -t vacuum-for-hire .
run-docker:
	docker run -it --rm --name my-running-app vacuum-for-hire