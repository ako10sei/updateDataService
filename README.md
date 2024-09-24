# Обновление данных Visiology из API ЦП
Версия Golang: 1.22
## Содержание

1. [Настройка проекта](#настройка-проекта)
2. [Переменные окружения](#переменные-окружения)
3. [Запуск проекта](#запуск-проекта)


## Настройка проекта

1. **Установите Go**: Скачайте и установите версию Go 1.22 с [golang.org/dl](https://golang.org/dl/).
2. **Настройте рабочую область**: Создайте каталог для ваших проектов Go (например, `$HOME/go`). Установите переменную
   окружения `GOPATH` в указанный каталог.
3. **Склонируйте проект**: Сделайте клонирование проекта по ссылке репозитория 
```bash
git clone https://github.com/ako10sei/updateDataService.git
```
4. **Установите пакеты Go** с помощью команды для загрузки необходимых пакетов:

```bash
go mod tidy
go mod vendor
````
## Переменные окружения

1. **Создайте файл `.env`** в корне вашего проекта и добавьте в него переменные окружения, указанные в приведённом ниже
   примере:

```env
API_BASE_URL=https: API endpoint получения организаций 
API_CLIENT_SECRET=client_secret
API_CLIENT_ID=client_id

REC_API_BASE_URL=URL
REC_API_USERNAME=username
REC_API_PASSWORD=password
REC_API_VERSION=3.11

DEBUG=True // Режим отладки. При значении флага False осуществляется функционал обновления данных портала.
```

## Запуск проекта

1. Выполните команду `go run cmd/data-update/main.go` в каталоге вашего проекта, чтобы запустить приложение.
