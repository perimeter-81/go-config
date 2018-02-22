## Example

A very basic example (check the `examples` folder):

```go
package main

import (
	"fmt"
	"time"
	
	"github.com/neliseev/go-config"
)

type params struct {
	Home         string        `env:"HOME"`
	Port         int           `env:"PORT" envDefault:"3000"`
	IsProduction bool          `env:"PRODUCTION"`
	Hosts        []string      `env:"HOSTS" envSeparator:":"`
	Duration     time.Duration `env:"DURATION"`
}

func main() {
	p := new(params)
	if err := config.Parse(nil, p); err != nil {
		fmt.Printf("%+v\n", err)
	}
	
	fmt.Printf("%+v\n", p)
}
```

You can run it like this:

```sh
$ PRODUCTION=true HOSTS="host1:host2:host3" DURATION=1s go run examples/first.go
{Home:/your/home Port:3000 IsProduction:true Hosts:[host1 host2 host3] Duration:1s}
```

## Supported types and defaults for ENV

The library has support for the following types:

* `string`
* `int`
* `int64`
* `bool`
* `float32`
* `float64`
* `[]string`
* `[]int`
* `[]bool`
* `[]float32`
* `[]float64`

If you set the `envDefault` tag for something, this value will be used in the
case of absence of it in the environment. If you don't do that AND the
environment variable is also not set, the zero-value
of the type will be used: empty for `string`s, `false` for `bool`s
and `0` for `int`s.

By default, slice types will split the environment value on `,`; you can change this behavior by setting the `envSeparator` tag.

Also u can provide []byte with json/yaml file, ordering for build config that: first parse config file and after ENV, so config will merged.

## Required fields

The `env` tag option `required` (e.g., `env:"tagKey,required"`) can be added
to ensure that some environment variable is set.  In the example above,
an error is returned if the `config` struct is changed to:


```go
type config struct {
    Home         string   `yaml:"home" env:"HOME"`
    Port         int      `json:"port" env:"PORT" envDefault:"3000"`
    IsProduction bool     `env:"PRODUCTION"`
    Hosts        []string `env:"HOSTS" envSeparator:":"`
    SecretKey    string   `env:"SECRET_KEY,required"`
}
```
