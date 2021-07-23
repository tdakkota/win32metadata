package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignatureReader_Method(t *testing.T) {
	// From Windows.Win32.winmd, may be various depending on file
	// Does not matter for testing purposes.
	var (
		NTSTATUS = ElementType{
			Kind:    ELEMENT_TYPE_VALUETYPE,
			TypeDef: ElementTypeTypeDef{Index: TypeDefOrRef(989)},
		}
		HANDLE = ElementType{
			Kind:    ELEMENT_TYPE_VALUETYPE,
			TypeDef: ElementTypeTypeDef{Index: TypeDefOrRef(181)},
		}
		BOOL = ElementType{
			Kind:    ELEMENT_TYPE_VALUETYPE,
			TypeDef: ElementTypeTypeDef{Index: TypeDefOrRef(89)},
		}
		DWORD = ElementType{
			Kind:    ELEMENT_TYPE_VALUETYPE,
			TypeDef: ElementTypeTypeDef{Index: TypeDefOrRef(997)},
		}
		KEY_VALUE_ENTRY = ElementType{
			Kind:    ELEMENT_TYPE_VALUETYPE,
			TypeDef: ElementTypeTypeDef{Index: TypeDefOrRef(4905)},
		}
	)

	tests := []struct {
		name   string
		sig    Signature
		expect MethodSignature
	}{
		{
			"RtlNtStatusToDosError",
			Signature{
				0,            // Calling convention
				1,            // Parameter count
				9,            // Return type
				17, 131, 221, // First param definition (ELEMENT_TYPE_VALUETYPE + TypeDefOrRef index)
			},
			MethodSignature{
				Convention: 0,
				Return: Element{
					Type: ElementType{Kind: ELEMENT_TYPE_U4},
				},
				Params: []Element{
					{Type: NTSTATUS},
				},
			},
		},
		{
			"DuplicateHandle",
			Signature{
				0,
				7,
				17, 89, // BOOL
				17, 128, 181, // HANDLE
				17, 128, 181, // HANDLE
				17, 128, 181, // HANDLE
				15, 17, 128, 181, // *HANDLE = LPHANDLE
				9,      // uint32
				17, 89, // BOOL
				17, 131, 229, // DWORD
			},
			MethodSignature{
				Convention: 0,
				Return: Element{
					Type: BOOL,
				},
				Params: []Element{
					{Type: HANDLE},
					{Type: HANDLE},
					{Type: HANDLE},
					{Type: HANDLE, Pointers: 1},
					{Type: ElementType{Kind: ELEMENT_TYPE_U4}},
					{Type: BOOL},
					{Type: DWORD},
				},
			},
		},
		{
			"NtQueryMultipleValueKey",
			Signature{
				0,
				6,
				17, 131, 221, // NTSTATUS
				17, 128, 181, // HANDLE
				15, 17, 147, 41, // *KEY_VALUE_ENTRY
				9,     // uint32
				15, 1, // *void
				15, 9, // *uint32
				15, 9, // *uint32
			},
			MethodSignature{
				Convention: 0,
				Return: Element{
					Type: NTSTATUS,
				},
				Params: []Element{
					{Type: HANDLE},
					{Type: KEY_VALUE_ENTRY, Pointers: 1},
					{Type: ElementType{Kind: ELEMENT_TYPE_U4}},
					{Type: ElementType{Kind: ELEMENT_TYPE_VOID}, Pointers: 1},
					{Type: ElementType{Kind: ELEMENT_TYPE_U4}, Pointers: 1},
					{Type: ElementType{Kind: ELEMENT_TYPE_U4}, Pointers: 1},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := require.New(t)
			r := test.sig.Reader()

			method, err := r.Method(nil) // Context must not be needed anyway
			a.NoError(err)

			a.Equal(test.expect, method)
		})
	}
}
