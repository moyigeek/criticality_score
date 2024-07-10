package main

import (
    "fmt"
    "os"
    "github.com/HUSTSeclab/criticality_score/pkg/collector/archlinux"
    "github.com/HUSTSeclab/criticality_score/pkg/collector/debian"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: main <archlinux|debian>")
        return
    }

    switch os.Args[1] {
    case "archlinux":
        archlinux.Archlinux()
    case "debian":
        debian.Debian() // 确保 debian.go 中有一个名为 Debian 的函数
    default:
        fmt.Println("Unknown command:", os.Args[1])
        fmt.Println("Usage: main <archlinux|debian>")
    }
}
