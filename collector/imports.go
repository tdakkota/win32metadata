package collector

import (
	"fmt"
	"strconv"
	"strings"
)

func getNamespaces(namespace string) []string {
	namespaces := strings.Split(namespace, ".")
	for i := range namespaces {
		namespaces[i] = strings.ToLower(namespaces[i])
	}
	if len(namespaces) == 1 && namespaces[0] == "" {
		return []string{"mdgen"}
	}

	return namespaces
}

func getPackage(namespace string) string {
	idx := strings.LastIndexByte(namespace, '.')
	if idx < 0 {
		return "other"
	}
	return strings.ToLower(namespace[idx+1:])
}

func getImportName(namespace string) string {
	const pkg = "github.com/tdakkota/win32metadata/collector/mdgen/"
	return strconv.Quote(pkg + strings.Join(getNamespaces(namespace), "/"))
}

type imports struct {
	types      map[namedType]Import
	currentPkg string
}

func (i *imports) Import(def namedType) string {
	pkg := getPackage(def.Namespace)
	if i.currentPkg == pkg {
		return publicName(def.Name)
	}

	i.types[def] = Import{
		Def:     def.Name,
		Package: pkg,
		Path:    getImportName(def.Namespace),
	}
	return fmt.Sprintf("%s.%s", pkg, publicName(def.Name))
}
