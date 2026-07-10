package contracts

type Contract interface {
	Validate() error
}

type SchemaAware interface {
	SchemaVersion() string
}

type Versioned interface {
	Version() string
}

type Identifiable interface {
	ID() string
}

type Named interface {
	Name() string
}
