package engine_test

import (
	"fmt"
	"testing"
)

func TestRT(t *testing.T) {
	//	np := engine.NewProcessed()
	//	chns := make(chan engine.NodeStatistics, 100000)
	//
	//	ctx := context.Background()
	//	ctx, cancel := context.WithCancel(ctx)
	//	go np.Add(ctx, chns)
	//	total := []int{}
	//	for i := 0; i < 100000; i++ {
	//		rand.Seed(time.Now().UnixNano())
	//		rt := rand.Intn(900) + 100
	//		total = append(total, int(rt))
	//		chns <- *engine.NewNodeStatistics(rt, rand.Int63n(100000)+1665735192000, common.Success, "")
	//	}
	// time.Sleep(1 * time.Second)
	//	cancel()
	// sort.Ints(total)
	// r, _ := np.GetPercent([]float64{0.90, 0.50})
	// fmt.Println(r)
	// fmt.Println(total[90000])
	// assert.Equal(t, total[90000], r)
}

func TestCh(t *testing.T) {
	ch := make(chan int, 10000)
	go func() {
		for {
			v, ok := <-ch
			if !ok {
				fmt.Println("ch is closed")
				return
			} else {
				fmt.Println(v)
			}
		}
	}()

	close(ch)

	fmt.Println("closed")
}
