package filter

import (
	"fmt"
	"testing"
)

func TestUtil(t *testing.T) {
	util := NewRouteUtil()
	var hosts = []string{"b.keruyun.com"}
	var paths = []string{"/**"}
	util.parseHostPath(hosts, paths, 1)
	var hosts2 = []string{"b.keruyun.com"}
	var paths2 = []string{"/aa/**"}
	util.parseHostPath(hosts2, paths2, 2)
	var hosts3 = []string{"b.keruyun.com"}
	var paths3 = []string{"/aa/ab/**"}
	util.parseHostPath(hosts3, paths3, 3)
	var hosts4 = []string{"b.keruyun.com"}
	var paths4 = []string{"/aa/ab"}
	util.parseHostPath(hosts4, paths4, 4)
	var hosts5 = []string{"kpib.keruyun.com"}
	var paths5 = []string{"/aa/**"}
	util.parseHostPath(hosts5, paths5, 5)
	var hosts6 = []string{"b.keruyun.com"}
	var paths6 = []string{"/aa/ab/ac/**", "//aaa"}
	util.parseHostPath(hosts6, paths6, 6)

	matchService(1, "b.keruyun.com/aa", *util, t)
	matchService(1, "b.keruyun.com/ab", *util, t)
	matchService(1, "kpib.keruyun.com/ab", *util, t)
	matchService(2, "b.keruyun.com/aa/ac", *util, t)
	matchService(5, "kpib.keruyun.com/aa/ac", *util, t)
	matchService(4, "b.keruyun.com/aa/ab", *util, t)
	matchService(3, "b.keruyun.com/aa/ab/ac", *util, t)
	matchService(4, "b.keruyun.com/aa/ab?sfaf=sda&asd=23", *util, t)
	matchService(6, "b.keruyun.com/aa/ab/ac/ad", *util, t)
	matchService(6, "aaa", *util, t)

}
func matchService(expectedValue int, key string, util RouteUtil, t *testing.T) {
	getValue := getSvc(key, util)
	if getValue == nil {
		t.Errorf("url:%v expected:%v, got:null", key, expectedValue)
		return
	}
	if expectedValue != getValue.(int) {
		t.Errorf("url:%v expected:%v, got:%v", key, expectedValue, getValue.(int))
	}
}
func getSvc(key string, util RouteUtil) interface{} {
	svc := util.MatchService(key)
	if svc != nil {
		fmt.Println("match ", key, " ", svc.(int))
	} else {
		fmt.Println(key + " not found")
	}
	return svc
}
