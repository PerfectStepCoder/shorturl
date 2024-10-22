# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

# Разное в проекте
### Запуск линтера
> /Users/dmitrii/go/bin/golangci-lint run

### Зависимости
Установит зависимости найденные в коде проекта
> go mod tidy

### Запуск тестов
> go test -v cmd/shortener/main_test.go

### Запуск с флагами
> go run ./cmd/shortener/main.go -a 0.0.0.0:9999 -b http://localhost:9999

> go run ./cmd/shortener/main.go -a 0.0.0.0:9999 -b http://localhost:9999 -d postgres://admin:password@localhost:6434/urlservice

### Форматирование кода
> gofmt -s -w .
-s simplifies the code
-w writes results directly

> /Users/dmitrii/go/bin/goimports -local "github.com/PerfectStepCoder/shorturl/cmd/shortener" -w main.go
> find . -name '*.go' | xargs /Users/dmitrii/go/bin/goimports -w -local 

/Users/dmitrii/go/bin/goimports
### Проверка стуктурных тегов
> go vet -structtag

## Тесты производительности
### Бенчмарки
> cd ./internal/tests/
> go test --bench .
> go test -bench . -benchmem
> go test -bench=. -cpuprofile=./base.pprof запись профиля
> go tool pprof bench.test base.pprof  анализ в консоле

### Профилирование
> go run ./cmd/shortener/main.go
> http://localhost:8080/debug/pprof/
Сравнение профилей:
> go tool pprof -top -diff_base=./base.pprof ./result.pprof

## Документация
> cd корень проекта
> /Users/dmitrii/go/bin/godoc -http=:8088 
-goroot=.
