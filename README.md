# SIBIRSKAYA KORONA

## Репозиторий бекенда

## Команда

## Менторы

## Ссылка на проект

## Репозиторий фронтенда

## Статусы ответа

# url: /join (POST)

  * 200 - OK (успешный запрос)
  * 400 - BadRequest (неверный запрос)
  * 409 - Conflict (пользователь с таким ником уже существует)
  
# url: /login (POST)
 
  * 200 - OK (успешный запрос)
  * 308 - PermanentRedirect (уже залогинен, редирект на главную)
  * 400 - BadRequest (неверный запрос)
  * 404 - NotFound (нет пользвателя с указанным ником)
  * 412 - PreconditionFailed (неверный пароль)
  
# url: /logout (DELETE)
 
  * 200 - OK (успешный запрос)
  * 303 - SeeOther (смотреть другое, редирект на логин)

# url: /profile?nickname=you_nick (GET)

  nickname в query string - просмотр страницы без возможности изменения

  * 200 - OK (успешный запрос)
  * 303 - SeeOther (не авторизован, случай без query string)
  * 400 - BadRequest (неверный запрос)
  * 404 - NotFound (нет пользвателя с указанным ником)
  
# url: /profile (PUT) 
