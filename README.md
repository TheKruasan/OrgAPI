# Organization Structure API

## Описание

REST API для управления организационной структурой:
- подразделения (иерархия)
- сотрудники внутри подразделений
- перемещение без циклов
- удаление (cascade / reassign)

---

## Запуск проекта

### Через Docker

docker compose up --build

---

### API будет доступно:

http://localhost:8080

---

## База данных

PostgreSQL поднимается автоматически через docker-compose.

---

## API

### Создать подразделение
POST /departments

---

### Получить подразделение (дерево)
GET /departments/{id}?depth=1&include_employees=true

- depth — глубина дерева (max 5)
- include_employees — включить сотрудников

---

### Обновить подразделение
PATCH /departments/{id}

---

### Удалить подразделение
DELETE /departments/{id}?mode=cascade  
DELETE /departments/{id}?mode=reassign&reassign_to_department_id=2

- cascade — удалить всё дерево
- reassign — перенести сотрудников и удалить отдел

---

### Создать сотрудника
POST /departments/{id}/employees

---

## Кратко о логике

- Иерархия подразделений (дерево)
- Запрет циклов
- До 5 уровней вложенности
- Сотрудники привязаны к department
- Cascade / reassign удаление