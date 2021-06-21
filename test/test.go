package test

import (
	testimp "github.com/caos/documentation/test/test2"
	"github.com/caos/orbos/pkg/kubernetes/k8s"
)

type Test struct {
	TestStruct1 TestStruct `yaml:"testStruct1"`
	TestStruct2 TestStruct `yaml:"testStruct2"`
}

// Test1: typedesc 1.1
// typedesc 1.2
type Test1 struct {
	Attributes      TestAttributes      `yaml:",inline"`
	Imports         TestImports         `yaml:",inline"`
	ExternalImports TestExternalImports `yaml:",inline"`
	Types           TestTypes           `yaml:",inline"`
	Maps            TestMaps            `yaml:",inline"`
}

type TestAttributes struct {
	// @default: testString
	// fielddesc testString
	TestString string `yaml:"testString"`
	// @default: testPointerString
	// fielddesc testPointerString
	TestPointerString *string `yaml:"testPointerString"`
	// @default: - testSlice
	// fielddesc testSlice
	TestSlice []string `yaml:"testSlice"`
	// @default: - testPointerSlice
	// fielddesc testPointerSlice
	TestPointerSlice []*string `yaml:"testPointerSlice"`
}

type TestImports struct {
	// @default: importsTest1
	// fielddesc importsTest1
	Test1 testimp.Test2 `yaml:"importsTest1"`
	// @default: *importsTest2
	// fielddesc importsTest2
	Test2 *testimp.Test2 `yaml:"importsTest2"`
	// @default: - importsTest3
	// fielddesc importsTest3
	TestSlice []testimp.Test2 `yaml:"importsTest3"`
	// @default: - *importsTest4
	// fielddesc importsTest4
	TestPointerSlice []*testimp.Test2 `yaml:"importsTest4"`
}

type TestExternalImports struct {
	//Tolerations to run on nodes
	Tolerations *k8s.Tolerations `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
}

type TestTypes struct {
	// @default: testString
	// fielddesc testString
	test1 TestString `yaml:"test1"`
	// @default: *testString
	// fielddesc *testString
	test2 *TestString `yaml:"test2"`
	// @default: testPointerString
	// fielddesc testPointerString
	test3 TestPointerString `yaml:"test3"`
	// @default: - testString
	// fielddesc []testString
	test4 []TestString `yaml:"test4"`
	// @default: - *testString
	// fielddesc []*TestString
	test5 []*TestString `yaml:"test5"`
	// @default: - testPointerString
	// fielddesc []testPointerString
	test6 []TestPointerString `yaml:"test6`
	// @default: - testSlice
	// fielddesc testSlice
	test7 TestSlice `yaml:"test7"`
	// @default: - testSlicePointer
	// fielddesc testSlicePointer
	test8 TestSlicePointer `yaml:"test8"`
	// @default: testStruct
	// fielddesc testStruct
	test9 TestStruct `yaml:"test9"`
	// @default: *testStruct
	// fielddesc *testStruct
	test10 *TestStruct `yaml:"test10"`
	// @default: testIface1
	// fielddesc testIface1
	TestIface1 interface{} `yaml:"testIface1"`
	// @default: testIface2
	// fielddesc testIface2
	TestIface2 TestInterface `yaml:"testIface2"`
}

type TestMaps struct {
	// @default: string: string
	//fielddesc testMap1
	TestMap1 map[string]string `yaml:"testMap1"`
	// @default: string: *string
	//fielddesc testMap2
	TestMap2 map[string]*string `yaml:"testMap2"`
	// @default: string: testString
	//fielddesc testMap3
	TestMap3 map[string]TestString `yaml:"testMap3"`
	// @default: string: testPointerString
	//fielddesc testMap4
	TestMap4 map[string]TestPointerString `yaml:"testMap4"`
	// @default: string: testStruct
	//fielddesc testMap5
	TestMap5 map[string]TestStruct `yaml:"testMap5"`
	// @default: string: *testStruct
	//fielddesc testMap6
	TestMap6 map[string]*TestStruct `yaml:"testMap6"`
}

type TestString string
type TestPointerString *string
type TestSlice []string
type TestSlicePointer []*string

/* TestStruct: typedesc testStruct */
// typedesc testStruct
type TestStruct struct {
	// fielddesc test4Test1
	Test4Test1 string `yaml:"test4Test1"`
	// fielddesc test4Test2
	Test4Test2 []string `yaml:"test4Test2"`
	// fielddesc test4Test3
	Test4Test3 *string `yaml:"test4Test3"`
	// fielddesc test4Test4
	Test4Test4 []*string `yaml:"test4Test4"`
}

type TestInterface interface{}
