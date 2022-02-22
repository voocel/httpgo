<p align="center">
    <h1 align="center">httpgo</h1>
    <p align="center">A fast and simple go http client </p>
</p>

## ⚙️ Installation

```
go get -u github.com/voocel/httpgo
```

## 👀 Example
#### 📖GET

```go
res, err := httpgo.Get("http://www.google.com").Do()
```

#### 📖POST
```go
res, err := httpgo.Post("http://www.google.com").Do()
```

#### 📖SetTimeout
```go
res, err := httpgo.Get("http://www.google.com").SetTimeout(5 * time.Second).Do()
```

## 🔥 Supported Methods
* [x] SetQuery(`key, value string`)
* [x] SetQueries(`m map[string]string`)
* [x] SetForm(`m map[string]string`)
* [x] SetJSON(`v string`)
* [x] SetTimeout(`t time.Duration`)
* [x] SetHeader(`key, value string`)
* [x] SetHeaders(`m map[string]string`)
* [x] AddHeader(`key, value string`)
* [x] AddHeaders(`m map[string]string`)
* [x] SetUA(`ua string`)
* [x] AddCookie(`c *http.Cookie`)
* [x] BasicAuth(`username, password string`)