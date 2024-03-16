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

## Регенерация схемы GraphQL

Для регенерации схемы достаточно перейти в директорию `/api/v1/` и запустить там:

```
$ go run github.com/99designs/gqlgen generate
```

## Регенерация protobuf

Переходим `/src/` и запускаем:

```
$ protoc --experimental_allow_proto3_optional --go_out=./protochats --go_opt=paths=source_relative --go-grpc_out=./protochats --go-grpc_opt=paths=source_relative chats.proto
```
