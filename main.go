package main

import (
    "flag"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "sort"
)

func main() {
    // 1. 读取命令行参数：-path 指定要分析的目录
    dirPath := flag.String("path", ".", "要分析的文件夹路径，默认为当前目录")
    flag.Parse()

    // 2. 确认目录是否存在
    info, err := os.Stat(*dirPath)
    if err != nil {
        fmt.Println("错误：无法访问路径:", err)
        return
    }
    if !info.IsDir() {
        fmt.Println("错误：指定路径不是文件夹")
        return
    }

    fmt.Println("正在分析文件夹:", *dirPath)

    // 3. 用 map 统计每种后缀的总大小
    extSizeMap := make(map[string]int64)

    err = filepath.WalkDir(*dirPath, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            // 某些文件可能权限错误，打印一下然后继续
            fmt.Println("访问文件出错:", err)
            return nil
        }

        if d.IsDir() {
            return nil
        }

        // 获取文件信息
        info, err := d.Info()
        if err != nil {
            fmt.Println("读取文件信息出错:", err)
            return nil
        }

        size := info.Size()

        // 文件后缀，比如 ".jpg"、".txt"，没有后缀就用 "<no_ext>"
        ext := filepath.Ext(path)
        if ext == "" {
            ext = "<no_ext>"
        }

        extSizeMap[ext] += size
        return nil
    })

    if err != nil {
        fmt.Println("遍历文件夹时出错:", err)
        return
    }

    // 4. 把结果按大小排序输出
    type extStat struct {
        Ext  string
        Size int64
    }

    var stats []extStat
    for ext, size := range extSizeMap {
        stats = append(stats, extStat{Ext: ext, Size: size})
    }

    sort.Slice(stats, func(i, j int) bool {
        return stats[i].Size > stats[j].Size
    })

    fmt.Println("\n按文件类型统计的大小：")
    fmt.Println("-----------------------------------")
    var totalSize int64
    for _, s := range stats {
        totalSize += s.Size
        fmt.Printf("%-10s  %10.2f MB\n", s.Ext, float64(s.Size)/1024.0/1024.0)
    }
    fmt.Println("-----------------------------------")
    fmt.Printf("总大小：%.2f MB\n", float64(totalSize)/1024.0/1024.0)
}
