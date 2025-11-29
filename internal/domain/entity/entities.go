package entity

// GroupName represents the name of a group of store entries, such as aliases or exports.
//
// Values possible:
//
// - entity.Aliases
//
// - entity.Exports
type GroupName string

const (
	Aliases GroupName = "aliases"
	Exports GroupName = "exports"
)

type Key = string
type Value = string

type DbMap map[Key]Value

type DbSnapshot map[GroupName]DbMap

type ShellCommand struct {
	Name string
	Args []string
}
