package core

import(
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestVM(t *testing.T){
	//push FOO to the stack (will be used as key)
	//push 3 to the stack
	//push 2 to the stack
	// perform 3-2
	// 1 is on the stack
	// [FOO, 1] on the stack
	//store (will store 1 on key "FOO")
	//data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0, 0x0d, 0x03,0x0a, 0x02, 0x0a, 0x0e}
	// F O O ==> [F O O]
	//data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0, 0x0d}

	//data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f} //store 5 with key "FOO" on the contractState Array
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())
	fmt.Printf("%+v\n", contractState)
	valueBytes, err := contractState.Get([]byte("FOO"))
	value := deserializeInt64(valueBytes)
	assert.Nil(t, err)
	assert.Equal(t, value, int64(5))
}

func TestVMGet(t *testing.T){
	pushFoo := []byte{0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0xae}
	// 2 3 + = 5 = FOO FOO=5
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}
	data = append(data, pushFoo...)
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())
	value := vm.stack.Pop().([]byte)
	valueSerialized := deserializeInt64(value)
	assert.Equal(t, valueSerialized, int64(5))
}



func TestVMMul(t *testing.T){
	//Multiple 3 by 2 = 6
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0xea}
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())
	result := vm.stack.Pop().(int)
	assert.Equal(t, result, 6)
}

func TestVMDiv(t *testing.T){
	// divide 6 by 2 = 3
	data := []byte{0x02, 0x0a, 0x06, 0x0a, 0xfd} 
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())
	result := vm.stack.Pop().(int)
	assert.Equal(t, result, 3)
}


func TestStack(t *testing.T) {
	s := NewStack(128)

	s.Push(1)
	s.Push(2)

	value:= s.Pop()
	assert.Equal(t, value, 2)
	value = s.Pop()
	assert.Equal(t, value, 1)
}