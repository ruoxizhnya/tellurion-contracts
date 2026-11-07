package broker

import "time"

// OrderRequest 是下单请求。只保留通道需要的字段。
//
// **不含** `BusinessType` / `ExecAlgo` / `StopPrice` / `Trail*`：券商私有字段留在
// 适配器内部，执行层语义留在 core。
type OrderRequest struct {
	// ClientOrderID 由调用方生成、重试不变；网关不改写、不生成（幂等键）。
	ClientOrderID string    `json:"client_order_id"`
	AccountID     string    `json:"account_id"`
	Symbol        string    `json:"symbol"`
	Side          Side      `json:"side"`
	OrderType     OrderType `json:"order_type"`
	Quantity      int       `json:"quantity"`
	LimitPrice    float64   `json:"limit_price,omitempty"`
}

// OrderAck 是下单受理回执。**受理 ≠ 成交。**
type OrderAck struct {
	ClientOrderID    string      `json:"client_order_id"`
	BrokerOrderID    string      `json:"broker_order_id,omitempty"`
	Status           OrderStatus `json:"status"`
	IdempotentReplay bool        `json:"idempotent_replay"`
	AcceptedAt       time.Time   `json:"accepted_at"`
	Code             string      `json:"code,omitempty"`
	Reason           string      `json:"reason,omitempty"`
}

// CancelRequest 是撤单请求：至少 `ClientOrderID` / `BrokerOrderID` 之一非空，
// 优先按 `ClientOrderID` 幂等。
type CancelRequest struct {
	ClientOrderID string `json:"client_order_id,omitempty"`
	BrokerOrderID string `json:"broker_order_id,omitempty"`
}

// OrderState 是订单的当前状态快照。
type OrderState struct {
	ClientOrderID     string      `json:"client_order_id"`
	BrokerOrderID     string      `json:"broker_order_id,omitempty"`
	Status            OrderStatus `json:"status"`
	FilledQuantity    int         `json:"filled_quantity"`
	AvgFillPrice      float64     `json:"avg_fill_price"`
	RemainingQuantity int         `json:"remaining_quantity"`
	LastUpdate        time.Time   `json:"last_update"`
	Code              string      `json:"code,omitempty"`
	Reason            string      `json:"reason,omitempty"`
}
