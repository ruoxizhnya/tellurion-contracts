// Package broker 是交易通道（tellurion-broker）的跨进程 wire 形状。
//
// 边界（CT-05 裁定）：
//   - 只收**报文里出现的值**：订单请求 / 回执 / 状态、成交、持仓、账户快照、
//     紧急清仓结果及其枚举；不含 XTP 适配器内部类型（BusinessType / XTPConfig 等）。
//   - `Quantity` 是**整股**（`int`），与 `domain.Position.Quantity` 的 `float64`
//     不是同一语义，**不许**用类型别名把两侧等同。
//   - `Side` 只有 buy / sell 两值，`OrderType` 只有 market / limit 两值；
//     `domain.Direction` 四值、`domain.OrderType` 四值的映射是**业务判断**，
//     留在消费方，不藏进 JSON 字段名。
//   - 12 个契约哨兵**不随迁**（详见 tellurion-contracts:docs/TASKS.md CT-05 施工记录）：
//     它们是 broker 进程内的错误值，HTTP 消费方按键值码（wire code）判定。
package broker

// Side 是买卖方向（两值）。
type Side string

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"
)

// Valid 报告 s 是否是已知方向。
func (s Side) Valid() bool { return s == SideBuy || s == SideSell }

// OrderType 是委托类型（只两种：市价 / 限价）。
type OrderType string

const (
	OrderMarket OrderType = "market"
	OrderLimit  OrderType = "limit"
)

// Valid 报告 t 是否是已知委托类型。
func (t OrderType) Valid() bool { return t == OrderMarket || t == OrderLimit }

// OrderStatus 是订单状态八值（**不是 run 三值**）。
//
// `unknown` 是**未决**状态，不是拒绝：调用方只能查单 / 对账，禁止重发。
type OrderStatus string

const (
	OrderPending         OrderStatus = "pending"
	OrderAccepted        OrderStatus = "accepted"
	OrderPartiallyFilled OrderStatus = "partially_filled"
	OrderFilled          OrderStatus = "filled"
	OrderCancelled       OrderStatus = "cancelled"
	OrderRejected        OrderStatus = "rejected"
	OrderExpired         OrderStatus = "expired"
	OrderUnknown         OrderStatus = "unknown"
)

// Valid 报告 s 是否是八值之一。
func (s OrderStatus) Valid() bool {
	switch s {
	case OrderPending, OrderAccepted, OrderPartiallyFilled, OrderFilled,
		OrderCancelled, OrderRejected, OrderExpired, OrderUnknown:
		return true
	}
	return false
}

// Terminal 报告 s 是否是终态（filled / cancelled / rejected / expired）。
// `unknown` **不是**终态：它表示还没有结论。
func (s OrderStatus) Terminal() bool {
	switch s {
	case OrderFilled, OrderCancelled, OrderRejected, OrderExpired:
		return true
	}
	return false
}
