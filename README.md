# gonfigs

Assign environment, argument and/or default values to struct variables using Golang tags.

## Description

[Gonfigs](https://github.com/Someone0nEarth/gonfigs) parses a *struct for "gonfigs" tags like `arg_name`, `env_name`, `default` and using them to set the value
of the corresponding fields.

If more than one gonfigs tag are matching for a field, the priority for setting the value is:

`arg_name` > `env_name` > `default`

Tags:
  - `arg_name` looks for a flag.Flag (command line argument) with the name corresponding to the tag value. Will be
    listed, when using '--help'  
  - `arg_env` looks for an environment variable with the name corresponding to the tag value
  - `default` will use the tags value for the field
  - `description` the description of the field. For documentation and/or if arg_name is although set, it will be
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
	User         *string `env_name:"DB_USER"          arg_name:"db_user"                         
	Password     string  `env_name:"DB_USER_PASSWORD" arg_name:"db_user_password"`
	Host         *string `env_name:"DB_HOST"          arg_name:"db_host"           default:"localhost"`
	Port         *uint   `env_name:"DB_PORT"          arg_name:"db_port"           default:"3306"`
	RootPassword string  `env_name:"DB_ROOT_PASSWORD" arg_name:"db_root_password"                       description:"Will be used for creating the database if necessary."`
	DatabaseName *string `env_name:"DB_NAME"          arg_name:"db_name"           default:"fancy_app"  description:"Database with this name will be created if not already existing."`
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


