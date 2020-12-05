# jpholiday

日本の祝日（国民の祝日）を判定するパッケージ。


## Example

	import (
		"time"
		"jpholiday"
		"fmt"
	)

	func main() {
		day := time.Parse("2006-01-02", "2013-11-23")
		name := jpholiday.Name(day)
		fmt.Println(name)    // 勤労感謝の日
	}


## Changelog

* 20201205: 2021年の東京オリンピック特別措置法の祝日に対応


## Lisence

New BSD License.
