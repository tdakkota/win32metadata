// Command printsig prints the Go signature of a requested method.
package main

import (
	"debug/pe"
	"flag"
	"fmt"
	"os"

	"github.com/tdakkota/win32metadata/types"
)

func run() error {
	fileName := flag.String("file", "", "path to metadata file")
	methodName := flag.String("method", "", "method to print")
	typeNamespace := flag.String("namespace", "", "method namespace")
	flag.Parse()

	if *methodName == "" {
		return fmt.Errorf("invalid method name: %q", *methodName)
	}

	file, err := pe.Open(*fileName)
	if err != nil {
		return fmt.Errorf("open PE file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	c, err := types.FromPE(file)
	if err != nil {
		return fmt.Errorf("parse metadata: %w", err)
	}

	methodIdx, method, err := findMethod(c, *typeNamespace, *methodName)
	if err != nil {
		return err
	}
	toPrint := map[types.TypeDefOrRef]types.TypeDef{}

	s, err := printMethod(c, methodIdx, method, toPrint)
	if err != nil {
		return fmt.Errorf("print method %q: %w", method.Name, err)
	}
	fmt.Println(s)

	for _, def := range toPrint {
		s, err := printTypeDef(c, def, toPrint)
		if err != nil {
			return fmt.Errorf("print type %q: %w", def.TypeName, err)
		}
		fmt.Println(s)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
}
