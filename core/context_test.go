package core

import (
	"fmt"
	"regexp"
	"testing"
)

func TestUpdateString(t *testing.T) {
	str := "L1(A:12970000:${scoring_id:97}${size:5<sort:weight-reatime-desc><size:5><sort:weight-reatime-desc><size:5>)"

	rex := regexp.MustCompile(`\$\{(.*?)\}`)
	out := rex.FindAllStringSubmatch(str, -1)

	for _, i := range out {
		fmt.Println(i[1])
	}
}
