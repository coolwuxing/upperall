package main

import (
    "path/filepath"
    "os"
    "strings"
    "flag"
    "fmt"
    "strconv"
    "sort"
)

var target string = ""
var upper string = ""
var m = []SearchMatch{}

type SearchMatch struct {
    file    string
    occurrences int
}

// implements sort.Interface for []SearchMatch
type ByOccurrences []SearchMatch

func (a ByOccurrences) Len() int           { return len(a) }
func (a ByOccurrences) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOccurrences) Less(i, j int) bool { return a[i].occurrences> a[j].occurrences}

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
                    sm := SearchMatch{path, matches}
                    m = append(m, sm)
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
    if len(target) < 1 || len(target) > 1023 {
        fmt.Printf("Please input a string (len 1~1024)\n")
        return
    }
    upper = strings.ToUpper(target)
    _ = filepath.Walk(".", visit)
    sort.Sort(ByOccurrences(m))
    output := ""
    for i := range m {
        output += strconv.Itoa(m[i].occurrences) + "\t" + m[i].file + "\n";
    }        
    fmt.Printf("Changes\tFile Name\n%s", output)
}
