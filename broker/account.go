package broker

import "time"

// Fill 是归一化成交回报。`FillID` 稳定、跨调用不重复。
type Fill struct {
	FillID        string    `json:"fill_id"`
	ClientOrderID string    `json:"client_order_id"`
	BrokerOrderID string    `json:"broker_order_id,omitempty"`
	Symbol        string    `json:"symbol"`
	Side          Side      `json:"side"`
	Quantity      int       `json:"quantity"`
	Price         float64   `json:"price"`
	Commission    float64   `json:"commission"`
	StampTax      float64   `json:"stamp_tax"`
	TransferFee   float64   `json:"transfer_fee"`
	FilledAt      time.Time `json:"filled_at"`
}

// Position 是柜台口径的持仓。`AvailableQuantity` 是柜台事实；
// 本地 T+1 记账字段不进 DTO。
type Position struct {
	AccountID         string    `json:"account_id"`
	Symbol            string    `json:"symbol"`
	Quantity          int       `json:"quantity"`
	AvailableQuantity int       `json:"available_quantity"`
	AvgCost           float64   `json:"avg_cost"`
	MarketValue       float64   `json:"market_value"`
	AsOf              time.Time `json:"as_of"`
}

// Cash 是资金口径快照。
type Cash struct {
	Total     float64 `json:"total"`
	Available float64 `json:"available"`
}

// Fees 是费用累计口径。
type Fees struct {
	Commission  float64 `json:"commission"`
	StampTax    float64 `json:"stamp_tax"`
	TransferFee float64 `json:"transfer_fee"`
}

// AccountSnapshot 是账户快照。账户查询与对账查询**共用同一形状**。
type AccountSnapshot struct {
	AccountID   string     `json:"account_id"`
	AsOf        time.Time  `json:"as_of"`
	Cash        Cash       `json:"cash"`
	TotalAssets float64    `json:"total_assets"`
	Positions   []Position `json:"positions"`
	Fees        Fees       `json:"fees"`
}

// FlattenedPosition 是一条**已真实提交并确认**的平仓成交。
type FlattenedPosition struct {
	AccountID     string    `json:"account_id"`
	Symbol        string    `json:"symbol"`
	Quantity      int       `json:"quantity"`
	ClientOrderID string    `json:"client_order_id"`
	BrokerOrderID string    `json:"broker_order_id,omitempty"`
	AvgFillPrice  float64   `json:"avg_fill_price"`
	FilledAt      time.Time `json:"filled_at"`
}

// SkippedPosition 是一条**没有确认成交**的持仓（含未提交 / 被拒 / 状态未知）。
type SkippedPosition struct {
	AccountID string `json:"account_id"`
	Symbol    string `json:"symbol"`
	Quantity  int    `json:"quantity"`
	Reason    string `json:"reason"`
}

// EmergencyFlattenResult 是紧急清仓结果。`Sold` 与 `Skipped` 互斥：
// 每个 `Sold` 都必须有真实成交依据，未确认的一律进 `Skipped`。
type EmergencyFlattenResult struct {
	Sold        []FlattenedPosition `json:"sold"`
	Skipped     []SkippedPosition   `json:"skipped"`
	SoldTotal   int                 `json:"sold_total"`
	StartedAt   time.Time           `json:"started_at"`
	CompletedAt time.Time           `json:"completed_at"`
	Reason      string              `json:"reason"`
}
