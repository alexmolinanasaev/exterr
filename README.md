## exterr
Пакет для отображения user-friendly ошибок как для разработчика, так и для пользователя

Мотивация. При использовании стандартных интерфейсов Error в Go порой бывает сложно тслеить тип ошибки и его место появления.
трейс недостаточно информативен. Проблема что на фронт уходи ошибка о состоянии базы данных
Ббилотека помогает упростить и унифицировать работу с ошибками. 

Перечислить особенности (кастомный стектрейс)

Сооответствие интерфейску Error
Эту ошибку можно ичспользовать как обычный Error

Хранение кода ошибка для распоззнования категории
Основной и альтернативный для внешних сервисов



# Установка
```bash
go get github.com/alexmolinanasaev/exterr
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
func New(msg string) ErrExtender

// Ошибка с форматированной строкой описания
func Newf(format string, a ...interface{}) ErrExtender

// Добавление к стандартному Error (err) описания msg "{msg}:{str}"
func NewWithErr(msg string, err error) ErrExtender

// 
func NewWithAlt(msg, altMsg string) ErrExtender

// Ошибка 
func NewWithType(msg, altMsg string, t ErrType) ErrExtender
```








// temp
Newf(string, {})
NewWithErr(string, string)
NewWithAlt(string, string)