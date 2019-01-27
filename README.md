# GoZix Redigo

## Dependencies

* [viper](https://github.com/gozix/viper)

## Configuration example

```json
{
  "redis": {
    "host": "127.0.0.1",
    "port": "6379",
    "password": "somepassword",
    "max_idle": 3,
    "max_active": 100,
    "idle_timeout": "4m"
  }
}
```
"password" field is optional and ignored if empty