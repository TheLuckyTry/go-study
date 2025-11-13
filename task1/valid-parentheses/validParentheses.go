package main

import "fmt"

/*
	给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。

有效字符串需满足：

左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。
每个右括号都有一个对应的相同类型的左括号。

示例 1：

输入：s = "()"

输出：true

示例 2：

输入：s = "()[]{}"

输出：true

示例 3：

输入：s = "(]"

输出：false

示例 4：

输入：s = "([])"

输出：true

示例 5：

输入：s = "([)]"

输出：false
*/
func isValid(s string) bool {
	stack := make([]string, 0)
	m1 := map[string]string{
		")": "(",
		"}": "{",
		"]": "[",
	}

	for _, str := range s {
		if str == '(' || str == '{' || str == '[' {
			stack = append(stack, string(str))
		} else {
			if len(stack) == 0 {
				return false
			}
			top := stack[len(stack)-1]

			if top != m1[string(str)] {
				return false
			}

			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func main() {
	fmt.Println("示例1:", isValid("()"))
	fmt.Println("示例2:", isValid("()[]{}"))
	fmt.Println("示例3:", isValid("(]"))
	fmt.Println("示例4:", isValid("([])"))
	fmt.Println("示例5:", isValid("([)]"))

}
