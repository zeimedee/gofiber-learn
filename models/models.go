package models

type Employee struct {
	Id     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}

type Admin struct {
	User     string `json:"user" bson:"user,omitempty"`
	Password string `json:"password" bson:"password,omitempty"`
}
