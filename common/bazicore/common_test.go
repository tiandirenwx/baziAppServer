package bazicore

import (
	"fmt"
	"testing"
)

var monthIdx689 = []string{"冬", "腊", "正", "二", "三", "四", "五", "六", "七", "八", "九", "九", "十", "正"}
var monthIdx690 = []string{"正", "腊", "一", "二", "三", "四", "五", "六", "七", "八", "九", "十", "正", "腊"}
var monIdx2114 = []string{"冬", "腊", "正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "冬", "腊"}

func Test_Index(t *testing.T) {
	i, _ := Index("正", monthIdx689)
	j, _ := Index("正", monthIdx690)
	k, _ := Index("1", monIdx2114)
	fmt.Printf("Test_Index: month index of 689 is %d, month index of 690 is %d\n", i, j)
	fmt.Printf("Test_Index: month index of 2114 is %d\n", k)

	fmt.Printf("#-------------------------------------------------------------------------------------------#\n")
}

func Test_IndexPlus(t *testing.T) {
	i, _ := IndexPlus("正", monthIdx689, 4, -1)
	j, _ := IndexPlus("正", monthIdx690, 4, -1)
	fmt.Printf("Test_IndexPlus: month index of 689 is %d, month index of 690 is %d\n", i, j)

	i, _ = IndexPlus("正", monthIdx689, 0, -1)
	j, _ = IndexPlus("正", monthIdx690, 0, -1)
	fmt.Printf("Test_IndexPlus: month index of 689 is %d, month index of 690 is %d\n", i, j)

	monIdx2114 := []string{"冬", "腊", "正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "冬", "腊"}
	i, _ = IndexPlus("正", monIdx2114, 4, -1)
	fmt.Printf("Test_IndexPlus: month index of 2114 is %d\n", i)
	fmt.Printf("#-------------------------------------------------------------------------------------------#\n")
}
