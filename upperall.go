package main

import (
    "path/filepath"
    "os"
    "strings"
    "flag"
    "fmt"
    "strconv"
)

var target string = ""
var upper string = ""
var m = map[int][]string{}

func visit(path string, f os.FileInfo, err error) error {
    if !f.IsDir() {
        //fmt.Printf("%s\n", path)
        fi, err := os.OpenFile(path, os.O_RDWR, 0666)
        if err != nil {
            return err    
        }
        fsize := f.Size()
        buffsize := int64(1024)
        buf := make([]byte, buffsize)
        matches := 0
        for it := int64(0); it < fsize; it+=buffsize {
            n := 0
            i := int64(0)
            if it == 0 {
                i = 0
            } else {
                i = it - int64(len(target)) + 1
            }
            n, err = fi.ReadAt(buf, i)
            if n > 0 {
                s := string(buf[:n])
                //fmt.Printf("data: %s\n", s)
                for {
                    offset := strings.Index(s, target) 
                    if offset >=0 {
                        matches++
                        s = s[offset+len(target):]
                        //fmt.Printf("n-data: %s\n", s)
                        //fmt.Printf("Found: %s - %d\n", path, matches)
                    } else {
                        break
                    }
                }
                if matches > 0 {
                    news := strings.Replace(string(buf[:n]), target, upper, -1)
                    _, err = fi.WriteAt([]byte(news), i) 
                    m[matches] = append(m[matches], path)
                }
            }
        }
        err = fi.Close()
    }
    return err
} 


func main() {
    flag.Parse()
    target = flag.Arg(0)
    upper = strings.ToUpper(target)
    _ = filepath.Walk(".", visit)
    //fmt.Printf("Walk ends. err=%v\n", err)
    output := ""
    for k,v := range m {
        for i := 0; i < len(v); i++ {
            output = strconv.Itoa(k) + "\t" + v[i] + "\n" + output;
        }    
    }        
    fmt.Printf("Changes\tFile Name\n%s", output)
}
