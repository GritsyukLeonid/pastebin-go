# Pastebin на Go с двумя микросервисами

Проект представляет собой систему Pastebin, разделённую на два микросервиса, реализованных на языке Go. Система позволяет создавать текстовые записи, генерировать короткие ссылки для их просмотра и управлять сроком хранения записей. Каждый микросервис выполняет свою задачу, взаимодействуя друг с другом через HTTP.

## Архитектура

1. **Paste Service (Сервис приёма и хранения)**:
   - Принимает текстовые записи через API.
   - Генерирует уникальные короткие хэши для записей.
   - Позволяет указать срок хранения записи (TTL).
   - Сохраняет записи сначала во временный файл, а затем в базу данных.

2. **API Service (Сервис предоставления данных)**:
   - Обрабатывает запросы на получение записей по ID (или хэшу).
   - Возвращает запись в формате JSON.
   - Проверяет TTL записей и удаляет устаревшие записи.
   - Учёт популярности записей через счётчик просмотров.

### Функционал

- **Создание текстовых записей**: Пользователи могут создать текстовые записи, которые будут сохранены в системе.
- **Генерация коротких ссылок**: Каждой записи присваивается уникальный короткий хэш.
- **Указание срока хранения**: При создании записи можно указать TTL (время жизни записи).
- **Просмотр записей**: Пользователи могут открыть ссылку и просматривать запись.
- **Удаление устаревших записей**: Записи автоматически удаляются по истечении TTL.
- **Поддержка популярности**: Каждый запрос на запись увеличивает её популярность, и наиболее популярные записи могут быть выделены.

## Архитектура и взаимодействие

Проект состоит из двух микросервисов, которые взаимодействуют следующим образом:

- **Paste Service** принимает данные и сохраняет их. После этого он генерирует хэш и сохраняет запись в базе данных.
- **API Service** отвечает на запросы на получение записи и проверяет её TTL. Если запись устарела, она удаляется.

### Требования

- Go (версия 1.18 и выше)
- Docker (для развертывания базы данных)
- База данных (например, PostgreSQL или SQLite)

### Установка и запуск

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/yourusername/pastebin-go.git
   cd pastebin-go
