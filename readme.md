# H
__http for Gophers, and human__

Go std lib `net/http` is easy to use, but verbose. Package `h` provides a simply yet flexible API for Gophers and other human. Readable, maintainable, debuggable, extendable and flexible API is the goal.

## Usage
### Basic
```go
client := NewClient().
    SetBaseURL("baseURL").
    SetHeader("Header1", "Value").
    SetHeader("Header2", "Value").
    Use(func(r *Request, res *http.Response, err error) (*http.Response, error) {
        if res.StatusCode == 500 {
        	log.Fatal(500)
        }
    }).
    SetTimeout(5 * time.Second)

res, err := client.Request(http.MethodGet, "www.google.com").
	SetHeader("Header3", "Value").
	SetBody(strings.NewReader("123")).
	Run()
```
This returns a standard `(http.Response, error)` pair, same return signature as `http.Client.Do` function.

As you can see, a convenient method chaining API is provided.

### Plugins
The power of `h` comes in with plugins. `h` itself has less than 120 lines of code. But it has a extendable design.

`h/hplug` package provides several powerful plugins, including retries. Please see [hplug examples](plugin/retry/retry_test.go).