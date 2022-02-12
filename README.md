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