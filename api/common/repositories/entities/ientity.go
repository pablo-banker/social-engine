package entities

type IEntity interface {
	TableName() string
	LoadAssociations() []string
	GetID() any
}
