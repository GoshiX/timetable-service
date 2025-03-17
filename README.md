### Описание проекта: Бэкенд для сервиса расписания транспорта

---

### Основные функции:
   - Возможность поиска городу отправки и городу прибытия
   - Возможность получать сложные маршруты

---

### Технологии:
- Golang
- PostgreSQL
- Docker
- Docker Compose

---

### Как запустить проект:

1. **Клонировать репозиторий:**
   ```bash
   git clone https://github.com/GoshiX/timetable-service.git
   cd timetable-service
   ```

2. **Установить Docker:**

   [Гайд](https://docs.docker.com/engine/install/)

3. **Настроить конфигурацию:**

   В файле docker-compose.yml ввести порт, на котором будет запущен сервер и указать API ключ.

4. **Запустить сервер:**
     ```bash
     cd deploy
     docker-compose up --build
     ```

### Сценарии использования:

1. Получить список доступных городов:

   ```http://localhost:10200/available_dest```

2. Поиск маршрута:

   ```http://localhost:10200/route?from=Москва&to=Псков```

Все ответы даются в формате JSON.