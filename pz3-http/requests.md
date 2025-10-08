# Тестовые curl-запросы для API (Windows PowerShell)

## 1. Проверка работоспособности сервера
```sh
curl http://localhost:8080/health
```

## 2. Создать задачу
```sh
curl -Method POST http://localhost:8080/tasks `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"title":"Listen to Post Malone"}'
```

## 3. Создать задачу (ошибка: слишком короткий title)
```sh
curl -Method POST http://localhost:8080/tasks `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"title":"no"}'
```

## 4. Получить список задач
```sh
curl http://localhost:8080/tasks
```

## 5. Фильтрация задач
```sh
curl "http://localhost:8080/tasks?q=Post Malone"
```

## 6. Получить задачу по id
```sh
curl http://localhost:8080/tasks/1
```

## 7. Получить несуществующую задачу
```sh
curl http://localhost:8080/tasks/999
```

## 8. Отметить задачу как выполненную
```sh
curl -Method PATCH http://localhost:8080/tasks/1
```

## 9. Удалить задачу
```sh
curl -Method DELETE http://localhost:8080/tasks/1
```

## 10. Удалить несуществующую задачу
```sh
curl -Method DELETE http://localhost:8080/tasks/999
```