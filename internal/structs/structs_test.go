package structs_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"pkg.glorieux.io/mantra/internal/structs"
)

type testStruct struct{}

func (testStruct) A() {}
func (testStruct) B() {}
func (testStruct) C() {}

func (*testStruct) D() {}

func TestName(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		type TestStruct struct{}
		test := &TestStruct{}
		assert.Equal(t, "TestStruct", structs.Name(test))
	})

	t.Run("Not struct", func(t *testing.T) {
		type TestInt int
		testInt := TestInt(42)
		assert.Equal(t, "TestInt", structs.Name(testInt))
	})

}

func TestMethods(t *testing.T) {
	methods := structs.Methods(&testStruct{})
	assert.Len(t, methods, 4)
	methodNames := func(methods []reflect.Method) (names []string) {
		for _, m := range methods {
			names = append(names, m.Name)
		}
		return
	}
	assert.ElementsMatch(t, methodNames(methods), []string{"A", "B", "C", "D"})
}

func testFunc() {}

func TestFuncName(t *testing.T) {
	t.Run("func", func(t *testing.T) {
		assert.Equal(t, "testFunc", structs.FuncName(testFunc))
	})

	t.Run("method", func(t *testing.T) {
		p := &testStruct{}
		assert.Equal(t, "A", structs.FuncName(p.A))
	})
}
