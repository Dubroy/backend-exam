package main

import (
	"fmt"
	"reflect"
)

func swap[T any](a, b T) {
	// a 和 b 是指針類型（例如 *int）
	// 使用反射來獲取指針指向的值並交換

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	// 檢查是否為指針類型，如果不是則 panic（顯式調用）
	if va.Kind() != reflect.Ptr {
		panic("swap: first argument must be a pointer")
	}
	if vb.Kind() != reflect.Ptr {
		panic("swap: second argument must be a pointer")
	}

	// 獲取指針指向的值
	elemA := va.Elem()
	elemB := vb.Elem()

	// 創建臨時變數來保存值（使用 unsafe 來複製任意類型的值）
	temp := reflect.New(elemA.Type()).Elem()
	temp.Set(elemA)

	// 交換值
	elemA.Set(elemB)
	elemB.Set(temp)
}

func main() {
	a := 10
	b := 20

	fmt.Printf("a = %d, &a = %p\n", a, &a)
	fmt.Printf("b = %d, &b = %p\n", b, &b)

	swap(&a, &b)

	fmt.Printf("a = %d, &a = %p\n", a, &a)
	fmt.Printf("b = %d, &b = %p\n", b, &b)
}
