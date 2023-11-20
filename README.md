# channeler

## Как запустить

Запускать надо `cmd/channeler/main.go`

### env-переменыне

* `TELEGRAM_APP_ID` - можно получить на https://my.telegram.org/apps
* `TELEGRAM_APP_HASH` - можно получить на https://my.telegram.org/apps
* `SESSION_STORAGE_PATH` - где хранить сессию, без этого надо будет логиниться каждый запуск приложения
* `TELEGRAM_CHANNEL_ID` - ID канала
* `TELEGRAM_CHANNEL_ACCESS_HASH` - AccessHash канала
* `YDB_CONNECTION_STRING` - DSN для подключения к YDB
* `TABLE_PATH_PREFIX` - tablePathPrefix, где хранить таблицы в YDB

### Авторизация в телеграме

При первом запуске попросит авторизоваться прямо в терминале, сессию сохранит в `SESSION_STORAGE_PATH`, если он указан

### Авторизация в YDB

Тут используется [ydb-go-sdk-auth-environ](https://github.com/ydb-platform/ydb-go-sdk-auth-environ)

Можно ничего не делать, если запускать на виртуалке в Облаке и привязать к ней сервисный аккаунт с правами доступа к базе данных

## Возможные проблемы

* Через 7 суток пропала связь до телеграма, помогло добавить `RuntimeMaxSec=1d` в systemd-сервис

## Ссылки

* Статья на Хабре: https://habr.com/ru/companies/yandex_cloud_and_infra/articles/774104/
* Мой телеграм-канал: http://t.me/impressionableracoonchannel
