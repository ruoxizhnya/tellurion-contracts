package broker

// CT-05 ⑤③：wire 字节对拍。
//
// 期望值不是按本包反推的：每一条都由 `tellurion-broker:internal/gateway` 的
// 生产类型 `json.Marshal` 打出来后逐字节抄入（取证命令与输出见
// `tellurion-contracts:docs/TASKS.md` CT-05 施工记录）。改字段名 / tag /
// 字段顺序 / 可空性，这里会立刻红。

import (
	"encoding/json"
	"testing"
	"time"
)

func fixedTS() time.Time { return time.Date(2026, 11, 7, 9, 30, 0, 0, time.UTC) }

func TestWireBytes_MatchBrokerProductionEncoding(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "order_request_limit",
			v: OrderRequest{
				ClientOrderID: "c-1", AccountID: "acc-1", Symbol: "600519",
				Side: SideBuy, OrderType: OrderLimit, Quantity: 100, LimitPrice: 1685.5,
			},
			want: `{"client_order_id":"c-1","account_id":"acc-1","symbol":"600519","side":"buy","order_type":"limit","quantity":100,"limit_price":1685.5}`,
		},
		{
			name: "order_request_market",
			v: OrderRequest{
				ClientOrderID: "c-2", AccountID: "acc-1", Symbol: "600519",
				Side: SideSell, OrderType: OrderMarket, Quantity: 200,
			},
			want: `{"client_order_id":"c-2","account_id":"acc-1","symbol":"600519","side":"sell","order_type":"market","quantity":200}`,
		},
		{
			name: "order_ack_accepted",
			v: OrderAck{
				ClientOrderID: "c-1", BrokerOrderID: "b-9", Status: OrderAccepted,
				IdempotentReplay: false, AcceptedAt: fixedTS(),
			},
			want: `{"client_order_id":"c-1","broker_order_id":"b-9","status":"accepted","idempotent_replay":false,"accepted_at":"2026-11-07T09:30:00Z"}`,
		},
		{
			name: "order_ack_rejected",
			v: OrderAck{
				ClientOrderID: "c-3", Status: OrderRejected, IdempotentReplay: true,
				AcceptedAt: fixedTS(), Code: "order_rejected", Reason: "counterparty rejected",
			},
			want: `{"client_order_id":"c-3","status":"rejected","idempotent_replay":true,"accepted_at":"2026-11-07T09:30:00Z","code":"order_rejected","reason":"counterparty rejected"}`,
		},
		{
			name: "cancel_request_client",
			v:    CancelRequest{ClientOrderID: "c-1"},
			want: `{"client_order_id":"c-1"}`,
		},
		{
			name: "cancel_request_broker",
			v:    CancelRequest{BrokerOrderID: "b-9"},
			want: `{"broker_order_id":"b-9"}`,
		},
		{
			name: "order_state",
			v: OrderState{
				ClientOrderID: "c-1", BrokerOrderID: "b-9", Status: OrderFilled,
				FilledQuantity: 100, AvgFillPrice: 1685.5, RemainingQuantity: 0, LastUpdate: fixedTS(),
			},
			want: `{"client_order_id":"c-1","broker_order_id":"b-9","status":"filled","filled_quantity":100,"avg_fill_price":1685.5,"remaining_quantity":0,"last_update":"2026-11-07T09:30:00Z"}`,
		},
		{
			name: "fill",
			v: Fill{
				FillID: "f-1", ClientOrderID: "c-1", BrokerOrderID: "b-9", Symbol: "600519",
				Side: SideBuy, Quantity: 100, Price: 1685.5, Commission: 4.2,
				StampTax: 0, TransferFee: 0.1, FilledAt: fixedTS(),
			},
			want: `{"fill_id":"f-1","client_order_id":"c-1","broker_order_id":"b-9","symbol":"600519","side":"buy","quantity":100,"price":1685.5,"commission":4.2,"stamp_tax":0,"transfer_fee":0.1,"filled_at":"2026-11-07T09:30:00Z"}`,
		},
		{
			name: "position",
			v: Position{
				AccountID: "acc-1", Symbol: "600519", Quantity: 100, AvailableQuantity: 0,
				AvgCost: 1685.5, MarketValue: 168550, AsOf: fixedTS(),
			},
			want: `{"account_id":"acc-1","symbol":"600519","quantity":100,"available_quantity":0,"avg_cost":1685.5,"market_value":168550,"as_of":"2026-11-07T09:30:00Z"}`,
		},
		{
			name: "account_snapshot_empty_positions",
			v: AccountSnapshot{
				AccountID: "acc-1", AsOf: fixedTS(), Cash: Cash{Total: 1000000, Available: 800000},
				TotalAssets: 1000000, Positions: []Position{}, Fees: Fees{Commission: 4.2, StampTax: 0, TransferFee: 0.1},
			},
			want: `{"account_id":"acc-1","as_of":"2026-11-07T09:30:00Z","cash":{"total":1000000,"available":800000},"total_assets":1000000,"positions":[],"fees":{"commission":4.2,"stamp_tax":0,"transfer_fee":0.1}}`,
		},
		{
			name: "account_snapshot_one_position",
			v: AccountSnapshot{
				AccountID: "acc-1", AsOf: fixedTS(), Cash: Cash{Total: 1000000, Available: 800000},
				TotalAssets: 1168550,
				Positions: []Position{{
					AccountID: "acc-1", Symbol: "600519", Quantity: 100, AvailableQuantity: 0,
					AvgCost: 1685.5, MarketValue: 168550, AsOf: fixedTS(),
				}},
				Fees: Fees{Commission: 4.2, TransferFee: 0.1},
			},
			want: `{"account_id":"acc-1","as_of":"2026-11-07T09:30:00Z","cash":{"total":1000000,"available":800000},"total_assets":1168550,"positions":[{"account_id":"acc-1","symbol":"600519","quantity":100,"available_quantity":0,"avg_cost":1685.5,"market_value":168550,"as_of":"2026-11-07T09:30:00Z"}],"fees":{"commission":4.2,"stamp_tax":0,"transfer_fee":0.1}}`,
		},
		{
			name: "emergency_flatten",
			v: EmergencyFlattenResult{
				Sold: []FlattenedPosition{{
					AccountID: "acc-1", Symbol: "600519", Quantity: 100, ClientOrderID: "c-9",
					BrokerOrderID: "b-9", AvgFillPrice: 1685.5, FilledAt: fixedTS(),
				}},
				Skipped:   []SkippedPosition{{AccountID: "acc-1", Symbol: "000001", Quantity: 200, Reason: "unknown_state"}},
				SoldTotal: 168550,
				StartedAt: fixedTS(), CompletedAt: fixedTS(), Reason: "operator drill",
			},
			want: `{"sold":[{"account_id":"acc-1","symbol":"600519","quantity":100,"client_order_id":"c-9","broker_order_id":"b-9","avg_fill_price":1685.5,"filled_at":"2026-11-07T09:30:00Z"}],"skipped":[{"account_id":"acc-1","symbol":"000001","quantity":200,"reason":"unknown_state"}],"sold_total":168550,"started_at":"2026-11-07T09:30:00Z","completed_at":"2026-11-07T09:30:00Z","reason":"operator drill"}`,
		},
		{
			name: "emergency_flatten_empty",
			v: EmergencyFlattenResult{
				Sold: []FlattenedPosition{}, Skipped: []SkippedPosition{},
				SoldTotal: 0, StartedAt: fixedTS(), CompletedAt: fixedTS(), Reason: "no positions",
			},
			want: `{"sold":[],"skipped":[],"sold_total":0,"started_at":"2026-11-07T09:30:00Z","completed_at":"2026-11-07T09:30:00Z","reason":"no positions"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := json.Marshal(c.v)
			if err != nil {
				t.Fatalf("json.Marshal: %v", err)
			}
			if string(got) != c.want {
				t.Fatalf("wire 字节与 broker 生产侧不一致:\n got: %s\nwant: %s", got, c.want)
			}
		})
	}
}

func TestEnums_BrokerVocabularyIsClosed(t *testing.T) {
	for _, s := range []Side{SideBuy, SideSell} {
		if !s.Valid() {
			t.Fatalf("%q 应合法", s)
		}
	}
	for _, s := range []Side{"", "BUY", "close"} {
		if s.Valid() {
			t.Fatalf("%q 不应合法（Side 只有 buy/sell 两值）", s)
		}
	}
	for _, v := range []OrderType{OrderMarket, OrderLimit} {
		if !v.Valid() {
			t.Fatalf("%q 应合法", v)
		}
	}
	for _, v := range []OrderType{"", "stop", "trailing"} {
		if v.Valid() {
			t.Fatalf("%q 不应合法（broker 不收 stop/trailing）", v)
		}
	}
	all := []OrderStatus{OrderPending, OrderAccepted, OrderPartiallyFilled, OrderFilled,
		OrderCancelled, OrderRejected, OrderExpired, OrderUnknown}
	for _, s := range all {
		if !s.Valid() {
			t.Fatalf("%q 应合法", s)
		}
	}
	if OrderUnknown.Terminal() {
		t.Fatal("unknown 不是终态（它表示还没有结论）")
	}
	for _, s := range []OrderStatus{OrderFilled, OrderCancelled, OrderRejected, OrderExpired} {
		if !s.Terminal() {
			t.Fatalf("%q 应是终态", s)
		}
	}
}
