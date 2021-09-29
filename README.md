## exterr
Пакет для отображения user-friendly ошибок как для разработчика, так и для пользователя

# Для чего
При использовании стандартного интерфейса Error в Golang сложно отследить тип ошибки и место его появления, а трассировка недостаточно информативна (или наоборот, даёт избыточную информацию). Также есть необходимость исключить выдачу во фронтенд код и ошибки состояния (напрмер, при неудачном подключении в базе данных). Библиотека позволяет упростить и унифицировать обработку ошибок.

# Особенности
- Соответствие стандартному интерфейсу Error;
- Возможность добавить код ошибки для удобной для обработки;
- Возможность добавить альтернативное описание для внешних сервисов;
- Кастомизированный stacktrace:
  -  Сырой (raw);
  -  Тегированный (tagged);
  -  В формате JSON (fancy).


# Установка
```bash
$ go get github.com/alexmolinanasaev/exterr
```

# Импорт в проект
```go
import (
	"github.com/alexmolinanasaev/exterr"
)
```

# Основные функции:
```go
// Ошибка с описанием
// Пример: exterr.New("SQL error!")
func New(msg string) ErrExtender

// Ошибка с форматированной строкой описания
// Пример: exterr.Newf("ERROR: %s", desc)
func Newf(format string, a ...interface{}) ErrExtender

// Добавление к стандартному Error (err) описания msg "{msg}:{err}"
// Пример: exterr.NewWithErr("Auth error", err)
func NewWithErr(msg string, err error) ErrExtender

// Добавление альтернативного описания
// Пример: exterr.NewWithAlt("sql no rows", "<SQL_NO_ROWS>")
func NewWithAlt(msg, altMsg string) ErrExtender

// Добавление номера (ErrType) к ошибке
// Пример: exterr.NewWithType("sql no rows", "user not found", 1005)
func NewWithType(msg, altMsg string, t ErrType) ErrExtender
```

# Работа со stacktrace
Начальный stacktrace формирует 1 строку в месте возникновения ошибки и может быть дополнен с помощью функции **AddTrace()**, например:
```go
func DatabaseInit() error {
    DB, err := sql.Open("mysql", "user:password@/test_db")
    if err != nil {
		return exterr.NewWithType("sql connection", "<SQL_CONNECTION_ERROR>", 1001)
	}
}

func main() {
    err := DatabaseInit()
    if err != nil {
        err.AddTrace()
        log.Fatal(err.Error())
    }
}
```

# Тестирование
Тесты расположены в каталоге **/tests**.
```bash
$ cd tests
$ go test
```
