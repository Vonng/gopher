package format

import (
	"testing"
	"fmt"
)

func TestChineseNumberToInt(t *testing.T) {
	ChineseNumberToInt("1.3亿")
}

type TestCase struct {
	Input  string
	Expect int64
}

func TestChineseSuffixStringToInt(t *testing.T) {
	testCase := []TestCase{
		TestCase{"2.56万", 25600},
		TestCase{"256万", 2560000},
		TestCase{"256.万", 2560000},
		TestCase{"25.6万", 256000},
		TestCase{"0.256万", 2560},
		TestCase{"1.256万", 12560},
		TestCase{"1.2567万", 12567},
		TestCase{"1.25678万", 12567},
		TestCase{"1.25678亿", 125678000},
		TestCase{"1.25678901亿", 125678901},
		TestCase{"1.256789019亿", 125678901},
		TestCase{"0.00001亿", 1000},
	}

	for _, c := range testCase {
		output, err := ChineseSuffixStringToInt(c.Input)
		if err != nil {
			t.Error(err)
		}

		if output != c.Expect {
			t.Errorf("Input[%s] Expect[%d] Got[%d]\n", c.Input, c.Expect, output)
		}
		fmt.Printf("Input[%s] Expect[%d] = Got[%d]\n", c.Input, c.Expect, output)
	}
}

func TestPrefixedBytesToInt(t *testing.T) {
	testCase := []TestCase{
		TestCase{"256", 256},
		TestCase{"256.128", 256},
		TestCase{"2KB", 2048},
		TestCase{"2.56KB", 2621},
		TestCase{"1024K", 1048576},
		TestCase{"2M", 2097152},
		TestCase{"2.5M", 2621440},
		TestCase{"2.5432M", 2666738},
	}

	for _, c := range testCase {
		output, err := PrefixedBytesToInt(c.Input)
		if err != nil {
			t.Error(err)
		}

		if output != c.Expect {
			t.Errorf("Input[%s] Expect[%d] Got[%d]\n", c.Input, c.Expect, output)
		}
		fmt.Printf("Input[%s] Expect[%d] = Got[%d]\n", c.Input, c.Expect, output)
	}
}
