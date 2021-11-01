package murshort

import (
	"strconv"

	"github.com/spaolacci/murmur3"
)

var tenToAny map[int]string = map[int]string{
	0: "A", 1: "B", 2: "C", 3: "D", 4: "E", 5: "F",
	6: "G", 7: "H", 8: "I", 9: "J", 10: "K",
	11: "L", 12: "M", 13: "N", 14: "O", 15: "P",
	16: "Q", 17: "R", 18: "S", 19: "T", 20: "U",
	21: "V", 22: "W", 23: "X", 24: "Y", 25: "Z",
	26: "a", 27: "b", 28: "c", 29: "d", 30: "e",
	31: "f", 32: "g", 33: "h", 34: "i", 35: "j",
	36: "k", 37: "l", 38: "m", 39: "n", 40: "o",
	41: "p", 42: "q", 43: "r", 44: "s", 45: "t",
	46: "u", 47: "v", 48: "w", 49: "x", 50: "y", 51: "z",
	52: "0", 53: "1", 54: "2", 55: "3", 56: "4",
	57: "5", 58: "6", 59: "7", 60: "8", 61: "9",
}

// 10进制转任意进制
func decimalToAny(num, n int) string {
	new_num_str := ""
	var remainder int
	var remainder_string string
	for num != 0 {
		remainder = num % n
		if 76 > remainder && remainder > 9 {
			remainder_string = tenToAny[remainder]
		} else {
			remainder_string = strconv.Itoa(remainder)
		}
		new_num_str = remainder_string + new_num_str
		num = num / n
	}
	return new_num_str
}

func Mur3h62(lng string) string {
	table := murmur3.Sum32([]byte(lng))
	tinyUrl := decimalToAny(int(table), 62)
	return tinyUrl
}
