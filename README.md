# Org Structure API

REST API для управления организационной структурой (подразделения и сотрудники).

---

## Стек
- Go (net/http)
- PostgreSQL
- GORM
- Goose (миграции)
- Docker / docker-compose

---

## Запуск проекта

docker compose up --build

---

## Что происходит при запуске

1. Поднимается PostgreSQL  
2. Дожидаемся готовности БД (healthcheck)  
3. Применяются миграции (goose)  
4. Запускается приложение  

---

## API

### Departments

Создать подразделение  
POST /departments/

Получить подразделение (дерево)  
GET /departments/{id}?depth=1&include_employees=true

Обновить подразделение  
PATCH /departments/{id}

Удалить подразделение  
DELETE /departments/{id}?mode=cascade|reassign

---

### Employees

Создать сотрудника  
POST /departments/{id}/employees/

---

## Бизнес-логика

- Иерархия подразделений (parent_id → self reference)  
- Запрещены циклы в дереве  
- Уникальность имени внутри одного parent  
- Удаление:
  - cascade — удаление всего дерева
  - reassign — перенос сотрудников  

---

## Миграции

Миграции применяются автоматически при старте контейнера migrate.

---

## Переменные окружения

DB_HOST=postgres  
DB_PORT=5432  
DB_USER=postgres  
DB_PASSWORD=postgres  
DB_NAME=org_db  

.env файл оставлен для удобства запуска и тестирования 
---

## Остановка

docker compose down

---

## Структура проекта

cmd/  
internal/  
  handler/  
  service/  
  repository/  
  dto/  
  models/  
migrations/  
Dockerfile  
docker-compose.yml  
migration.sh  