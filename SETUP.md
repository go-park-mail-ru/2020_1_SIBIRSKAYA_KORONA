# ИНСТРУКЦИЯ ПО ЗАПУСКУ БЭКЕНД ЧАСТИ ПРИЛОЖЕНИЯ 

### Если запускаем всё на локалке

1) Cчитаем, что postgres и memcached установлены

2) Поднимаем их
    * service memcached start
    * service postgresql start

3) Настраиваем конфиг
    * Выставляем в database - host: localhost
    * Выставляем в memcached - host: localhost

4) bash scripts/run.sh

5) go run cmd/main.go --config config.yaml

### Работа с приложением внутри docker-контейнеров

1) Перед всем этим нужно отключить локальные memcached и postgres, потому что они занимают порты
    * service memcached stop
    * service postgresql stop

2) Создаем образ для нашего бэкенд-контейнера: make docker-image

3) Настраиваем конфиг
    * Выставляем в database - host: postgres
    * Выставляем в memcached - host: memcached


4) Запускаем наше приложение (оно состоит из трех-контейнеров) через docker-compose: make start

5) Останавливаем приложение: make stop

### Чистка

Посмотреть все образы, которые есть сейчас в докере: docker images

Удалить образ по id: docker rmi {id}

Удалить все фиговые <none> образы (это надо будет чинить): docker rmi $(docker images -f "dangling=true" -q)

Посмотреть контейнеры (действующие): docker ps

Посмотреть все контейнеры (упавшие в том числе): docker ps -a

Вывевсти только их id: docker ps -a -q

Остановить все контейнеры: docker stop $(docker ps -a -q)

Удалить все контейнеры вместе с их разделами (volume): docker rm -v $(docker ps -a -q)
