API_DOC_TARGET=api.yaml

MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))

doc-prepare:
	npm install speccy -g
	docker pull swaggerapi/swagger-ui

doc-create:
	speccy resolve docs/main.yaml -o $(API_DOC_TARGET)

doc-host:	
	docker run -p 80:8080 -e SWAGGER_JSON=/api.yaml -v $(MAKEFILE_PATH)/../api.yaml:/api.yaml swaggerapi/swagger-ui