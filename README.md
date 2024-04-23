## О проекте

Данный проект это финальная задача в курсе от Яндекс по GoLang.
В этом проекте представлена масштабируемая система распределительного вычислителя арифметических выражений на GoLang с графическим интерфейсом на React.

Используемые системы:

* HTTP сервер работает на фреймворке [**Gin**](https://github.com/gin-gonic/gin)
* База данных используется [**Postgres**](https://www.postgresql.org/)
* Распределение задач реализовано через брокер сообщений [**RabbitMQ**](https://www.rabbitmq.com/).
* Для разворачивания системы используется [**Docker**](https://www.docker.com/).
* Live обновление на React реализовано через [**WebSocket**](https://github.com/gorilla/websocket).

Данный проект требует улучшений для использования в реальных задачах.

## Стек технологий

[![GoLang][GoLang]][GoLang-url] [![React][React.js]][React-url] [![Postgres][Postgres]][Postgres-url] [![Rabbit][Rabbit]][Rabbit-url]

## Как запустить проект?

> [!NOTE]
> В проекте предусмотрен запуск для Production и для локальной разработки. 
>
> Просто используйте `docker-compose.yml` для полного разворачивания системы или `docker-compose-dev.yml` для локальной разработки.
>
> Если у вас останутся вопросы по запуску проекта, пишите в Telegram: @solopenkovdmitriy 

### Перед установкой

Для легкого запуска системы советую использовать [Docker](https://www.docker.com/get-started/)

Если вы собираетесь выполнять локальную разработку, то для тестирования web интерфейса, написанного на React. Вам потребуется так же установить [NodeJS](https://nodejs.org)

Для разработки BackEnd части проекта на GoLang потребуется установить [GoLang](https://go.dev/learn/)

### Базовые действия

1. Скачайте репозиторий

```sh
git clone https://github.com/DmitriySolopenkov/distribCalc.rev2
```

2. Создайте новый `.env` из `.env.example` и при необходимости отредактируйте его

```sh
cp .env.example .env
```

3. Произведите те же операции для `.env` файла в папке frontend

> [!TIP]
> В файле .env нам необходимо указать только один параметр `REACT_APP_API_SERVER`<br>Укажите адрес сервера из основого .env параметр `SERVER_ADDR`

```sh
cd frontend; cp .env.example .env
```

### Запуск

4. Запустите docker-compose из главной папки

> Данный docker-compose запускает Postgres, RabbitMQ, React, Nginx и Оркестратор

```sh
docker-compose up
```

<a name="env-agent-params"></a>
### Настройка .env и параметры агента

В `.env` вам нужно отредактировать только эти поля, остальные поля нужны для продвинутой настройки
* `AGENT_TIMEOUT` - максимальное время в секундах для ожидания агента, если агент будет неактивен спустя `AGENT_TIMEOUT` секунд, то он будет удален из списка
* `AGENT_PING` - время в секундах для отправления сигнала ping от агента, если агент не отправит в течении этого времени сообщение о пинге, то его статус будет изменен на `reconnected`
* `AGENT_RESOLVE_TIME` - максимальное время в секундах для решения одной задачи, если в течении этого времени не будет получен ответ, то задача пересоздастся
* `SERVER_ADDR` - IP:PORT для запуска HTTP сервера
* `MODE` - `release` или `debug` режим запуска Gin
* `POSTGRES_HOST` & `RABBIT_HOST` - измените на `localhost` только при локальной разработке!

Параметры запуска агента
* `-agent string` - имя агента для отображения в списке серверов
* `-ping int` - время в секундах, раз в которое будет отправляться сообщение с пингов **(данный параметр должен быть меньше или равен параметру `AGENT_RESOLVE_TIME` из `.env`)**
* `-threads int` - количество потоков (goroutine) для параллельного решения задач на одном агенте
* `-wait int` - задержка решения задач для эмуляции выполнения долгих запросов
* `-debug` - включить режим отладки
* `-queue string` - имя очереди RabbitMQ с заданиями, при необходимости не менять
* `-server string` - имя очереди RabbitMQ для ответов, при необходимости не менять
* `-url string` - DSN строка для подключения RabbitMQ, при необходимости не менять


## Использование

Для использования web интерфейса перейдите по адресу
* http://127.0.0.1:8001

Документация по HTTP серверу доступна по адресу SERVER_ADDR/swagger/index.html

Во вкладке `СЕРВЕРА` вы можете увидеть список всех запущенных агентов.


[GoLang]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[GoLang-url]: https://go.dev/
[React.js]: https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB
[React-url]: https://reactjs.org/
[Postgres]: https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white
[Postgres-url]: https://www.postgresql.org/
[Rabbit]: https://img.shields.io/badge/rabbitmq-%23FF6600.svg?&style=for-the-badge&logo=rabbitmq&logoColor=white
[Rabbit-url]: https://www.rabbitmq.com/
