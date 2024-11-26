package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Println("timeout ao requisitar a cotacao")
        }
        panic(err)
    }
	defer res.Body.Close()
    body, err := io.ReadAll(res.Body)
    if err != nil {
        panic(err)
    }


	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	_, err = f.Write([]byte(fmt.Sprintf("Dólar:%s", body)))
	if err != nil {
		panic(err)
	}
	f.Close()

	fmt.Print(string(body))
}
