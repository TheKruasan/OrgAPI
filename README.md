# Organization Structure API

## Описание

REST API для управления организационной структурой:

- подразделения с иерархией (дерево)
- сотрудники внутри подразделений
- перемещение подразделений без циклов
- удаление подразделений (cascade / reassign)

---

## Запуск проекта

### 1. Поднять проект

make up

---

### 2. Остановить проект

make down

---

## Что делает make up

Команда:

make up

выполняет:

1. запуск PostgreSQL через docker compose  
2. ожидание готовности базы  
3. установка goose (если не установлен)  
4. накатывание миграций  
5. запуск приложения

---

## API

### Создать подразделение

POST /departments

Body:
- name (string)
- parent_id (int, optional)

---

### Получить подразделение (дерево)

GET /departments/{id}?depth=1&include_employees=true

Query параметры:
- depth (int, max 5) — глубина дерева
- include_employees (bool) — включать сотрудников

---

### Обновить подразделение

PATCH /departments/{id}

Body:
- name (string, optional)
- parent_id (int, optional)

---

### Удалить подразделение

DELETE /departments/{id}

Query:
- mode=cascade — удалить всё дерево
- mode=reassign — перенести сотрудников

Дополнительно:
- reassign_to_department_id (обязателен при mode=reassign)

---

### Создать сотрудника

POST /departments/{id}/employees

Body:
- full_name (string)
- position (string)
- hired_at (date, optional)

---

## Бизнес-логика

- подразделения образуют дерево через parent_id
- запрещено создавать циклы в структуре
- имя подразделения уникально в рамках одного родителя
- глубина дерева ограничена (до 5 уровней)
- сотрудники привязаны к подразделениям
- поддерживается cascade и reassign удаление

---

## Требования

- Go
- net/http
- GORM
- PostgreSQL
- goose (миграции)
- Docker / docker-compose

---

## Примечание

Все миграции выполняются автоматически через makefile при запуске проекта.