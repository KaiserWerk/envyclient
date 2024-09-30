# envyclient

A simple client to comfortably set and get variables from a remote [envy](https://github.com/KaiserWerk/envy) instance.

### API

Create a new client with a base address, an environment name and the auth key (bearer token).

```golang
client := envyclient.NewClient("https://somehost:7000", "app-name", "default")
```

Set a variable in the environment. Only strings allowed.

```golang
name := "varname"
value := "some data"
err := client.SetVar(name, value)
```

Get a variable from the environment. It is automatically set as an environment variable and also returned.

```golang
name := "varname"
_, err := client.GetVar(name) // returns Var, error
```


Get all variables from the environment. They are automatically set as environment variables and also returned.

```golang
_, err := client.GetAllVars()  // returns []Var, error
```

### Example

```golang
client := envyclient.NewClient("http://domain.com:3030", "my_app", "default")
_ = client.SetVar("API_KEY", "w49t861w6t")
_ = client.SetVar("baseUrl", "http://cloud.host.org:8080")

_, _ = client.GetAllVars()

fmt.Println(os.Getenv("API_KEY")) // outputs "w49t861w6t"
fmt.Println(os.Getenv("baseUrl")) // outputs "http://cloud.host.org:8080"
```