# gonfigs

Assign environment, argument and/or default values to struct variables using Golang tags.

## Description

[Gonfigs](https://github.com/Someone0nEarth/gonfigs) parses a *struct for "gonfigs" tags like `argName`, `envName`, `default` and using them to set the value
of the corresponding fields.

If more than one gonfigs tag are matching for a field, the priority for setting the value is:

`argName` > `envName` > `default`

Tags:
  - `argName` looks for a flag.Flag (command line argument) with the name corresponding to the tag value. Will be
    listed, when using '--help'  
  - `envName` looks for an environment variable with the name corresponding to the tag value
  - `default` will use the tags value for the field
  - `description` the description of the field. For documentation and/or if argName is although set, it will be
    shown in the `usage` field when using `--help` in the command line.

Struct fields will be set only, if their current values are zero (Value.IsZero()==true).

At the moment, the following struct fields types are supported: unit, *uint, string and *string

Gonfigs will panic, if the `config` is not a struct, `config` is not a pointer and the tags are used on unsupported
field types.

### Example

```go
import (
  "github.com/Someone0nEarth/gonfigs"
)

...

type DatabaseConfig struct {
	User         *string `envName:"DB_USER"          argName:"db_user"                         
	Password     string  `envName:"DB_USER_PASSWORD" argName:"db_user_password"`
	Host         *string `envName:"DB_HOST"          argName:"db_host"           default:"localhost"`
	Port         *uint   `envName:"DB_PORT"          argName:"db_port"           default:"3306"`
	RootPassword string  `envName:"DB_ROOT_PASSWORD" argName:"db_root_password"                       description:"Will be used for creating the database if necessary."`
	DatabaseName *string `envName:"DB_NAME"          argName:"db_name"           default:"fancy_app"  description:"Database with this name will be created if not already existing."`
}

...

dbConfig := DatabaseConfig{}

gonfigs.Parse(&dbConfig)

...
```

## Similar Projects

If gonfigs is not fitting your needs, maybe one of the following projects will do it:

- <https://github.com/caarlos0/env>
- <https://github.com/tkanos/gonfig> (I "created" the name `gonfigs` before I knew about `gonfig` ;-) 
- <https://github.com/spf13/viper>

## Feedback

Feel free to suggests, fork, contribute, request, improve, criticize and so on :)

Have fun,

someone.earth


