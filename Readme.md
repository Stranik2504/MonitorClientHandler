
# ClientHandler

Клиенское приложение на Go для работы с MonitorHandler (https://github.com/Stranik2504/MonitorHandler).

## Зависимости

- Go 1.18 или новее
- Внешние зависимости, указанные в `go.mod`

Для установки зависимостей выполните:
```sh
go mod download
```

## Запуск

Для запуска сервера используйте:
```sh
go run .
```
или соберите бинарный файл:
```sh
go build -o clienthandler
./clienthandler
```

## Сборка документации

Для генерации документации используйте стандартную команду Go:
```sh
go doc ./...
```
или для HTML-документации (требуется godoc):
```sh
godoc -http=:6060
```
Документация будет доступна по адресу [http://localhost:6060/pkg/](http://localhost:6060/pkg/).
