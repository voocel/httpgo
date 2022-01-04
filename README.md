# httpgo
a fast and simple go http client

## Installation

```
go get -u github.com/voocel/httpgo
```

## Example
#### GET

```go
res, err := httpgo.Get("http://www.google.com").Do()
```

#### POST
```go
res, err := httpgo.Post("http://www.google.com").Do()
```
