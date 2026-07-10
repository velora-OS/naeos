package types

type ID string

type Reference struct {
	ID   ID
	Name string
}

type ErrorInfo struct {
	Code    string
	Message string
}

type Artifact struct {
	Path    string
	Content []byte
}

type Task struct {
	ID           string
	Name         string
	Dependencies []string
	Priority     int
}

type PolicyRule struct {
	RuleID    string
	Condition string
	Priority  int
	Action    string
	Scope     string
}

type KnowledgeEntry struct {
	Topic     string
	Component string
	Version   string
	Rationale string
}

type TelemetryEvent struct {
	Name      string
	Timestamp int64
	Payload   map[string]any
}

type ValidationResult struct {
	Valid  bool
	Errors []ErrorInfo
}

type ReviewResult struct {
	Approved bool
	Comments []string
}

type SpecDocument struct {
	Raw      string
	Project  string
	Modules  []ModuleDef
	Services []ServiceDef
}

type ModuleDef struct {
	Name string
	Path string
}

type ServiceDef struct {
	Name string
	Kind string
	Port int
}
