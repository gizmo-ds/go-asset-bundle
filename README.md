# Go Asset Bundle

Just like [asar](https://github.com/electron/asar) 😂

## Example

[examples](example/main.go)

### Create an AssetBundle

```go
ab, _ := bundle.NewAssetBundle("./public.ab")
defer ab.Close()

version := 1000
ab.Bundle("./public", version)
```

### Use an AssetBundle

```go
ab, _ := bundle.OpenAssetBundle("./public.ab")
defer ab.Close()

fmt.Println("Version", ab.Version)

http.Handle("/", http.StripPrefix("/", http.FileServer(ab)))
addr := fmt.Sprintf("127.0.0.1:%d", 3000)
fmt.Println("http server started on", addr)
http.ListenAndServe(addr, nil)
```

## License

[MIT](./LICENSE)
