//+build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

type Value struct {
	Name    string
	Mask    int
	Value   int
	Flag    int
	Denotes string
}

type Attributes struct {
	Name       string
	Represents string
	Type       string
	Values     []Value
}

const indexTemplate = `// {{ $.Name }} represents {{ $.Represents }}.
type {{ $.Name }} {{ $.Type }}

{{ range $value := .Values }}
// {{ $value.Name }} check {{ $value.Name }} flag. {{ if $value.Denotes }}
// Denotes: {{ $value.Denotes }}{{ end }}
func (f {{ $.Name }}) {{ $value.Name }}() bool {
	{{- if $value.Value }}
	return f == {{ $value.Value }}
	{{- else }}

	{{- if $value.Mask }}
	return f & {{ $value.Mask }} == {{ $value.Flag }}
	{{- else }}
	return f & {{ $value.Flag }} != 0
	{{- end }}
	
	{{- end }}
}
{{ end }}`

func run() error {
	out := &bytes.Buffer{}

	t := template.Must(template.New("gen").Parse(indexTemplate))
	attrs := []Attributes{
		{
			Name:       "AssemblyHashAlgorithm",
			Represents: "II.23.1.1 Values for AssemblyHashAlgorithm",
			Type:       "uint32",
			Values: []Value{
				{Name: "None", Value: 0x0000},
				{Name: "MD5", Value: 0x8003},
				{Name: "SHA1", Value: 0x8004},
			},
		},
		{
			Name:       "AssemblyFlags",
			Represents: "II.23.1.2 Values for AssemblyFlags",
			Type:       "uint32",
			Values: []Value{
				{Name: "PublicKey", Flag: 0x0001, Denotes: "The assembly reference holds the full (unhashed)public key."},
				{Name: "Retargetable", Flag: 0x0100, Denotes: "The implementation of this assembly used at runtime isnot expected to match the version seen at compile time."},
				{Name: "DisableJITcompileOptimizer", Flag: 0x4000},
				{Name: "EnableJITcompileTracking", Flag: 0x8000},
			},
		},
		{
			Name:       "EventAttributes",
			Represents: "II.23.1.4 Flags for events [EventAttributes]",
			Type:       "uint16",
			Values: []Value{
				{Name: "SpecialName", Flag: 0x0200, Denotes: "Event is special."},
				{Name: "RTSpecialName", Flag: 0x0400, Denotes: "CLI provides 'special' behavior, depending upon the name of the event"},
			},
		},
		{
			Name:       "FieldAttributes",
			Represents: "II.23.1.5 Flags for fields [FieldAttributes]",
			Type:       "uint16",
			Values: []Value{
				{Name: "CompilerControlled", Mask: 0x7, Flag: 0x0000, Denotes: "Member, not, referenceable"},
				{Name: "Private", Mask: 0x7, Flag: 0x0001, Denotes: "Accessible only by the parent type"},
				{Name: "FamANDAssem", Mask: 0x7, Flag: 0x0002, Denotes: "Accessible by sub-types only in this Assembly"},
				{Name: "Assembly", Mask: 0x7, Flag: 0x0003, Denotes: "Accessibly by anyone in the Assembly"},
				{Name: "Family", Mask: 0x7, Flag: 0x0004, Denotes: "Accessible only by type and sub-types"},
				{Name: "FamORAssem", Mask: 0x7, Flag: 0x0005, Denotes: "Accessibly by sub-types anywhere, plus anyone in assembly"},
				{Name: "Public", Mask: 0x7, Flag: 0x0006, Denotes: "Accessibly by anyone who has visibility to this scope field contract attributes"},

				{Name: "Static", Flag: 0x0010, Denotes: "Defined on type, else per instance"},
				{Name: "InitOnly", Flag: 0x0020, Denotes: "Field can only be initialized, not written to after init"},
				{Name: "Literal", Flag: 0x0040, Denotes: "Value is compile time constant"},
				{Name: "NotSerialized", Flag: 0x0080, Denotes: "Reserved (to indicate this field should not be serialized when type is remoted)"},
				{Name: "SpecialName", Flag: 0x0200, Denotes: "Field is special"},
				// Interop Attributes
				{Name: "PInvokeImpl", Flag: 0x2000, Denotes: "Implementation is forwarded through PInvoke."},
				// Additional flags
				{Name: "RTSpecialName", Flag: 0x0400, Denotes: "CLI provides 'special' behavior, depending upon the name of the field"},
				{Name: "HasFieldMarshal", Flag: 0x1000, Denotes: "Field has marshalling information"},
				{Name: "HasDefault", Flag: 0x8000, Denotes: "Field has default"},
				{Name: "HasFieldRVA", Flag: 0x0100, Denotes: "Field has RVA"},
			},
		},
		{
			Name:       "FileAttributes",
			Represents: "II.23.1.6 Flags for files [FileAttributes]",
			Type:       "uint32",
			Values: []Value{
				{Name: "ContainsMetaData", Flag: 0x0000, Denotes: "This is not a resource file"},
				{Name: "ContainsNoMetaData", Flag: 0x0001, Denotes: "This is a resource file or other non-metadata-containing file"},
			},
		},
		{
			Name:       "GenericParamAttributes",
			Represents: "II.23.1.7 Flags for Generic Parameters [GenericParamAttributes]",
			Type:       "uint16",
			Values: []Value{
				{Name: "None", Mask: 0x0003, Flag: 0x0000, Denotes: "The generic parameter is non-variant and has no special constraints"},
				{Name: "Covariant", Mask: 0x0003, Flag: 0x0001, Denotes: "The generic parameter is covariant"},
				{Name: "Contravariant", Mask: 0x0003, Flag: 0x0002, Denotes: "The generic parameter is contravariant"},

				{Name: "ReferenceTypeConstraint", Mask: 0x001C, Flag: 0x0004, Denotes: "The generic parameter has the class special constraint"},
				{Name: "NotNullableValueTypeConstraint", Mask: 0x001C, Flag: 0x0008, Denotes: "The generic parameter has the valuetype special constraint"},
				{Name: "DefaultConstructorConstraint", Mask: 0x001C, Flag: 0x0010, Denotes: "The generic parameter has the .ctor special constraint"},
			},
		},
		{
			Name:       "PInvokeAttributes",
			Represents: "II.23.1.8 Flags for ImplMap [PInvokeAttributes]",
			Type:       "uint16",
			Values: []Value{
				{Name: "NoMangle", Flag: 0x0001, Denotes: "PInvoke is to use the member name as specified"},
				// Character set
				{Name: "CharSetNotSpec", Mask: 0x0006, Flag: 0x0000, Denotes: ""},
				{Name: "CharSetAnsi", Mask: 0x0006, Flag: 0x0002, Denotes: ""},
				{Name: "CharSetUnicode", Mask: 0x0006, Flag: 0x0004, Denotes: ""},
				{Name: "CharSetAuto", Mask: 0x0006, Flag: 0x0006, Denotes: ""},

				{Name: "SupportsLastError", Flag: 0x0040, Denotes: "Information about target function. Not relevant for fields"},
				// Calling convention
				{Name: "CallConvPlatformapi", Mask: 0x0700, Flag: 0x0100, Denotes: ""},
				{Name: "CallConvCdecl", Mask: 0x0700, Flag: 0x0200, Denotes: ""},
				{Name: "CallConvStdcall", Mask: 0x0700, Flag: 0x0300, Denotes: ""},
				{Name: "CallConvThiscall", Mask: 0x0700, Flag: 0x0400, Denotes: ""},
				{Name: "CallConvFastcall", Mask: 0x0700, Flag: 0x0500, Denotes: ""},
			},
		},
		{

			Name:       "ManifestResourceAttributes",
			Represents: "II.23.1.9 Flags for ManifestResource [ManifestResourceAttributes]",
			Type:       "uint32",
			Values: []Value{
				{Name: "Public", Mask: 0x0007, Flag: 0x0001, Denotes: "The Resource is exported from the Assembly"},
				{Name: "Private", Mask: 0x0007, Flag: 0x0002, Denotes: "The Resource is private to the Assembly"},
			},
		},
		{
			Name:       "MethodAttributes",
			Represents: "II.23.1.10 Flags for methods [MethodAttributes]",
			Type:       "uint16",
			Values: []Value{
				{Name: "CompilerControlled", Mask: 0x7, Flag: 0x0000, Denotes: "Member, not, referenceable"},
				{Name: "Private", Mask: 0x7, Flag: 0x0001, Denotes: "Accessible only by the parent type"},
				{Name: "FamANDAssem", Mask: 0x7, Flag: 0x0002, Denotes: "Accessible by sub-types only in this Assembly"},
				{Name: "Assembly", Mask: 0x7, Flag: 0x0003, Denotes: "Accessibly by anyone in the Assembly"},
				{Name: "Family", Mask: 0x7, Flag: 0x0004, Denotes: "Accessible only by type and sub-types"},
				{Name: "FamORAssem", Mask: 0x7, Flag: 0x0005, Denotes: "Accessibly by sub-types anywhere, plus anyone in assembly"},
				{Name: "Public", Mask: 0x7, Flag: 0x0006, Denotes: "Accessibly by anyone who has visibility to this scope field contract attributes"},

				{Name: "Static", Flag: 0x0010, Denotes: "Defined on type, else per instance"},
				{Name: "Final", Flag: 0x0020, Denotes: "Method cannot be overridden"},
				{Name: "Virtual", Flag: 0x0040, Denotes: "Method is virtual"},
				{Name: "HideBySig", Flag: 0x0080, Denotes: "Method hides by name+sig, else just by name"},

				{Name: "ReuseSlot", Mask: 0x0100, Flag: 0x0000, Denotes: "Method reuses existing slot in vtable"},
				{Name: "NewSlot", Mask: 0x0100, Flag: 0x0100, Denotes: "Method always gets a new slot in the vtable"},

				{Name: "Strict", Flag: 0x0200, Denotes: "Method can only be overriden if also accessible"},
				{Name: "Abstract", Flag: 0x0400, Denotes: "Method does not provide an implementation"},
				{Name: "SpecialName", Flag: 0x0800, Denotes: "Method is special"},
				// Interop attributes
				{Name: "PInvokeImpl", Flag: 0x2000, Denotes: "Implementation is forwarded through PInvoke"},
				{Name: "UnmanagedExport", Flag: 0x0008, Denotes: "Reserved: shall be zero for conforming implementations"},
				// Additional flags
				{Name: "RTSpecialName", Flag: 0x1000, Denotes: "CLI provides 'special' behavior, depending upon the name of the method"},
				{Name: "HasSecurity", Flag: 0x4000, Denotes: "Method has security associate with it"},
				{Name: "RequireSecObject", Flag: 0x8000, Denotes: "Method calls another method containing security code."},
			},
		},
		{
			Name:       "MethodImplAttributes",
			Represents: "II.23.1.11 Flags for methods [MethodImplAttributes]",
			Type:       "uint16",
			Values: []Value{
				{Name: "IL", Mask: 0x0003, Flag: 0x0000, Denotes: "Method impl is CIL"},
				{Name: "Native", Mask: 0x0003, Flag: 0x0001, Denotes: "Method impl is native"},
				{Name: "OPTIL", Mask: 0x0003, Flag: 0x0002, Denotes: "Reserved: shall be zero in conforming implementations"},
				{Name: "Runtime", Mask: 0x0003, Flag: 0x0003, Denotes: "Method impl is provided by the runtime"},

				{Name: "Unmanaged", Mask: 0x0004, Flag: 0x0004, Denotes: "Method impl is unmanaged, otherwise managed"},
				{Name: "Managed", Mask: 0x0004, Flag: 0x0000, Denotes: "Method impl is managed"},
				// Implementation info and interop
				{Name: "ForwardRef", Flag: 0x0010, Denotes: "Indicates method is defined; used primarily in merge scenarios"},
				{Name: "PreserveSig", Flag: 0x0080, Denotes: "Reserved: conforming implementations can ignore"},
				{Name: "InternalCall", Flag: 0x1000, Denotes: "Reserved: shall be zero in conforming implementations"},
				{Name: "Synchronized", Flag: 0x0020, Denotes: "Method is single threaded through the body"},
				{Name: "NoInlining", Flag: 0x0008, Denotes: "Method cannot be inlined"},
				{Name: "MaxMethodImplVal", Flag: 0xffff, Denotes: "Range check value"},
				{Name: "NoOptimization", Flag: 0x0040, Denotes: "Method will not be optimized when generating native code"},
			},
		},
		{
			Name:       "MethodSemanticsAttributes",
			Represents: "II.23.1.12 Flags for MethodSemantics [MethodSemanticsAttributes]",
			Type:       "uint16",
			Values: []Value{
				{Name: "Setter", Flag: 0x0001, Denotes: "Setter for property"},
				{Name: "Getter", Flag: 0x0002, Denotes: "Getter for property"},
				{Name: "Other", Flag: 0x0004, Denotes: "Other method for property or event"},
				{Name: "AddOn", Flag: 0x0008, Denotes: "AddOn method for event. This refers to the required add_ method for events (§22.13)"},
				{Name: "RemoveOn", Flag: 0x0010, Denotes: "RemoveOn method for event. . This refers to the required remove_ method for events (§22.13)"},
				{Name: "Fire", Flag: 0x0020, Denotes: "Fire method for event. This refers to the optional raise_ method for events (§22.13)"},
			},
		},
		{
			Name:       "ParamAttributes",
			Represents: "II.23.1.13 Flags for params [ParamAttributes]",
			Type:       "uint16",
			Values: []Value{
				{Name: "In", Flag: 0x0001, Denotes: "Param is [In]"},
				{Name: "Out", Flag: 0x0002, Denotes: "Param is [out]"},
				{Name: "Optional", Flag: 0x0010, Denotes: "Param is optional"},
				{Name: "HasDefault", Flag: 0x1000, Denotes: "Param has default value"},
				{Name: "HasFieldMarshal", Flag: 0x2000, Denotes: "Param has FieldMarshal"},
				{Name: "Unused", Flag: 0xcfe0, Denotes: "Reserved: shall be zero in a conforming implementation"},
			},
		},
		{
			Name:       "PropertyAttributes",
			Represents: "II.23.1.14 Flags for properties [PropertyAttributes]",
			Type:       "uint16",
			Values: []Value{
				{Name: "SpecialName", Flag: 0x0200, Denotes: "Property is special"},
				{Name: "RTSpecialName", Flag: 0x0400, Denotes: "Runtime(metadata internal APIs) should check name encoding"},
				{Name: "HasDefault", Flag: 0x1000, Denotes: "Property has default"},
				{Name: "Unused", Flag: 0xe9ff, Denotes: "Reserved: shall be zero in a conforming implementation"},
			},
		},
		{
			Name:       "TypeAttributes",
			Represents: "II.23.1.15 Flags for types [TypeAttributes]",
			Type:       "uint32",
			Values: []Value{
				// Visibility attributes
				{Name: "NotPublic", Mask: 0x7, Flag: 0x0, Denotes: "Class has no public scope"},
				{Name: "Public", Mask: 0x7, Flag: 0x1, Denotes: "Class has public scope"},
				{Name: "NestedPublic", Mask: 0x7, Flag: 0x2, Denotes: "Class is nested with public visibility"},
				{Name: "NestedPrivate", Mask: 0x7, Flag: 0x3, Denotes: "Class is nested with private visibility"},
				{Name: "NestedFamily", Mask: 0x7, Flag: 0x4, Denotes: "Class is nested with family visibility"},
				{Name: "NestedAssembly", Mask: 0x7, Flag: 0x5, Denotes: "Class is nested with assembly visibility"},
				{Name: "NestedFamANDAssem", Mask: 0x7, Flag: 0x6, Denotes: "Class is nested with family and assembly visibility"},
				{Name: "NestedFamORAssem", Mask: 0x7, Flag: 0x7, Denotes: "Class is nested with family or assembly visibility"},
				// Class layout attributes
				{Name: "AutoLayout", Mask: 0x18, Flag: 0x0, Denotes: "Class fields are auto-laid out"},
				{Name: "SequentialLayout", Mask: 0x18, Flag: 0x8, Denotes: "Class fields are laid out sequentially"},
				{Name: "ExplicitLayout", Mask: 0x18, Flag: 0x10, Denotes: "Layout is supplied explicitly"},
				// Class semantics attributes
				{Name: "Class", Mask: 0x00000020, Flag: 0x0, Denotes: "Type is a class"},
				{Name: "Interface", Mask: 0x00000020, Flag: 0x20, Denotes: "Type is an interface"},
				// Special semantics in addition to class semantics
				{Name: "Abstract", Flag: 0x80, Denotes: "Class is abstract"},
				{Name: "Sealed", Flag: 0x100, Denotes: "Class cannot be extended"},
				{Name: "SpecialName", Flag: 0x400, Denotes: "Class name is special"},
				// Implementation Attributes
				{Name: "Import", Flag: 0x1000, Denotes: "Class/Interface is imported"},
				{Name: "Serializable", Flag: 0x2000, Denotes: "Reserved (Class is serializable)"},
				// String formatting Attributes
				{Name: "AnsiClass", Mask: 0x30000, Flag: 0x0, Denotes: "LPSTR is interpreted as ANSI"},
				{Name: "UnicodeClass", Mask: 0x30000, Flag: 0x10000, Denotes: "LPSTR is interpreted as Unicode"},
				{Name: "AutoClass", Mask: 0x30000, Flag: 0x20000, Denotes: "LPSTR is interpreted automatically"},
				// FIXME(tdakkota): See II.23.1.15 Flags for types [TypeAttributes], String formatting Attributes
				// Quote:
				// Use this mask to retrieve non-standard encoding information for native interop.
				// The meaning of the values of these 2 bits is unspecified.
				{Name: "CustomFormatClass", Mask: 0x30000, Flag: 0x30000, Denotes: "A non-standard encoding specified by CustomStringFormatMask"},
				// Class Initialization Attributes
				{Name: "BeforeFieldInit", Flag: 0x00100000, Denotes: "Initialize the class before first static field access"},
				// Additional Flags
				{Name: "RTSpecialName", Flag: 0x00000800, Denotes: "CLI provides 'special' behavior, depending upon the name of the Type"},
				{Name: "HasSecurity", Flag: 0x00040000, Denotes: "Type has security associate with it"},
				{Name: "IsTypeForwarder", Flag: 0x00200000, Denotes: "This ExportedType entry is a type forwarder"},
			},
		},
	}

	if _, err := io.WriteString(out, `// Code generated by mk_attributes.go, DO NOT EDIT.
package types

import (
	"fmt"
)

var _ fmt.Stringer
`); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	for _, attr := range attrs {
		for i := range attr.Values {
			if attr.Values[i].Denotes == "" {
				continue
			}
			// Trim if exist and add dot, to ensure dot in the end of comment.
			attr.Values[i].Denotes = strings.TrimRight(attr.Values[i].Denotes, ".") + "."
		}

		if err := t.Execute(out, attr); err != nil {
			return fmt.Errorf("generate %q: %w", attr.Name, err)
		}
	}

	formatted, err := format.Source(out.Bytes())
	if err != nil {
		io.Copy(os.Stderr, out)
		return fmt.Errorf("format: %w", err)
	}

	return os.WriteFile("attributes.gen.go", formatted, 0o600)
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
