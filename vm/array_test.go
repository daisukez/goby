package vm

import (
	"testing"
)

func TestArrayClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Array.class.name`, "Class"},
		{`Array.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestArrayEvaluation(t *testing.T) {
	input := `
	[1, "234", true]
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input, getFilename())
	vm.checkCFP(t, 0, 0)

	arr, ok := evaluated.(*ArrayObject)
	if !ok {
		t.Fatalf("Expect evaluated value to be an array. got=%T", evaluated)
	}

	checkExpected(t, 0, arr.Elements[0], 1)
	checkExpected(t, 0, arr.Elements[1], "234")
	checkExpected(t, 0, arr.Elements[2], true)
}

func TestArrayComparisonOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`[1, "String", true, 2..5] == 123`, false},
		{`[1, "String", true, 2..5] == "123"`, false},
		{`[1, "String", true, 2..5] == "124"`, false},
		{`[1, "String", true, 2..5] == (1..3)`, false},
		{`[1, "String", true, 2..5] == { a: 1, b: 2 }`, false},
		{`[1, "String", true, 2..5] == [1, "String", true, 2..5]`, true},
		{`[1, "String", true, 2..5] == [1, "String", false, 2..5]`, false},
		{`[1, "String", true, 2..5] == ["String", 1, false, 2..5]`, false}, // Array has order issue
		{`[1, { a: 1, b: 2 }, "Goby" ] == [1, { a: 1, b: 2 }, "Goby"]`, true},
		{`[1, { a: 1, b: 2 }, "Goby" ] == [1, { b: 2, a: 1 }, "Goby"]`, true},
		{`[1, { a: 1, b: 2 }, "Goby" ] == [1, { a: 1, b: 2, c: 3 }, "Goby"]`, false}, // Array of hash has no order issue
		{`[1, { a: 1, b: 2 }, "Goby" ] == [1, { a: 2, b: 2, a: 1 }, "Goby"]`, true},  // Array of hash key will be overwritten if duplicated
		{`[1, "String", true, 2..5] == Integer`, false},
		{`[1, "String", true, 2..5] != 123`, true},
		{`[1, "String", true, 2..5] != "123"`, true},
		{`[1, "String", true, 2..5] != "124"`, true},
		{`[1, "String", true, 2..5] != (1..3)`, true},
		{`[1, "String", true, 2..5] != { a: 1, b: 2 }`, true},
		{`[1, "String", true, 2..5] != [1, "String", true, 2..5]`, false},
		{`[1, "String", true, 2..5] != [1, "String", false, 2..5]`, true},
		{`[1, "String", true, 2..5] != ["String", 1, false, 2..5]`, true}, // Array has order issue
		{`[1, { a: 1, b: 2 }, "Goby" ] != [1, { a: 1, b: 2 }, "Goby"]`, false},
		{`[1, { a: 1, b: 2 }, "Goby" ] != [1, { b: 2, a: 1 }, "Goby"]`, false},
		{`[1, { a: 1, b: 2 }, "Goby" ] != [1, { a: 1, b: 2, c: 3 }, "Goby"]`, true},  // Array of hash has no order issue
		{`[1, { a: 1, b: 2 }, "Goby" ] != [1, { a: 2, b: 2, a: 1 }, "Goby"]`, false}, // Array of hash key will be overwritten if duplicated
		{`[1, "String", true, 2..5] != Integer`, true},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestArrayIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
			[][1]
		`, nil},
		{`
			[1, 2, 3][100]
		`, nil},
		{`
			[1, 2, 10, 5][2]
		`, 10},
		{`
			[1, "a", 10, 5][1]
		`, "a"},
		{`
		    [1, "a", 10, "b"][-2]
		`, 10},
		{`
		    [1, "a", 10, "b"][-5]
		`, nil},
		{`
			a = [1, "a", 10, 5]
			a[0]
		`, 1},
		{`
			a = [1, "a", 10, 5]
			a[2] = a[1]
			a[2]
		`, "a"},
		{`
			a = [1, "a", 10, 5]
			a[-2] = a[1]
			a[-2]
		`, "a"},
		{`
			a = [1, "a", 10, 5]
			a[-4] = a[1]
			a[-4]
		`, "a"},
		{`
			a = []
			a[10] = 100
			a[10]
		`, 100},
		{`
			a = []
			a[10] = 100
			a[0]
		`, nil},
		{`
			a = [1, 2 ,3 ,5 , 10]
			a[0] = a[1] + a[2] + a[3] * a[4]
			a[0]
		`, 55},
		{
			`
			code = []
			code[100] = 'Continue'
			code[101] = 'Switching Protocols'
			code[102] = 'Processing'
			code[200] = 'OK'
			code.to_s
			`,
			`[nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, "Continue", "Switching Protocols", "Processing", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, "OK"]`},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
	}
}

func TestArrayAtMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
			[].at(1)
		`, nil},
		{`
			[1, 2, 10, 5].at(2)
		`, 10},
		{`
			[1, "a", 10, 5].at(1)
		`, "a"},
		{`
			[1, "a", 10, 5].at(4)
		`, nil},
		{`
			[1, "a", 10, 5].at(-2)
		`, 10},
		{`
			[1, "a", 10, 5].at(-5)
		`, nil},
		{`
			a = [1, "a", 10, 5]
			a.at(0)
		`, 1},
		{`
			a = [1, "a", 10, 5]
			a[2] = a.at(1)
			a[2]
		`, "a"},
		{`
			a = [1, 2, 3, 5, 10]
			a[0] = a.at(1) + a.at(2) + a.at(3) * a.at(4)
			a.at(0)
		`, 55},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
	}
}

func TestArrayClearMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		a = [1, 2, 3]
		a.clear
		`, []interface{}{}},
		{`
		a = []
		a.clear
		`, []interface{}{}},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestArrayConcatMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		a = [1, 2]
		a.concat([3], [4])
		`, []interface{}{1, 2, 3, 4}},
		{`
		a = []
		a.concat([1], [2], ["a", "b"], [3], [4])
		`, []interface{}{1, 2, "a", "b", 3, 4}},
		{`
		a = [1, 2]
		a.concat()
		`, []interface{}{1, 2}},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestArrayConcatMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`a = [1, 2]
		a.concat(3)
		`, "TypeError: Expect argument to be Array. got: Integer", 2},
		{`a = []
		a.concat("a")
		`, "TypeError: Expect argument to be Array. got: String", 2},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestArrayCountMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		a = [1, 2]
		a.count
		`, 2},
		{`
		a = [1, 2]
		a.count(1)
		`, 1},
		{`
		a = ["a", "bb", "c", "db", "bb", 2]
		a.count("bb")
		`, 2},
		{`
		a = [true, true, true, false, true]
		a.count(true)
		`, 4},
		{`
		a = []
		a.count(true)
		`, 0},
		{`
		a = [1, 2, 3, 4, 5, 6, 7, 8]
		a.count do |i|
			i > 3
		end
		`, 5},
		{`
		a = ["a", "bb", "c", "db", "bb"]
		a.count do |i|
			i.size > 1
		end
		`, 3},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
	}
}

func TestArrayCountMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{
			`a = [1, 2]
		a.count(3, 3)
		`, "ArgumentError: Expect 1 argument, got=2", 2},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestArrayEachMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		sum = 0
		[1, 2, 3, 4, 5].each do |i|
		  sum = sum + i
		end
		sum
		`, 15},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
	}
}

func TestArrayEachIndexMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		sum = 0
		[2, 3, 40, 5, 22].each_index do |i|
		  sum = sum + i
		end
		sum
		`, 10},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
	}
}

func TestArrayEmptyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			`
			[1, 2, 3].empty?
			`, false},
		{
			`
			[nil].empty?
			`, false},
		{
			`
			[].empty?
			`, true},
		{
			`
			a = [[]]
			a.empty?
			`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestArrayFirstMethod(t *testing.T) {
	testsInt := []struct {
		input    string
		expected interface{}
	}{
		{`
		a = [1, 2]
		a.first
		`, 1},
	}

	for i, tt := range testsInt {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
	}

	testsArray := []struct {
		input    string
		expected []interface{}
	}{
		{`
		a = [3, 4, 5, 1, 6]
		a.first(2)
		`, []interface{}{3, 4}},
		{`
		a = ["a", "b", "d", "q"]
		a.first(2)
		`, []interface{}{"a", "b"}},
	}

	for i, tt := range testsArray {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
	}
}

func TestArrayFirstMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`a = [1, 2]
		a.first("a")
		`, "TypeError: Expect argument to be Integer. got: String", 2},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestArrayFlattenMethod(t *testing.T) {
	testsArray := []struct {
		input    string
		expected []interface{}
	}{
		{`
		[1, 2].flatten
		`, []interface{}{1, 2}},
		{`
		[1, 2, [3, 4, 5]].flatten
		`, []interface{}{1, 2, 3, 4, 5}},
		{`
		[[1, 2], [3, 4], [5, 6]].flatten
		`, []interface{}{1, 2, 3, 4, 5, 6}},
		{`
		[[[1, 2], [[[3, 4]], [5, 6]]]].flatten
		`, []interface{}{1, 2, 3, 4, 5, 6}},
	}

	for i, tt := range testsArray {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestArrayFlattenMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`a = [1, 2]
		a.flatten(1)
		`, "ArgumentError: Expect 0 argument. got=1", 2},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestArrayJoinMethod(t *testing.T) {
	testsInt := []struct {
		input    string
		expected string
	}{
		{`
		[1, 2].join
		`, "12"},
		{`
		["1", 2].join
		`, "12"},
		{`
		[1, 2].join(",")
		`, "1,2"},
		{`
		[1, 2, [3, 4]].join(",")
		`, "1,2,3,4"},
	}

	for i, tt := range testsInt {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
	}
}

func TestArrayJoinMethodFailWithArgumentError(t *testing.T) {
	testsFail := []errorTestCase{
		{`a = [1, 2]
		a.join(",", "-")
		`, "ArgumentError: Expect 0 or 1 argument. got=2", 2},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestArrayJoinMethodFailWithTypeError(t *testing.T) {
	testsFail := []errorTestCase{
		{`a = [1, 2]
		a.join(1)
		`, "TypeError: Expect argument to be String. got: Integer", 2},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestArrayLastMethod(t *testing.T) {
	testsArray := []struct {
		input    string
		expected []interface{}
	}{
		{`
		a = [3, 4, 5, 1, 6]
		a.last(3)
		`, []interface{}{5, 1, 6}},
		{`
		a = ["a", "b", "d", "q"]
		a.last(2)
		`, []interface{}{"d", "q"}},
	}

	for i, tt := range testsArray {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestArrayLastMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`a = [1, 2]
		a.last("l")
		`, "TypeError: Expect argument to be Integer. got: String", 2},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
	}
}

func TestArrayLengthMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{
			`
			[1, 2, 3].length
			`, 3},
		{
			`
			[nil].length
			`, 1},
		{
			`
			[].length
			`, 0},
		{
			`
			a = [-10, "123", [1,2,3], 1, 2, 3]
			a.length
			`, 6},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestArrayMapMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		a = [1, 2, 7]
		a.map do |i|
			i + 3
		end
		`, []interface{}{4, 5, 10}},
		{`
		a = [true, false, true, false, true ]
		a.map do |i|
			!i
		end
		`, []interface{}{false, true, false, true, false}},
		{`
		a = ["1", "sss", "qwe"]
		a.map do |i|
			i + "1"
		end
		`, []interface{}{"11", "sss1", "qwe1"}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestArrayPopMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`
			a = [1, 2, 3].pop
			a
			`, 3},
		{
			`
			a = [1, 2, 3]
			a.pop
			a.length
			`, 2},
		{
			`
			[].pop
		`, nil},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestArrayPushMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`
			a = [1, 2, 3]
			a.push("test")
			a[3]
			`, "test"},
		{
			`
			a = [1, 2, 3]
			a.push("test")
			a.length
			`, 4},
		{
			`
			a = []
			a.push(nil)
			a[0]
			`, nil},
		{
			`
			a = []
			a.push("foo")
			a.push(1)
			a.push(234)
			a[0]
			`, "foo"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestArrayReduceMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		a = [1, 2, 7]
		a.reduce do |sum, n|
			sum + n
		end
		`, 10},
		{`
		a = [1, 2, 7]
		a.reduce(10) do |sum, n|
			sum + n
		end
		`, 20},
		{`
		a = ["This ", "is a ", "test!"]
		a.reduce do |prev, s|
			prev + s
		end
		`, "This is a test!"},
		{`
		a = ["this ", "is a ", "test!"]
		a.reduce("Yes, ") do |prev, s|
			prev + s
		end
		`, "Yes, this is a test!"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestArrayReduceMethodFailWithInternalError(t *testing.T) {
	testsFail := []errorTestCase{
		{`a = [1, 2]
		a.reduce(1)
		`,
			"InternalError: Can't yield without a block",
			2},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestArrayReduceMethodFailWithArgumentError(t *testing.T) {
	testsFail := []errorTestCase{
		{`a = [1, 2]
		a.reduce(1, 2) do |prev, n|
			prev + n
		end
		`,
			"ArgumentError: Expect 0 or 1 argument. got=2",
			2},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestArrayRotateMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		a = [1, 2]
		a.rotate
		`, []interface{}{2, 1}},
		{`
		a = [1, 2, 3, 4]
		a.rotate(2)
		`, []interface{}{3, 4, 1, 2}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestArrayRotateMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`a = [1, 2]
		a.rotate("a")
		`,
			"TypeError: Expect argument to be Integer. got: String",
			2},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)

		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestArraySelectMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		a = [1, 2, 3, 4, 5]
		a.select do |i|
			i > 3
		end
		`, []interface{}{4, 5}},
		{`
		a = [true, false, true, false, true ]
		a.select do |i|
			i
		end
		`, []interface{}{true, true, true}},
		{`
		a = ["test", "not2", "3", "test", "5"]
		a.select do |i|
			i == "test"
		end
		`, []interface{}{"test", "test"}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestArrayShiftMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`
			a = [1, 2, 3].shift
			a
			`, 1},
		{
			`
			a = [1, 2, 3]
			a.pop
			a.length
			`, 2},
		{
			`
				[].shift
			`, nil},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestArrayShiftMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`a = [1, 2]
		a.shift(3, 3, 4, 5)
		`,
			"ArgumentError: Expect 0 argument. got=4",
			2},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)

		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}
