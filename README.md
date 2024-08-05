# Обновление данных Visiology из API ЦП
Версия Golang: 1.22
## Содержание

1. [Настройка проекта](#настройка-проекта)
2. [Переменные окружения](#переменные-окружения)
3. [Запуск проекта](#запуск-проекта)
4. [Отладка](#отладка)

## Настройка проекта

1. **Установите Go**: Скачайте и установите версию Go 1.22 с [golang.org/dl](https://golang.org/dl/).
2. **Настройте рабочую область**: Создайте каталог для ваших проектов Go (например, `$HOME/go`). Установите переменную
   окружения `GOPATH` в указанный каталог.
3. **Склонируйте проект**: Сделайте клонирование проекта по ссылке репозитория 
```bash
git clone https://bitbucket.webizi.ru/scm/~sashalom666/visiologyupdategolang.git
```
4. **Установите пакеты Go** с помощью команды `go get` для загрузки необходимых пакетов:

```bash
go get github.com/joho/godotenv
go get github.com/fatih/color
````
## Переменные окружения

1. **Создайте файл `.env`** в корне вашего проекта и добавьте в него переменные окружения, указанные в приведённом ниже
   примере:

```env
DIGITAL_PROFILE_BASE_URL=https: API endpoint получения организаций ЭК
DIGITAL_PROFILE_CLIENT_SECRET= client_secret
DIGITAL_PROFILE_CLIENT_ID= client_id

VISIOLOGY_BASE_URL=https://bi.xn--33-6kcadhwnl3cfdx.xn--p1ai/
VISIOLOGY_USERNAME= username
VISIOLOGY_PASSWORD= password
VISIOLOGY_API_VERSION=3.11

DEBUG=True // Режим отладки. При значении флага False осуществляется функционал обновления данных портала.
```

## Запуск проекта

1. Выполните команду `go run cmd/data-update/data-update.go` в каталоге вашего проекта, чтобы запустить приложение.
