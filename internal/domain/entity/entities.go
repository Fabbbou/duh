package entity

type GroupName string

const (
	Aliases GroupName = "aliases"
	Exports GroupName = "exports"
)

type Key = string
type Value = string

type StoreEntries map[Key]Value

type Store map[GroupName]StoreEntries
