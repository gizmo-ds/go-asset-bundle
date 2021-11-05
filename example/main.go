package main

import (
	"fmt"
	"mime"
	"net/http"

	bundle "github.com/GizmoOAO/go-asset-bundle"
)

func init() {
	// https://github.com/golang/go/issues/32350
	_ = mime.AddExtensionType(".js", "application/javascript; charset=utf-8")
}

func main() {
	ab, err := bundle.NewAssetBundle("../public.ab")
	if err != nil {
		panic(err)
	}
	err = ab.Bundle("../public", 1000)
	if err != nil {
		panic(err)
	}
	ab.Close()

	ab2, err := bundle.OpenAssetBundle("../public.ab")
	if err != nil {
		panic(err)
	}
	defer ab2.Close()

	fmt.Println("Version", ab2.Version)

	http.Handle("/", http.StripPrefix("/", http.FileServer(ab2)))
	addr := fmt.Sprintf("127.0.0.1:%d", 3000)
	fmt.Println("http server started on", addr)
	http.ListenAndServe(addr, nil)
}
