package internal

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoja(t *testing.T) {
	vm := goja.New()
	err := vm.GlobalObject().Set("load", func(name string) goja.Value {
		fmt.Printf("load %s\n", name)
		return vm.ToValue(map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "test",
				"namespace": "test",
			},
			"data": map[string]interface{}{
				"foo": "bar",
			},
		})
	})
	assert.NoError(t, err)

	err = vm.GlobalObject().Set("apply", func(name string, obj any) (any, error) {
		return obj, nil
	})
	assert.NoError(t, err)

	v, err := vm.RunString(`
var x = load("pera")
x.metadata.name`)
	assert.NoError(t, err)
	fmt.Println(v)

	v, err = vm.RunString(`
var x = {a:1, b:"b"}
apply('name', x)
`)
	assert.NoError(t, err)
	fmt.Println(v)
}
