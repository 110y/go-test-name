package analysis

type TestInfo struct {
	TestFuncName string   `json:"test_func_name"`
	SubTestNames []string `json:"sub_test_names"`
}
