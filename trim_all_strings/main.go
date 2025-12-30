package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func TrimAllStrings(a any) {
	if a == nil {
		return
	}

	// 使用 map 記錄已訪問的指針地址，避免循環引用
	visited := make(map[uintptr]bool)
	trimAllStringsRecursive(reflect.ValueOf(a), visited)
}

func trimAllStringsRecursive(v reflect.Value, visited map[uintptr]bool) {
	if !v.IsValid() {
		return
	}

	// 處理指針類型
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}

		// 檢查是否已訪問過（避免循環引用）
		ptr := v.Pointer()
		if visited[ptr] {
			return
		}
		visited[ptr] = true

		// 遞迴處理指針指向的值
		trimAllStringsRecursive(v.Elem(), visited)
		return
	}

	// 處理字串類型
	if v.Kind() == reflect.String {
		if v.CanSet() {
			trimmed := strings.TrimSpace(v.String())
			v.SetString(trimmed)
		}
		return
	}

	// 處理結構體
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.CanInterface() {
				trimAllStringsRecursive(field, visited)
			}
		}
		return
	}

	// 處理切片
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			trimAllStringsRecursive(v.Index(i), visited)
		}
		return
	}

	// 處理陣列
	if v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			trimAllStringsRecursive(v.Index(i), visited)
		}
		return
	}

	// 處理 map
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			// 處理 key（如果是字串）
			trimAllStringsRecursive(key, visited)
			// 處理 value
			trimAllStringsRecursive(value, visited)
		}
		return
	}

	// 處理 interface{} 類型
	if v.Kind() == reflect.Interface {
		if !v.IsNil() {
			trimAllStringsRecursive(v.Elem(), visited)
		}
		return
	}
}

func main() {
	type Person struct {
		Name string
		Age  int
		Next *Person
	}

	a := &Person{
		Name: " name ",
		Age:  20,
		Next: &Person{
			Name: " name2 ",
			Age:  21,
			Next: &Person{
				Name: " name3 ",
				Age:  22,
			},
		},
	}

	TrimAllStrings(&a)

	m, _ := json.Marshal(a)

	fmt.Println(string(m))

	a.Next = a

	TrimAllStrings(&a)

	fmt.Println(a.Next.Next.Name == "name")
}
