package main

import "fmt"

/*
设计模式之抽象工厂模式
从简单模式，再抽象出一个角色来，这个角色就叫发书人Assigner,任何人实现GetBook方法
都是发书人
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

// 发书人
type Assigner interface {
	GetBook(name string) Book
}

type assigner struct {
}

func (a *assigner) GetBook(name string) Book {
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
	var a assigner
	fmt.Println(a.GetBook("math").Name())
	fmt.Println(a.GetBook("english").Name())
	fmt.Println(a.GetBook("chinese").Name())
}
