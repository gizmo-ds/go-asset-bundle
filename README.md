# Go Asset Bundle

[![Go Report Card](https://goreportcard.com/badge/github.com/project-vrcat/go-asset-bundle?style=flat-square)](https://goreportcard.com/report/github.com/project-vrcat/go-asset-bundle)
[![License](https://img.shields.io/github/license/project-vrcat/go-asset-bundle?style=flat-square)](./LICENSE)

Just like [asar](https://github.com/electron/asar) 😂

## Example

[example](example/main.go)

### Create an AssetBundle

```go
ab, _ := bundle.NewAssetBundle("./public.ab")
defer ab.Close()

var version uint16 = 1000
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

## Thanks

Thanks to [JetBrains](https://jb.gg/OpenSource) for the open source license(s).

[![JetBrains Logo](https://raw.githubusercontent.com/project-vrcat/VRChatConfigurationEditor/main/images/jetbrains.svg)](https://jb.gg/OpenSource)

## License

Code is distributed under [MIT](./LICENSE) license, feel free to use it in your proprietary projects as well.
