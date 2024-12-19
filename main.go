package main

import (
    "fmt"
    "flag"
    "os"
)

var (
    language *string
    list_languages *bool
    version *bool
)

//Initilize utility flag
func init() {
    version = flag.Bool("version", false, "Check program version")
    language = flag.String("lang", "cpp", "Specify a language")
    list_languages = flag.Bool("listlang", false, "List all available languages")
}


func main() {
    //Call the flag and it will handle args
    flag.Parse()

    if *version {
        fmt.Printf("monkeycode version 0.1.0\n");
        os.Exit(0)
    }
    
    fmt.Println(*language, *version);
}
