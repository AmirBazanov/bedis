[![Test Resp Writer Reader](https://github.com/AmirBazanov/bedis/actions/workflows/go_unit_read_write_resp_fuzz.yml/badge.svg)](https://github.com/AmirBazanov/bedis/actions/workflows/go_unit_read_write_resp_fuzz.yml)
[![Test Server](https://github.com/AmirBazanov/bedis/actions/workflows/go_test_server.yml/badge.svg)](https://github.com/AmirBazanov/bedis/actions/workflows/go_test_server.yml)
[![Test Storage](https://github.com/AmirBazanov/bedis/actions/workflows/go_unit_storage.yml/badge.svg)](https://github.com/AmirBazanov/bedis/actions/workflows/go_unit_storage.yml)

# Mini Redis

Минималистичная реализация Redis на Go, предназначенная для обучения и экспериментов с базами данных в памяти.

## Описание

Mini Redis — это упрощённый in-memory key-value store, вдохновлённый Redis.  
Проект поддерживает базовые команды для работы с ключами и значениями и может использоваться для понимания принципов работы высокопроизводительных баз данных в памяти.

Особенности проекта:

- Хранение данных в памяти
- Поддержка основных команд: `SET`, `GET`, `DEL`, `EXISTS`
- Простой TCP-сервер для работы с клиентами
- Лёгкий и понятный код для обучения

## Установка

Склонируйте репозиторий и соберите проект:

```bash
git clone https://github.com/amirbazanov/bedis.git
cd bedis
go build -o bedis
```

По умолчанию сервер слушает порт 6379.
