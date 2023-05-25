package models

type Profession struct {
	Id         int
	Name       string
	OtherNames []string
	Level      int
	ParentId   int
	ProfRoleId int
}

type Position struct {
	Id         int
	Name       string
	OtherNames []string
}
