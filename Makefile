# Не писал мэйкфайлы лет сто
API_DOC_TARGET=api.yaml

# Проверки на 

# TODO добавить проверки на наличие данных приколюх, если они уже установлены, то ничего не делать
doc-prepare:
	npm install speccy -g
	docker pull swaggerapi/swagger-ui

doc-create:
	speccy resolve docs/main.yaml -o $(API_DOC_TARGET) 

# Тупая локальная херовина
doc-host:	
	docker run -p 80:8080 -e SWAGGER_JSON=/api.yaml -v /home/q273/Techno/secondSem/TrelloBackend/2020_1_SIBIRSKAYA_KORONA/api.yaml:/api.yaml swaggerapi/swagger-ui