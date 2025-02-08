# pastebin-go

# 📝 Pastebin на Go  

Pastebin-сервис на Go с уникальными короткими хэшами для URL, поддержкой популярных записей и адаптацией к различной активности пользователей.  

## 🔧 Функционал  
- Пользователь может создать текстовый блок и загрузить его в систему.  
- Генерация уникальных и коротких хэшей для ссылки на посты.  
- Другие пользователи могут открыть ссылку и увидеть этот же текст.  
- Возможность указать срок хранения записи (TTL).  
- Автоматическое удаление записей после истечения времени.  
- Некоторые записи будут более популярными и часто просматриваемыми.  
- Один пользователь может создавать посты чаще других.  

---

## 🏗 Архитектура  

![image](https://github.com/user-attachments/assets/ace7460e-b420-4bc0-bca8-8f9f42641eca)


### **1. Paste Service (Сервис приёма и хранения)**  
- Принимает текстовые записи через API.  
- Позволяет указать TTL (время жизни записи).  
- Сохраняет данные временно в файл, затем в базу.  
- Генерирует короткие уникальные хэши для URL.  
- Отправляет сообщение в брокер для последующей обработки.  

### **2. API Service (Сервис предоставления данных)**  
- Обрабатывает запросы на получение записей по ID.  
- Возвращает paste в формате JSON или Protobuf.  
- Проверяет TTL записей и удаляет устаревшие.  
- Поддерживает механизмы учёта популярности записей (например, через счетчик просмотров).  

---

## 🚀 Установка и запуск  

### **1. Клонирование репозитория**  
```bash
git clone https://github.com/GritsyukLeonid/pastebin-go.git
cd pastebin-go
```

### **2. Запуск с Docker Compose**  
```bash
docker-compose up --build
```

### **3. Запуск вручную**  
```bash
go run ./cmd/paste_service/main.go
go run ./cmd/api_service/main.go
```

---

## ⚙ Конфигурация  

### **.env (пример)**  
```ini
DB_HOST=localhost
DB_PORT=5432
DB_USER=youruser
DB_PASS=yourpassword
DB_NAME=pastebin
BROKER_URL=nats://localhost:4222
DEFAULT_TTL=3600  # Время жизни записей в секундах
```

---

## 📖 API  

### **Создание paste**  
```http
POST /api/v1/paste
Content-Type: application/json

{
  "content": "Привет, мир!",
  "expiration": 86400  # Время хранения (в секундах)
}
```
**Ответ:**  
```json
{
  "id": "abc123",  # Уникальный короткий хэш
  "link": "http://localhost:8080/api/v1/paste/abc123",
  "message": "Paste сохранён"
}
```

### **Получение paste по ID**  
```http
GET /api/v1/paste/abc123
Accept: application/json
```
**Ответ:**  
```json
{
  "id": "abc123",
  "content": "Привет, мир!",
  "created_at": "2025-02-08T12:00:00Z",
  "expires_at": "2025-02-09T12:00:00Z",
  "views": 15  # Количество просмотров
}
```

### **Ошибка при истекшем сроке**  
```json
{
  "error": "Paste not found or expired"
}
```

---

## 🤝 Контрибьюция  
1. Форкните репозиторий  
2. Создайте новую ветку (`git checkout -b feature-branch`)  
3. Внесите изменения и зафиксируйте (`git commit -m "Добавлен новый функционал"`)  
4. Отправьте изменения (`git push origin feature-branch`)  
5. Откройте Pull Request  

---

## 📜 Лицензия  
MIT License. Свободно используйте и модифицируйте. 
