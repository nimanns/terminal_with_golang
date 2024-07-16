package main 

import "fmt"


type Product struct {
	id int
	name string
	price int
	description string
}

func (p *Product) show(){
	fmt.Print(p.name)
	fmt.Print(" *** ")
	fmt.Print("Price: ")
	fmt.Println(p.price)
	fmt.Println(p.description)
}

func main(){
	var product Product = Product{0, "new product", 30, "new product"}
	product.show()
}
