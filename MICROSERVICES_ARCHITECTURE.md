# Микросервисная архитектура блога

## Архитектура

```
┌─────────────┐
│   Frontend  │
└──────┬──────┘
       │
       ↓
┌─────────────────────────┐
│     API Gateway         │  (REST → gRPC)
│  Port: 8080             │
└────┬────────────────────┘
     │
     ├──→ ┌────────────────────┐
     │    │  Auth Service      │  gRPC :50051
     │    │  - JWT токены      │
     │    │  - Пользователи    │
     │    └────────────────────┘
     │
     ├──→ ┌────────────────────┐
     │    │  Article Service   │  gRPC :50052
     │    │  - CRUD статей     │
     │    │  - Public/Private/ │
     │    │    Link visibility │
     │    └────────┬───────────┘
     │             │
     └──→ ┌───────┴────────────┐
          │  Stats Service     │  gRPC :50053
          │  - Просмотры       │
          │  - Лайки           │
          └────────────────────┘
                   ↑
                   │
          ┌────────┴─────────┐
          │   RabbitMQ       │
          │  Events:         │
          │  - article.viewed│
          │  - article.liked │
          │  - article.unliked
          └──────────────────┘
```

## Сервисы

### 1. Auth Service (50051)
- **Ответственность**: Аутентификация и управление пользователями
- **База данных**: PostgreSQL (таблица users)
- **gRPC методы**:
  - Register
  - Login
  - ValidateToken
  - GetUserByID
  - GetUserByEmail

### 2. Article Service (50052)
- **Ответственность**: Управление статьями с поддержкой visibility
- **База данных**: PostgreSQL (таблица articles)
- **RabbitMQ**: Публикует события (article.created, article.viewed)
- **Типы видимости**:
  - **PUBLIC**: Доступна всем пользователям
  - **PRIVATE**: Только автор
  - **LINK**: Доступна по уникальной ссылке (access_token)
- **gRPC методы**:
  - CreateArticle
  - GetArticle (с проверкой доступа)
  - UpdateArticle
  - DeleteArticle
  - ListArticles (только публичные)
  - GetArticlesByUser (фильтрация по правам)
  - CheckArticleAccess

### 3. Stats Service (50053)
- **Ответственность**: Статистика и аналитика
- **База данных**: PostgreSQL (таблицы article_views, article_likes)
- **RabbitMQ**: Подписка на события просмотров и лайков
- **gRPC методы**:
  - RecordView
  - RecordLike
  - RemoveLike
  - GetArticleStats
  - GetUserLikeStatus
  - GetArticlesWithStats

### 4. API Gateway (8080)
- **Ответственность**: REST API для фронтенда, маршрутизация к микросервисам
- **Endpoints**:
  - `POST /api/auth/register`
  - `POST /api/auth/login`
  - `GET /api/articles` - публичные статьи
  - `GET /api/articles/:id?access_token=xxx` - статья (с токеном для LINK)
  - `POST /api/articles` - создать статью (требует auth)
  - `PUT /api/articles/:id` - обновить статью
  - `DELETE /api/articles/:id` - удалить статью
  - `POST /api/articles/:id/like` - лайк/анлайк

## База данных

### Изменения в схеме:

```sql
-- Добавить в таблицу articles
ALTER TABLE articles ADD COLUMN visibility VARCHAR(20) DEFAULT 'public';
ALTER TABLE articles ADD COLUMN access_token VARCHAR(255);

CREATE INDEX idx_articles_visibility ON articles(visibility);
CREATE INDEX idx_articles_access_token ON articles(access_token);
```

## RabbitMQ События

### Exchange: "articles" (topic)

**События**:
1. `article.created` - Новая статья создана
2. `article.viewed` - Статья просмотрена
3. `article.liked` - Лайк поставлен
4. `article.unliked` - Лайк снят

## Будущие расширения

### AI Service (планируется)
- Помощь в написании статей
- Генерация заголовков
- Грамматические проверки
- Улучшение текста

### Notification Service (планируется)
- Email уведомления
- Push уведомления
- Подписка на авторов

## Запуск

### Генерация proto файлов:
```bash
make proto
```

### Запуск всех сервисов:
```bash
docker-compose up -d
```

### Запуск отдельных сервисов для разработки:
```bash
# Auth Service
cd services/auth-service && go run cmd/main.go

# Article Service  
cd services/article-service && go run cmd/main.go

# Stats Service
cd services/stats-service && go run cmd/main.go

# API Gateway
cd services/api-gateway && go run cmd/main.go
```

## Переменные окружения

См. `docker-compose.yml` для полного списка переменных окружения для каждого сервиса.

## Преимущества архитектуры

1. **Масштабируемость**: Каждый сервис можно масштабировать независимо
2. **Изоляция**: Проблемы в одном сервисе не влияют на другие
3. **Гибкость**: Легко добавлять новые сервисы (AI, уведомления)
4. **Производительность**: gRPC для быстрой коммуникации между сервисами
5. **Асинхронность**: RabbitMQ для неблокирующих операций
6. **Безопасность**: Централизованная аутентификация через Auth Service
7. **Приватность**: Гибкая система видимости статей (public/private/link)

## Миграция с монолита

Текущий монолитный backend сохранен как `backend-legacy` на порту 8081. После тестирования микросервисов его можно удалить.
