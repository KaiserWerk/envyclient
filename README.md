# envyclient

A simple client to comfortably set and get variables from a remote [envy](https://github.com/KaiserWerk/envy) instance.

### API

Create a new client with a base address (without trailing slash), a scope name and the auth key (bearer token).

```golang
// no trailing slash
client := envyclient.NewClient("https://somehost:7000", "myScope", "default")
```

Permanently set a variable in the scope. Only strings are allowed.
An existing variable is overwritten.

```golang
name := "varname"
value := "some data"
err := client.SetVar(name, value)
```

For the purpose of returning variables, they implemented as a `Var`:

```golang
type Var struct {
    Name string
    Value string
}
```

Get a variable from the scope. It is automatically set as an environment variable and also returned.

```golang
name := "varname"
_, err := client.GetVar(name) // returns Var, error
```


Get all variables from the scope. They are automatically set as environment variables and also returned.

```golang
_, err := client.GetAllVars()  // returns []Var, error
```

### Example

```golang
// no trailing slash
client := envyclient.NewClient("http://domain.com:3030", "my_app", "default")
_ = client.SetVar("API_KEY", "w49t861w6t")
_ = client.SetVar("baseUrl", "http://cloud.host.org:8080")

_, _ = client.GetAllVars()

fmt.Println(os.Getenv("API_KEY")) // outputs "w49t861w6t"
fmt.Println(os.Getenv("baseUrl")) // outputs "http://cloud.host.org:8080"
```