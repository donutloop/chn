package object

const (
	// SortNatural use natural order
	SortNatural Sort = iota
	// SortCreatedDesc created newest to oldest
	SortCreatedDesc
	// SortCreatedAsc created oldest to newset
	SortCreatedAsc
	// SortUpdatedDesc updated newest to oldest
	SortUpdatedDesc
	// SortUpdatedAsc updated oldest to newset
	SortUpdatedAsc
)

type (
	// Sort type for storage handlers
	Sort int

	// Interfaces same as Interface but for slices
	Interfaces interface {
		GetNamespace() string
	}

	// Interface very basic object interface
	Interface interface {
		GetNamespace() string
		GetId() string
		IdSetter
	}

	// TimeTracker help storage handlers set created and updated time when needed.
	TimeTracker interface {
		SetCreated(t int64)
		GetCreated() int64
		SetUpdated(t int64)
		GetUpdated() int64
	}

	// IdSetter helps the storage handler creating new object with fresh uuid
	IdSetter interface {
		SetId(string)
	}

	// ListOpt options for listing objects
	ListOpt struct {
		Page  int64
		Limit int64
		Sort  Sort
	}
)
