# Go-Config

An annotation-based, lightweight and customizable util package to read configuration from JSON, YAML or Environment.

## Usage

See the minimal example below to read configuration of any type from environment variables and print to StdOut:

```golang
import "github.com/sbreit1/go-config"

var conf Config
config.FromEnvironment("EXAMPLE", &conf)

config.Print("Conf", &conf)
```

### Supported Types

| Type Class | Types |
| ---------- | ----- |
| Base Types | `string`, `bool`, `int` |
| Complex Types | Struct, Array, Slice |
| Default Types | `time.Duration`, `time.Time` |

### Planned Types

| Type Class | Types |
| ---------- | ----- |
| Base Types | All remaining types |
| Complex Types | Map |
| `[]byte` | Read from base64 string |
| Custom Types | Interfaces `FromEnv` and `FromEnvValue` |

## Read from Environment

The most common way of configuring an application is by defining structs with the corresponding configuration parameters, but you can use any datatype for configuration. Use `FromEnvironment` to fill configuration values from environment variables. The variable name is assembled by the type hierarchy. Use the `config` Tag to customize naming.

```golang
// Config is the main struct for configuration
type Config struct {
    // MainDatabase is read from {PREFIX}_DB
    // without the explicit env name "DB", the field name in uppercase MAINDATABASE is used.
    MainDatabase Database `config:"env:DB"`
}

// Database is nested in Config.
type Database struct {
    // Address is read from {DB-PREFIX}_ADDRESS
    Address string
    // User is read from {DB-PREFIX}_USER
    User string
    // Pass is read from {DB-PREFIX}_PASS
    Pass string
}

var conf Config
// read configuration from environment with prefix MAIN:
//   MAIN_DB_ADDRESS
//   MAIN_DB_USER
//   MAIN_DB_PASS
// the flat environment variable names resemble the type hierarchy
config.FromEnvironment("MAIN", &conf)
```

### Slices and Arrays from Environment

Slices have variable length, which is also read from environment. See the following example to read a slice with two entries:

```golang
// Environment:
//   MAIN_LIST_NUM = "2"  (required for slice initialization)
//   MAIN_LIST_0 = "42"
//   MAIN_LIST_1 = "1337"

type Config struct {
    List []int
}

var conf Config
config.FromEnvironment("MAIN", &conf)
// conf.List = []int{ 42, 1337 }
```

The `NUM` entry is not required for Arrays as their length is fixed.

### Duration from Environment

Values of type `time.Duration` can be initialized by an [ISO 8601 Duration String](https://en.wikipedia.org/wiki/ISO_8601#Durations) or a similar short form:

```golang
// Environment:
//   MAIN = "1y 4d 13m 5s"  (required for slice initialization)
//   MAIN = "P1Y4DT13M5S"   (same duration in ISO 8601 format)

// you do not need to encapsulate all values in structs
var d time.Duration
config.FromEnvironment("MAIN", &d)
// d = 1 year + 4 days + 13 minutes + 5 seconds
```

### Time from Environment

**TODO**

## Default Values

You can define default values to use, when no configuration values are available. These values should have the same format as the corresponding environment values.

```golang
type Config struct {
    User string            `config:"default:Jon Doe"`
    Age int                `config:"default:42"`
    Male bool              `config:"default:yes"`
    Patience time.Duration `config:"default:1m 8s"`
}
```


## Pretty Print

The `go-config` allows you to print out the configuration to StdOut for logging and debugging purposes omitting sensitive values.

```golang
type Config struct {
    // User is printed out as {PREFIX}.User
    User string
    // Pass is not printed out at all
    Pass string `config:"print:[none]"`
    // PhoneNumber will print out "******" as {PREFIX}.Phone when it is not empty.
    PhoneNumber string `config:"print:Phone:[mask]"`
    // Children is printed out as {PREFIX}.Children containing only the number of elements.
    Children []string `config:"print:[len]"`
}

conf := Config{ "Jon Doe", "secret", "0123456789", []string{ "Jane", "Joe" }}
config.Print("Conf", &conf)
// Example output:
//   Conf.User:     Jon Doe
//   Conf.Phone:    ******
//   Conf.Children: 2
```