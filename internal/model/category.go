package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Category struct {
	Id   int    `gorm:"column:id;primaryKey" json:"id"`
	Name string `gorm:"column:name" json:"name"`
	gorm.Model
}

// Constructors
func NewCategory(id int, name string) *Category {
	return &Category{Id: id, Name: name}
}

// Methods
func (c *Category) DisplayCategory() {
	fmt.Printf("Category Id: %d\nCategory Name: %s\n", c.Id, c.Name)
}
