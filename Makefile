BINARY=drello_binary
API_DOC_TARGET=api.yaml
PROJECT_DIR := ${CURDIR}
DOCUMENTATION_CONTAINER_NAME=documentation
TEST_FLAGS = \
	-covermode=atomic ./... \
	--test-config="$(PROJECT_DIR)/config.yaml"

# тесты
generate-mocks:
	go generate ./...
test-cover:
	go test -v -cover $(TEST_FLAGS)
test-coverpkg:
	go test -v -coverpkg=./... $(TEST_FLAGS)

# docker
build-binary:
	go build -o ${BINARY} cmd/main.go

docker-image:
	docker build -t drello-backend .

docker-container-clean:
	./scripts/rm_container.sh

docker-volume-clean:
	docker volume prune

docker-image-clean:
	./scripts/clean_images.sh

start:
	docker-compose up -d

stop:
	docker-compose down

# документация
doc-prepare:
	npm install speccy -g
	docker pull swaggerapi/swagger-ui

doc-create:
	speccy resolve docs/main.yaml -o $(API_DOC_TARGET)

doc-host:	
	docker run --name=documentation -d -p 80:8080 -e SWAGGER_JSON=/api.yaml -v $(PROJECT_DIR)/api.yaml:/api.yaml swaggerapi/swagger-ui

doc-stop:
	docker stop ${DOCUMENTATION_CONTAINER_NAME}
	docker rm ${DOCUMENTATION_CONTAINER_NAME}

.PHONY:
	start stop
