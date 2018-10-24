/*

  Copyright 2017 Loopring Project Ltd (Loopring Foundation).

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package motan_test

import (
	"fmt"
	"github.com/expanse-org/relay-lib/motan"
	"github.com/expanse-org/relay-lib/motan/demo"
	"testing"
)

func TestInitClient(t *testing.T) {

	options := motan.MotanClientOptions{}
	options.ConfFile = "./clientZkDemo.yaml"
	options.ClientId = "mytest-demo"
	mclient := motan.InitClient(options)
	//mccontext := weibomotan.GetClientContext("./clientZkDemo.yaml")
	//mccontext.Start(nil)
	//mclient := mccontext.GetClient("mytest-motan2")

	//args := make(map[string]string, 16)
	//args["name"] = "ray"
	//args["id"] = "xxxx"
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

	person := &demo.Person{}
	person.Id = 100
	person.Email = "email@sss.com"
	person.Name = "NameDemo"
	phoneNumber := new(demo.Person_PhoneNumber)
	phoneNumber.Number = "12232323"
	person.Phones = append(person.Phones, phoneNumber)
	address := &demo.AddressBook{}
	address.People = append(address.People, person)
	//address := make(map[string]string)
	//address["aaa"] = "address111111"
	addressRes := &demo.AddressBook{}

	err := mclient.Call("hello", []interface{}{address}, addressRes)
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
