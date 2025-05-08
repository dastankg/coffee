
## 🔍 Обзор проекта

Проект представляет собой REST API сервис для управления базой данных кофе. Сервис позволяет создавать, получать, обновлять и удалять записи о различных сортах кофе, а также включает систему авторизации пользователей и отправки уведомлений.

---
## 💻 Технологии

- Go 1.24
- PostgreSQL 16
- Docker & Docker Compose
- JWT для авторизации
- Swagger для документации API
- GORM как ORM
- QR-код генератор
- SMTP для email-уведомлений

---
## ⚙️ Запуск проекта

### 📌 Предварительные требования

- Установленный Docker и Docker Compose
- Go 

```bash
    git clone https://github.com/dastankg/coffee
    cd coffee
```
---

### 🔧 Настройка переменных окружения

Создайте файл `.env` в корне проекта со следующим содержимым:

```env
# База данных
DB_HOST=postgres
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=coffee_db
DATABASE_URL=postgresql://your_username:your_password@postgres:5432/coffee_db

# JWT токены
TOKEN=your_secure_token_secret

# SMTP для отправки email
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_EMAIL=your_email@example.com
SMTP_PASSWORD=your_email_password
```

##  Сборка и запуск контейнеров
docker-compose up --build -d


## Документация API доступна через Swagger по адресу:

http://localhost:8081/docs/