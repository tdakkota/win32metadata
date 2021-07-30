package collector

type Field struct {
	Name   string
	GoType string
}

type Constant struct {
	Name   string
	GoType string
	Value  interface{}
}

type Param struct {
	Name   string
	GoType string
}

type Method struct {
	Name   string
	Return Param
	Params []Param

	DLLName   string
	DLLImport string
}

type Import struct {
	Path    string
	Package string
	Def     string
}

type Type struct {
	Name      string
	Namespace string
	Package   string
	IsNewType bool
	Methods   []Method
	Fields    []Field
	Constants []Constant
	Imports   []Import
}
