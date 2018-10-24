package main

import (
	"fmt"
	"github.com/expanse-org/motan-go/main/demo"
	motan "github.com/expanse-org/motan-go"
)

func main() {
	runClientDemo()
}

func runClientDemo() {
	mccontext := motan.GetClientContext("./main/clientZkDemo.yaml")
	mccontext.Start(nil)
	//mclient := mccontext.GetClient("mytest-motan2")

	args := make(map[string]string, 16)
	args["name"] = "ray"
	args["id"] = "xxxx"
	//var reply string
	//err := mclient.Call("hello", []interface{}{args}, &reply)
	//if err != nil {
	//	fmt.Printf("motan call fail! err:%v\n", err)
	//} else {
	//	fmt.Printf("motan call success! reply:%s\n", reply)
	//}
	//
	//// async call
	//args["key"] = "test async"
	//result := mclient.Go("hello", []interface{}{args}, &reply, make(chan *motancore.AsyncResult, 1))
	//res := <-result.Done
	//if res.Error != nil {
	//	fmt.Printf("motan async call fail! err:%v\n", res.Error)
	//} else {
	//	fmt.Printf("motan async call success! reply:%+v\n", reply)
	//}

	mclient2 := mccontext.GetClient("mytest-demo")

	person := &demo.Person{}
	person.Id = 100
	person.Email = "email@sss.com"
	person.Name = "NameDemo"
	phoneNumber := new (demo.Person_PhoneNumber)
	phoneNumber.Number = "12232323"
	person.Phones = append(person.Phones, phoneNumber)
	address := &demo.AddressBook{}
	address.People = append(address.People, person)

	addressRes := &demo.AddressBook{}

	err := mclient2.Call("hello", []interface{}{address}, addressRes)
	if err != nil {
		fmt.Printf("motan call fail! err:%v\n", err)
	} else {
		fmt.Printf("motan call success!\n")
		fmt.Printf("addressRes.People.length :%d\n", len(addressRes.People))
		fmt.Printf("addressRes.People[0].Email :%s\n", addressRes.People[0].Email)
		fmt.Printf("addressRes.People[0].name :%s\n", addressRes.People[0].Name)
		fmt.Printf("addressRes.People[0].Phones[0].Number :%s\n", addressRes.People[0].Phones[0].Number)
	}

}
