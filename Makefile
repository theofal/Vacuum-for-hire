ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# LOCAL
run:
	go run  $(PROJECT_PATH) $(WEBDRIVER_PATH)
update-dependencies:
	go get -u $(PROJECT_PATH)/...
verify-dependencies:
	go test all

# DOCKER
build-docker:
	docker build -t vacuum-for-hire .
run-docker:
	docker run -it --rm --name vacuum-for-hire vacuum-for-hire
shell-docker:
	docker container run -it vacuum-for-hire /bin/bash
docker-clean-images:
	pushd utils/ && \
	./clean.sh && \
	sh popd