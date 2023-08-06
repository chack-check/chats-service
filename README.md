# Chack Check Chats Service

Сервис чатов для chack check, написанный на Go

## Запуск

Перед запуском всех сервисов chack check локально, нужно создать docker network, чтобы
сервисы могли общаться между собой:

```
$ docker network create chack-check-network
```

> Network нужно создать только один раз

После этого можно запустить приложение через `make`:

```
$ make dev
```

## GraphiQL

После запуска локально, сервис будет доступен по адресу http://localhost:8001/api/v1/chats
