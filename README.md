# backdev_go

backdev_go - это пример RESTful HTTP API сервиса для аутентификации на стеке Go + JWT + PSQL.

Вся коммуникация сделана на POST запросах, где данные - это JSON в теле запросов. Почему не GET + заголовки? Потому что не было такого требования, и переделать в случае чего не сложно.

## Запуск

### docker-compose

```yml
version: '3.8'

services:
  backdev:
    image: ussuratoncachi/backdev_go:0.2
    container_name: backdev_go
    ports:
      - "9000:9000"
    depends_on:
      - postgres
    environment:
      - BACKDEV_SECRET=My Very Precious Secret Phrase
      - BACKDEV_DB_TYPE=postgresql
      - BACKDEV_LISTEN_IP=0.0.0.0:9000

      - BACKDEV_PSQL_HOST=postgres
      - BACKDEV_PSQL_PORT=5432
      - BACKDEV_PSQL_USER=admin
      - BACKDEV_PSQL_PASSWORD=admin
      - BACKDEV_PSQL_DB_NAME=backdev_db

      #- BACKDEV_SMTP_HOST=smtp.gmail.com
      #- BACKDEV_SMTP_PORT=587
      - BACKDEV_SMTP_USER=<Your SMTP Email>
      - BACKDEV_SMTP_FROM_EMAIL=<Your SMTP Email>
      - BACKDEV_SMTP_PASSWORD=<Your SMTP password>
    restart: always
    networks:
      - app-network

  postgres:
    image: postgres:15
    container_name: postgres
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: admin
      POSTGRES_DB: backdev_db
    volumes:
      - ./postgres-data:/var/lib/postgresql/data
    networks:
      - app-network

networks:
  app-network:
    driver: bridge
```

### Альтернативно: вручную собрать программу и запустить

```bash
$ go build .
$ nano config.toml # Не забудтье настроить конфиг под себя!
$ ./backdev_go config.toml 
```

## Простой UI

На главном руте (`/`) доступна HTML страница с интерфейсом для выполнения запросов на каждый эндпоинт.

## Эндпоинты

### POST `/authorize`. 
Выдает новую пару JWT + Refresh токенов.

Принимает JSON: `user_uuid`, `user_email`. 

Возвращает JSON: `access_token`, `refresh_token`.

### POST `/refresh`. 

Обновляет JWT токен. Работает только с тем Refresh токеном, который был выдан к этому Access токену. Работает только один раз. При попытке использования одной пары токенов несколько раз - выдаст ошибку.

При изменении IP адреса пользователя - он отправляет предупреждение на Email. Для предупреждения нужно настроить SMTP сервер в `config.toml`.

Принимает JSON: `access_token`, `refresh_token`, `emulate_ip` (опционально).

Возвращает JSON: `access_token`, `refresh_token`.

### POST `/validate`

Проверяет корректность JWT токена. 

Принимает JSON: `access_token`.

Возвращает - статус 200 если токен корректный, статус 4XX - если некорректный.

