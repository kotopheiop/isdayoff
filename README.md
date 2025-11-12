# Isdayoff

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org)
[![Test Coverage](https://img.shields.io/badge/coverage-81.1%25-brightgreen.svg)](./isdayoff_test.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/kotopheiop/isdayoff)](https://goreportcard.com/report/github.com/kotopheiop/isdayoff)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/kotopheiop/isdayoff.svg)](https://pkg.go.dev/github.com/kotopheiop/isdayoff)

Клиент для работы с [Isdayoff API](https://isdayoff.ru/)

## Требования

Go 1.25+

## Установка

Убедитесь, что ваш проект использует Go Modules (если он уже использует, в корне будет файл `go.mod`):

``` sh
go mod init
```

Затем подключите модуль isdayoff в вашей Go программе с помощью `import`:

``` go
import (
    "github.com/kotopheiop/isdayoff"
)
```

Запустите любую из обычных команд `go` (`build`/`install`/`test`). Инструментарий Go автоматически разрешит и загрузит модуль.

Альтернативно, вы также можете явно выполнить `go get` для пакета:

```
go get -u github.com/kotopheiop/isdayoff
```

## Пример

```go
package main

import (
	"fmt"
	"github.com/kotopheiop/isdayoff"
)

func main() {
	dayOff := isdayoff.New()
	countryCode := isdayoff.CountryCodeKazakhstan
	pre := false
	covid := false
	day, err := dayOff.Tomorrow(isdayoff.Params{
		CountryCode: &countryCode,
		Pre:         &pre,
		Covid:       &covid,
	})    

	fmt.Println(day) // 0
}
```

## Примечание: 
- Названия часовых поясов (TZ) должны быть взяты из [IANA](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List)

## Лицензия

Этот проект распространяется под лицензией MIT. См. файл [LICENSE](LICENSE) для подробностей.

Вы можете свободно использовать, изменять и распространять этот код без каких-либо ограничений.
