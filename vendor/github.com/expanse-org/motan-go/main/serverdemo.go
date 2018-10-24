package main

import (
	"bytes"
	"fmt"
	"time"
	"github.com/expanse-org/motan-go/main/demo"
	motan "github.com/expanse-org/motan-go"
)

func main() {
	runServerDemo()
}

func runServerDemo() {
	mscontext := motan.GetMotanServerContext("./main/serverZkDemo.yaml")
	//mscontext.RegisterService(&Motan2TestService{}, "")
	mscontext.RegisterService(&MotanDemoService{}, "")
	mscontext.Start(nil)
	mscontext.ServicesAvailable() //注册服务后，默认并不提供服务，调用此方法后才会正式提供服务。需要根据实际使用场景决定提供服务的时机。作用与java版本中的服务端心跳开关一致。
	time.Sleep(time.Second * 50000000)
}

type MotanDemoService struct{}

func (m *MotanDemoService) Hello(book *demo.AddressBook) *demo.AddressBook {
	fmt.Printf("book.People.length :%d\n", len(book.People))
	fmt.Printf("book.People[0].Email :%s\n", book.People[0].Email)
	fmt.Printf("book.People[0].name :%s\n", book.People[0].Name)
	book.People[0].Name = book.People[0].Name + " >>>"
	book.People[0].Email = book.People[0].Email + " >>>"
	book.People[0].Phones[0].Number = book.People[0].Phones[0].Number + " >>>"
	return book
}

type Motan2TestService struct{}

func (m *Motan2TestService) Hello(params map[string]string) string {
	if params == nil {
		return "params is nil!"
	}
	var buffer bytes.Buffer
	for k, v := range params {
		if buffer.Len() > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(k)
		buffer.WriteString("=")
		buffer.WriteString(v)

	}
	fmt.Printf("Motan2TestService hello:%s\n", buffer.String())
	return buffer.String()
}
