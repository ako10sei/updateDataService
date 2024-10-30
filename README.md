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
git clone https://bitbucket.webizi.ru/scm/~sashalom666/visiologyupdategolang.git
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
где * - конфиденциальные параметры

DIGITAL_PROFILE_BASE_URL=https://xn--n1abf.xn--33-6kcadhwnl3cfdx.xn--p1ai/digital_profile/api/v1.0.0/
DIGITAL_PROFILE_CLIENT_SECRET= ******
DIGITAL_PROFILE_CLIENT_ID= ******

VISIOLOGY_BASE_URL=https://bi.xn--33-6kcadhwnl3cfdx.xn--p1ai/
VISIOLOGY_USERNAME= *****
VISIOLOGY_PASSWORD= *****
VISIOLOGY_API_VERSION=3.11

ENVIRONMENT=local

DEBUG=True // Режим отладки. При значении флага False осуществляется функционал обновления данных портала.
```

## Запуск проекта

1. Выполните команду `go run cmd/data-update/main.go` в каталоге вашего проекта, чтобы запустить приложение.
