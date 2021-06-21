package test2

/* Test2: typedesc 2.1
typedesc 2.2 */
// typedesc 2.3
type Test2 struct {
	/*blaaaa*/
	Test2 string `yaml:"test2"`
}

type Test2TypeString string
type Test2TypeSlice []string
type Test2TypePointerString *string
type Test2TypePointerSlice []*string
