<div align="center">
  <h1>📚 StudySpot</h1>
  <p><strong>Сервис поиска мероприятий для студентов на Go с современным веб-интерфейсом</strong></p>
  <p>
    <img src="https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go" alt="Go 1.22">
    <img src="https://img.shields.io/badge/PostgreSQL-16-336791?style=flat-square&logo=postgresql" alt="PostgreSQL">
    <img src="https://img.shields.io/badge/Redis-7.4-DC382D?style=flat-square&logo=redis" alt="Redis">
    <img src="https://img.shields.io/badge/Docker-✓-2496ED?style=flat-square&logo=docker" alt="Docker">
    <img src="https://img.shields.io/badge/JWT-✓-000000?style=flat-square&logo=jsonwebtokens" alt="JWT">
    <img src="https://img.shields.io/badge/license-MIT-green?style=flat-square" alt="License MIT">
    <img src="https://img.shields.io/badge/HTML5-CSS3-orange?style=flat-square&logo=html5" alt="HTML5/CSS3">
    <img src="https://img.shields.io/badge/JavaScript-Vanilla-F7DF1E?style=flat-square&logo=javascript" alt="JavaScript">
    <img src="https://img.shields.io/badge/Nginx-1.27-009639?style=flat-square&logo=nginx" alt="Nginx">
    <img src="https://img.shields.io/badge/Gin-1.10-00ADD8?style=flat-square&logo=go" alt="Gin">
  </p>
</div>

## ✨ О проекте

**StudySpot** — это полностью контейнеризированный сервис для поиска студенческих мероприятий, написанный на современном стеке технологий. Проект создан с упором на производительность, масштабируемость и удобство использования.

### Основные возможности:
- 🎓 **Единая платформа** для всех мероприятий университета
- 🔍 **Быстрый поиск** по названию с полнотекстовым поиском PostgreSQL
- 📂 **Фильтрация** по категориям (Олимпиады, Хакатоны, Курсы, Стажировки)
- ⚡ **Кэширование** в Redis для мгновенных ответов
- 🔐 **Безопасная аутентификация** через JWT с ролевой моделью
- 👑 **Администрирование** — создание, редактирование и удаление мероприятий
- 🎨 **Современный веб-интерфейс** с тёмной темой и адаптивным дизайном
- 🐳 **Docker-first подход** — весь стек поднимается одной командой
- 🗄️ **Автоматические миграции** при запуске контейнеров

## 🛠 Стек технологий

| Компонент | Технология |
|-----------|------------|
| **Язык (Backend)** | [Go 1.22](https://go.dev/) |
| **Веб-фреймворк** | [Gin 1.10](https://gin-gonic.com/) |
| **База данных** | [PostgreSQL 16](https://www.postgresql.org/) |
| **Кэш** | [Redis 7.4](https://redis.io/) |
| **Аутентификация** | [JWT](https://jwt.io/) |
| **Веб-сервер** | [Nginx](https://nginx.org/) |
| **Фронтенд** | HTML5, CSS3, JavaScript (Vanilla) |
| **Контейнеризация** | [Docker](https://www.docker.com/) + [Docker Compose](https://docs.docker.com/compose/) |

## 🚀 Быстрый старт

### Предварительные требования
- Установленные [Docker](https://docs.docker.com/get-docker/) и [Docker Compose](https://docs.docker.com/compose/install/)
- (Опционально) [Go 1.22+](https://go.dev/dl/) для локальной разработки

### Установка и запуск

1. **Клонируйте репозиторий**
   ```bash
   git clone https://github.com/yourusername/studyspot.git
   cd studyspot
   ```

2. **Настройте переменные окружения**
   
   Скопируйте файл с примером конфигурации:
   ```bash
   cp .env.example .env
   ```
   
   Минимально необходимые настройки:
   ```env
   # Сервер
   SERVER_PORT=8080
   GIN_MODE=release
   
   # PostgreSQL
   DB_HOST=postgres
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=studyspot
   DB_SSL_MODE=disable
   
   # Redis
   REDIS_HOST=redis
   REDIS_PORT=6379
   REDIS_PASSWORD=
   REDIS_DB=0
   
   # JWT
   JWT_SECRET=your-super-secret-key-change-in-production
   JWT_EXPIRE_HOURS=24
   ```

3. **Запустите все сервисы**
   ```bash
   docker compose up -d --build
   ```

4. **Проверьте работу**
   - Веб-интерфейс: http://localhost
   - Health check: http://localhost:8080/health
   - API документация: http://localhost:8080/api/events

### Тестовый аккаунт администратора

```
Email: admin@studyspot.com
Password: admin123
```

## 📖 Использование

### 🌐 Веб-интерфейс

Главная страница сервиса предлагает интуитивно понятный интерфейс для поиска мероприятий.

| Функция | Описание |
|---------|----------|
| **Поиск** | Введите название мероприятия в строку поиска |
| **Фильтр** | Выберите категорию из выпадающего списка |
| **Просмотр** | Карточки мероприятий с названием, описанием, датой и местом |
| **Авторизация** | Вход и регистрация через модальное окно |
| **Создание** | Администраторы могут создавать новые мероприятия |

### 🔗 API Endpoints

**Публичные маршруты (без авторизации):**

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/auth/register` | Регистрация пользователя |
| POST | `/api/auth/login` | Авторизация (JWT) |

**Защищённые маршруты (JWT):**

| Метод | Путь | Описание | Роль |
|-------|------|----------|------|
| GET | `/api/events` | Список мероприятий | user/admin |
| GET | `/api/events/search?q=...` | Поиск по названию | user/admin |
| GET | `/api/events/search?category=...` | Фильтр по категории | user/admin |
| GET | `/api/events/:id` | Карточка мероприятия | user/admin |
| GET | `/api/categories` | Список категорий | user/admin |

**Административные маршруты (только admin):**

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/events` | Создать мероприятие |
| PUT | `/api/events/:id` | Обновить мероприятие |
| DELETE | `/api/events/:id` | Удалить мероприятие |
| POST | `/api/categories` | Создать категорию |
| PUT | `/api/categories/:id` | Обновить категорию |
| DELETE | `/api/categories/:id` | Удалить категорию |

### 📝 Примеры API запросов

**Регистрация:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"student@example.com","password":"password123"}'
```

**Вход (получение JWT):**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@studyspot.com","password":"admin123"}'
```

**Получение списка мероприятий:**
```bash
curl http://localhost:8080/api/events \
  -H "Authorization: Bearer <your_jwt_token>"
```

**Поиск мероприятий:**
```bash
curl "http://localhost:8080/api/events/search?q=хакатон" \
  -H "Authorization: Bearer <your_jwt_token>"
```

**Создание мероприятия (админ):**
```bash
curl -X POST http://localhost:8080/api/events \
  -H "Authorization: Bearer <your_jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Хакатон 2026",
    "description": "Грандиозный хакатон для студентов",
    "category_id": "UUID_категории",
    "date": "2026-07-15T10:00:00Z",
    "location": "Главный корпус, ауд. 101"
  }'
```

## 📁 Структура проекта

```
studyspot/
├── cmd/                              # Точка входа
│   └── app/
│       └── main.go                   # Запуск сервера
├── internal/                         # Внутренняя логика
│   ├── cache/                        # Redis клиент
│   │   └── redis.go
│   ├── config/                       # Конфигурация
│   │   └── config.go
│   ├── domain/                       # Модели данных
│   │   ├── user.go
│   │   ├── event.go
│   │   └── category.go
│   ├── handler/                      # HTTP-обработчики
│   │   ├── auth_handler.go
│   │   ├── event_handler.go
│   │   └── category_handler.go
│   ├── middleware/                   # JWT и CORS
│   │   └── auth.go
│   ├── repository/                   # Работа с БД
│   │   ├── user_repo.go
│   │   ├── event_repo.go
│   │   └── category_repo.go
│   └── service/                      # Бизнес-логика
│       ├── auth_service.go
│       ├── event_service.go
│       └── category_service.go
├── pkg/                              # Вспомогательные пакеты
│   ├── jwt/                          # JWT утилиты
│   │   └── jwt.go
│   ├── password/                     # Хеширование паролей
│   │   └── password.go
│   └── response/                     # Стандартизация ответов
│       └── response.go
├── frontend/                         # Фронтенд
│   └── index.html                    # Весь фронтенд в одном файле
├── migrations/                       # SQL-миграции
│   └── 001_init_schema.sql
├── docker/
│   └── Dockerfile                    # Dockerfile для бэкенда
├── docker-compose.yml                # Оркестрация всех сервисов
├── .env.example                      # Пример конфигурации
├── go.mod                            # Зависимости Go
├── go.sum                            # Хеши зависимостей
└── README.md                         # Документация
```

## 🗄️ Модели данных

### User (пользователь)
```sql
id         UUID PRIMARY KEY
email      VARCHAR(255) UNIQUE NOT NULL
password   VARCHAR(255) NOT NULL
role       VARCHAR(20) DEFAULT 'user'
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

### Category (категория)
```sql
id         UUID PRIMARY KEY
name       VARCHAR(100) UNIQUE NOT NULL
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

### Event (мероприятие)
```sql
id          UUID PRIMARY KEY
title       VARCHAR(255) NOT NULL
description TEXT
category_id UUID REFERENCES categories(id)
date        TIMESTAMP NOT NULL
location    VARCHAR(255)
created_by  UUID REFERENCES users(id)
created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

## 🎨 Особенности интерфейса

- **Тёмная тема** для комфортной работы в любое время
- **Адаптивный дизайн** для всех устройств
- **Плавные анимации** при взаимодействии
- **Карточки мероприятий** с всей необходимой информацией
- **Модальные окна** для создания и редактирования
- **Всплывающие уведомления** (toast-сообщения)

## 🗄️ Управление миграциями

Миграции применяются автоматически при первом запуске контейнера PostgreSQL.

### Ручное применение миграций

```bash
# Подключиться к БД
docker exec -it studyspot-db psql -U postgres -d studyspot

# Применить миграцию
\i /docker-entrypoint-initdb.d/001_init_schema.sql
```

## 🤝 Вклад в проект

Будем рады вашим идеям и улучшениям! Чтобы внести вклад:

1. Форкните репозиторий
2. Создайте ветку для фичи (`git checkout -b feature/amazing-feature`)
3. Закоммитьте изменения (`git commit -m '✨ Add some amazing feature'`)
4. Запушьте ветку (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📄 Лицензия

Проект распространяется под лицензией MIT. Подробности в файле LICENSE.

---

<div align="center">
  <p>🌟 Если вам понравился проект, поставьте звезду на GitHub!</p>
  <p>📧 По всем вопросам: your@email.com</p>
</div>
