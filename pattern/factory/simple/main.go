package main

import "fmt"

/*
设计模式之简单工厂模式
每次开学发书，分别为语文书，数学书，英语书，老师忙不过来，指定某个同学去发书，
这个同学就是工厂，别人都去那领书
*/

// 有多种书，可以获取书名
type Book interface {
	Name() string
}

// 定义多种书，都实现book的获取数名方法
type chineseBook struct {
	name string
}

func (cb *chineseBook) Name() string {
	return cb.name
}

type mathBook struct {
	name string
}

func (mb *mathBook) Name() string {
	return mb.name
}

type englishBook struct {
	name string
}

func (eb *englishBook) Name() string {
	return eb.name
}

// 现在有个同学发书（可以发多种书）
func GetBook(name string) Book {
	switch name {
	case "chinese":
		return &chineseBook{name: "语文书"}
	case "math":
		return &mathBook{name: "数学书"}
	case "english":
		return &englishBook{name: "英语书"}
	default:
		return nil
	}
}

func main() {
	book := GetBook("chinese")
	fmt.Println(book.Name())
}
