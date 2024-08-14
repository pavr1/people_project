package models

import (
	"github.com/pavr1/people/config"
)

type Person struct {
	config   *config.Config
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Age      int32  `json:"age"`
}

func NewPerson(config *config.Config) Person {
	return Person{
		config: config,
	}
}

func (p *Person) Populate(name string, lastName string, age int32) {
	p.Name = name
	p.LastName = lastName
	p.Age = age
}
