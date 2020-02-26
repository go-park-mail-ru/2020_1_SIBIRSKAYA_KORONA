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
  * 409 - Conflict (пользователь с таким логином уже существует)
  
 # url: /login (POST)
 
  * 200 - OK (успешный запрос)
  * 308 - PermanentRedirect (уже залогинен, редирект на главную)
  * 400 - BadRequest (неверный запрос)
  * 404 - NotFound (нет пользвателя с указанным логином)
  * 412 - PreconditionFailed (неверный пароль)
  
 # url: /logout (POST) // TODO : DELETE
 
   * 200 - OK (успешный запрос)
   * 303 - SeeOther (смотреть другое, редирект на логин)

 # url: /profile?nickname=you_nick (GET)
 
  * 200 - OK (успешный запрос)
  * 400 - BadRequest (неверный запрос)
  * 404 - NotFound (нет пользвателя с указанным логином)
  
   # url: /profile (GET) // TODO: (PUT)
