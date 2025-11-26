package entity

type GroupName string

const (
	Aliases GroupName = "aliases"
	Exports GroupName = "exports"
)

type Key = string
type Value = string

type StorageEntry struct {
	Key   Key
	Value Value
}

type StoreGroup map[GroupName][]StorageEntry
