# url-storage

Микросервис хранения url-alias для сервиса [url-shortener](https://github.com/quibex/url-storage)

## API
Умеет 2 ручки:
 - SetUrl(url, alias)
 - GetUrl(alias) url

grpc: [github.com/quibex/url-storage-api](https://github.com/quibex/url-storage-api)

## Стек
- Redis: используется для сквозного кэширования
- Postgres: для хранения данных
