package main

import (
    "path/filepath"
    "os"
    "flag"
    "fmt"
)

func visit(path string, f os.FileInfo, err error) error {
    fmt.Printf("%s\n", path)
    return nil
} 


func main() {
    flag.Parse()
    root := flag.Arg(0)
    err := filepath.Walk(root, visit)
    fmt.Printf("Walk ends. err=%v\n", err)
}
