# AI Service - План разработки

## Обзор

AI Service будет предоставлять интеллектуальные функции для помощи авторам в написании статей.

## Архитектура

```
┌─────────────┐
│ API Gateway │
└──────┬──────┘
       │
       ↓
┌────────────────────┐
│   AI Service       │  gRPC :50054
│  - Text generation │
│  - Suggestions     │
│  - Grammar check   │
│  - Title ideas     │
└────────┬───────────┘
         │
         ↓
┌────────────────────┐
│  AI Provider       │
│  - OpenAI API      │
│  - Anthropic       │
│  - Local LLM       │
└────────────────────┘
```

## Proto Schema

```protobuf
syntax = "proto3";

package ai;

option go_package = "github.com/XRS0/blog/proto/ai";

service AIService {
  rpc GenerateContent(GenerateContentRequest) returns (GenerateContentResponse);
  rpc SuggestTitle(SuggestTitleRequest) returns (SuggestTitleResponse);
  rpc ImproveText(ImproveTextRequest) returns (ImproveTextResponse);
  rpc CheckGrammar(CheckGrammarRequest) returns (CheckGrammarResponse);
  rpc ContinueWriting(ContinueWritingRequest) returns (ContinueWritingResponse);
  rpc Summarize(SummarizeRequest) returns (SummarizeResponse);
}

message GenerateContentRequest {
  string topic = 1;
  string style = 2; // "technical", "casual", "formal"
  int32 max_length = 3;
}

message GenerateContentResponse {
  string content = 1;
  string error = 2;
}

message SuggestTitleRequest {
  string content = 1;
  int32 count = 2; // Количество вариантов
}

message SuggestTitleResponse {
  repeated string titles = 1;
  string error = 2;
}

message ImproveTextRequest {
  string text = 1;
  string focus = 2; // "clarity", "engagement", "brevity"
}

message ImproveTextResponse {
  string improved_text = 1;
  repeated Suggestion suggestions = 2;
  string error = 3;
}

message Suggestion {
  string original = 1;
  string suggested = 2;
  string reason = 3;
}

message CheckGrammarRequest {
  string text = 1;
}

message CheckGrammarResponse {
  repeated GrammarIssue issues = 1;
  string error = 2;
}

message GrammarIssue {
  string text = 1;
  int32 offset = 2;
  int32 length = 3;
  string type = 4; // "spelling", "grammar", "style"
  string message = 5;
  repeated string suggestions = 6;
}

message ContinueWritingRequest {
  string current_text = 1;
  int32 max_length = 2;
}

message ContinueWritingResponse {
  string continuation = 1;
  string error = 2;
}

message SummarizeRequest {
  string text = 1;
  int32 max_length = 2;
}

message SummarizeResponse {
  string summary = 1;
  string error = 2;
}
```

## Основные функции

### 1. Генерация контента
- Генерация статьи по теме
- Разные стили: технический, casual, формальный
- Контроль длины

### 2. Предложение заголовков
- Генерация 3-5 вариантов заголовков
- На основе содержимого статьи
- Оптимизация для SEO

### 3. Улучшение текста
- Улучшение ясности
- Повышение вовлеченности
- Сокращение многословия
- Конкретные предложения по изменениям

### 4. Проверка грамматики
- Орфография
- Грамматические ошибки
- Стилистические предложения
- Интеграция с LanguageTool или аналогами

### 5. Продолжение написания
- AI дописывает текст
- Сохраняет стиль и тон
- Контекстное продолжение

### 6. Суммаризация
- Краткое изложение длинных текстов
- Для превью и описаний

## API Gateway endpoints

```
POST /api/ai/generate
POST /api/ai/suggest-titles
POST /api/ai/improve
POST /api/ai/check-grammar
POST /api/ai/continue
POST /api/ai/summarize
```

## Интеграция с провайдерами

### OpenAI
```go
client := openai.NewClient(apiKey)
resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
    Model: "gpt-4",
    Messages: []openai.ChatCompletionMessage{
        {Role: "user", Content: prompt},
    },
})
```

### Anthropic Claude
```go
client := anthropic.NewClient(apiKey)
resp, err := client.Messages.Create(ctx, anthropic.MessageRequest{
    Model: "claude-3-opus-20240229",
    Messages: []anthropic.Message{
        {Role: "user", Content: prompt},
    },
})
```

### Local LLM (Ollama)
```go
client := ollama.NewClient("http://localhost:11434")
resp, err := client.Generate(ctx, ollama.GenerateRequest{
    Model: "llama2",
    Prompt: prompt,
})
```

## Кеширование

Используйте Redis для кеширования результатов:

```go
// Проверка кеша
cacheKey := fmt.Sprintf("ai:suggest-titles:%s", hashContent(content))
if cached, err := redis.Get(ctx, cacheKey).Result(); err == nil {
    return cached
}

// Запрос к AI
result := callAIProvider(content)

// Сохранение в кеш
redis.Set(ctx, cacheKey, result, 24*time.Hour)
```

## Rate Limiting

Ограничение запросов для предотвращения злоупотреблений:

```go
// Per user limit
limiter := rate.NewLimiter(rate.Every(time.Minute), 10) // 10 req/min
if !limiter.Allow() {
    return ErrRateLimitExceeded
}
```

## Стоимость и квоты

- Отслеживание использования токенов
- Квоты для пользователей
- Оплата за дополнительные запросы

```sql
CREATE TABLE ai_usage (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    operation VARCHAR(50),
    tokens_used INTEGER,
    cost_cents INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);
```

## Frontend интеграция

### Компонент AI Assistant

```typescript
interface AIAssistantProps {
  content: string;
  onSuggestion: (text: string) => void;
}

function AIAssistant({ content, onSuggestion }: AIAssistantProps) {
  const [suggestions, setSuggestions] = useState<string[]>([]);
  
  const generateTitles = async () => {
    const response = await api.post('/ai/suggest-titles', { content });
    setSuggestions(response.data.titles);
  };
  
  const improveText = async () => {
    const response = await api.post('/ai/improve', { 
      text: content,
      focus: 'clarity'
    });
    onSuggestion(response.data.improved_text);
  };
  
  return (
    <div className="ai-assistant">
      <button onClick={generateTitles}>Suggest Titles</button>
      <button onClick={improveText}>Improve Text</button>
      {suggestions.map(title => (
        <div key={title} onClick={() => onSuggestion(title)}>
          {title}
        </div>
      ))}
    </div>
  );
}
```

## Мониторинг

- Логирование всех AI запросов
- Метрики производительности
- Стоимость запросов
- Качество результатов (feedback от пользователей)

## Безопасность

1. **Фильтрация контента**: Проверка на вредоносный контент
2. **Валидация входных данных**: Ограничение длины, формата
3. **Аутентификация**: Только авторизованные пользователи
4. **Rate limiting**: Защита от злоупотреблений

## План разработки

### Фаза 1: MVP
- [ ] Создать proto схему
- [ ] Настроить OpenAI клиент
- [ ] Реализовать генерацию заголовков
- [ ] Реализовать улучшение текста
- [ ] Добавить endpoints в API Gateway

### Фаза 2: Расширение
- [ ] Продолжение написания
- [ ] Проверка грамматики
- [ ] Суммаризация
- [ ] Кеширование результатов

### Фаза 3: Оптимизация
- [ ] Поддержка нескольких провайдеров
- [ ] Балансировка нагрузки
- [ ] Мониторинг и аналитика
- [ ] Система квот и биллинга

### Фаза 4: Advanced
- [ ] Тонкая настройка под блоговый контент
- [ ] Персонализация стиля
- [ ] Контекстное обучение на статьях пользователя
- [ ] Multi-modal (изображения, код)

## Пример использования

```bash
# Предложить заголовки
curl -X POST http://localhost:8080/api/ai/suggest-titles \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "В этой статье я расскажу о микросервисной архитектуре...",
    "count": 5
  }'

# Улучшить текст
curl -X POST http://localhost:8080/api/ai/improve \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Микросервисы это хорошо потому что...",
    "focus": "clarity"
  }'

# Продолжить написание
curl -X POST http://localhost:8080/api/ai/continue \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "current_text": "Начало моей статьи про Go...",
    "max_length": 500
  }'
```

## Стоимость

### OpenAI GPT-4
- Input: $0.03 / 1K tokens
- Output: $0.06 / 1K tokens

### Claude 3
- Input: $0.015 / 1K tokens
- Output: $0.075 / 1K tokens

### Локальный LLM
- Бесплатно, но требует мощного железа
- Ollama с Llama 2 или Mistral

Рекомендация: начать с Claude 3 Haiku (дешевле) для MVP, потом добавить GPT-4 для премиум пользователей.
