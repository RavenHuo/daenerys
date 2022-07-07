/**
 * @Author raven
 * @Description
 * @Date 2022/6/30
 **/
package core

import (
	"sort"
)

type Order interface {
	// 设置权重
	Order() int
	// 名字
	Name() string
}

// 根据权重排序
func orderSort(orders ...Order) {
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Order() > orders[j].Order()
	})
}

type MySort struct{}

func (s MySort) Order() int {
	return 1
}
func (s MySort) Name() string {
	return "my sort"
}
