// Command marketprovider 是 `market.Provider` 的最小用法示例（QT-09 ②）。
//
// 它只演示调用形状：五个方法各给一个恒返回零值的桩，真正取数的实现住在服务仓里。
// 下游不需要读本仓的实现就知道「要实现哪些方法、参数与返回值是什么」。
// 对应契约：domain/market 的 Provider（docs/API.md §4 `domain/market/` → Provider）
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ruoxizhnya/tellurion-contracts/domain/market"
)

// stub 是编译期可用的空实现：不做业务计算，只证明方法集完整。
type stub struct{}

func (stub) GetOHLCV(ctx context.Context, symbol string, start, end time.Time) ([]market.OHLCV, error) {
	return nil, nil
}

func (stub) GetFundamental(ctx context.Context, symbol string, date time.Time) (*market.Fundamental, error) {
	return nil, nil
}

func (stub) GetStocks(ctx context.Context, exchange string) ([]market.Stock, error) {
	return nil, nil
}

func (stub) GetLatestPrice(ctx context.Context, symbol string) (float64, error) {
	return 0, nil
}

func (stub) GetIndexConstituents(ctx context.Context, indexCode string) ([]string, error) {
	return nil, nil
}

// 编译期断言：桩实现了端口。方法集变了这里先红，不用等下游实现方。
var _ market.Provider = stub{}

func main() {
	var provider market.Provider = stub{}
	bars, err := provider.GetOHLCV(context.Background(), "600000.SH", time.Time{}, time.Time{})
	fmt.Printf("market.Provider 示例：GetOHLCV 返回 %d 根 bar，err=%v（桩不做业务）\n", len(bars), err)
}
