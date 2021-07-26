package main

import (
	"context"
	"debug/pe"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/tdakkota/win32metadata/types"
)

func run(context.Context) error {
	fileName := flag.String("file", "", "path to metadata file")
	methodName := flag.String("method", "", "method to print")
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

	method, err := findMethod(c, *methodName)
	if err != nil {
		return err
	}
	toPrint := map[types.TypeDefOrRef]types.TypeDef{}

	s, err := printMethod(c, method, toPrint)
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Println(err)
		return
	}
}
