BINARY=drello_binary
API_DOC_TARGET=api.yaml
PROJECT_DIR := ${CURDIR}

# тесты
test-cover:
	go test -v -cover -covermode=atomic ./... --test-config="$(PROJECT_DIR)/config.yaml"

test-coverpkg:
	go test -v -coverpkg=./... -covermode=atomic ./... --test-config="$(PROJECT_DIR)/config.yaml"

# docker
build-binary:
	go build -o ${BINARY} cmd/main.go

docker-image:
	docker build -t drello-backend .

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
	echo ${PROJECT_DIR}

doc-host:	
	docker run -p 80:8080 -e SWAGGER_JSON=/api.yaml -v $(PROJECT_DIR)/api.yaml:/api.yaml swaggerapi/swagger-ui

.PHONY:
	start stop
