package docu

type Type struct {
	Name  string
	Kinds []*Info
}

type Info struct {
	Path     string
	Kind     string
	Versions []*Version
}

type Version struct {
	Struct   string
	Version  string
	SubKinds map[string]string
}
