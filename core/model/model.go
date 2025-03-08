package model

type IModel interface {
	TableName() string
	Read() []string
	Write() []string
	Update() []string
	Delete() []string
}
