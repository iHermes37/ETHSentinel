package Analysis

import (
	"fmt"
	"time"
)

func main() {
	stopChans := make(map[string]chan struct{})

	// 启动 CrossDEX
	arbitrage.StartStrategy("CrossDEX", stopChans)

	// 模拟运行一会儿
	time.Sleep(5 * time.Second)

	// 停止 CrossDEX
	close(stopChans["CrossDEX"])
	delete(stopChans, "CrossDEX")

	fmt.Println("done")
}
