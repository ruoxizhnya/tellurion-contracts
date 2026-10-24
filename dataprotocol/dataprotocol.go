// Package dataprotocol — 数据面的共享协议形状（点位 / 去重坐标 / 序列化）。
//
// 2026-10-24 随 QT-40 X1a 新建：把三处已经漂开的同名形状收成一份定义 ——
// `tellurion-data:internal/ingest/point.go`（带 VisibleAt、时钟由调用方注入）、
// `tellurion:pkg/data/source/unified.go`（有 DeduplicateKey、构造器读墙钟）、
// `tellurion:pkg/storage/bulk_insert.go`（存储投影，因 import 成环而另抄一份）。
//
// 三条边界：
//  1. 零 I/O、零时钟：本包不提供构造器 —— IngestTime 由调用方注入，否则
//     「同输入同字节」失效（ADR-0004：值对象只带纯函数行为）。
//  2. 不重复定义时间语义：TimeSemantics 住 `domain/market`，这里只承载点位本身。
//  3. 不放聚合器 / 端口：Aggregator 今天无实现、无消费方，随 X1b 定形。
package dataprotocol

import (
	"encoding/json"
	"time"
)

// DataType* 是 UnifiedDataPoint.DataType 的取值集，决定写哪张表。
//
// 字段今天仍是裸 string（具名化会碰 wire，等 X1b 的 schema 冻结再评估）；
// 这组常量是取值集的共享定义。迁移期 core / data 的本地副本在各自单收编前
// 仍并存，取值必须与这里逐字一致。
const (
	DataTypeOHLCDaily   = "ohlcv_daily"
	DataTypeOHLCMinute  = "ohlcv_minute"
	DataTypeRealtime    = "realtime_quote"
	DataTypeCapitalFlow = "capital_flow"
	DataTypeFundamental = "fundamentals"
	DataTypeSectors     = "sectors"
	DataTypeStockSector = "stock_sector"
	DataTypeTopList     = "top_list"
	DataTypeLimitUpPool = "limit_up_pool"
	DataTypeAnnounce    = "announcements"
	DataTypeNews        = "news"
	DataTypeHotSearch   = "hot_search"
	DataTypeGlobalOHLCV = "global_ohlcv"
	DataTypeMacro       = "macro_series"
	DataTypeDailyBasic  = "daily_basic"
)

// UnifiedDataPoint 是数据在 ETL / 写入路径上的共享形状。
//
// 字段顺序与 `tellurion-data:internal/ingest/point.go` 的镜像一致
// （VisibleAt 在末位）：使用位置化字面量的调用点按这个顺序书写。
type UnifiedDataPoint struct {
	// Symbol 是项目约定的规范代码（"600519.SH" / "AAPL"）。
	// 交易所后缀归一在 connector，不在这里做。
	Symbol string

	// TradeTime 是数据点归属时刻：日频 = 交易日零点；实时 = 快照时刻；
	// 财报 = 报告期截止日（PIT 可见时刻另走 VisibleAt）。
	TradeTime time.Time

	// Source 是 connector 名（"tushare" / "mootdx"），落 `source` 列。
	// 同一张表里多源口径靠它区分。
	Source string

	// DataType 决定写哪张表，取值见 DataType* 常量。
	DataType string

	// Data 是该类型的类型化字段（JSONB 形状）。必含键由 connector 声明与
	// normalizer 决定，不在类型上表达。
	Data map[string]interface{}

	// IngestTime 是落库时刻，**必须由调用方注入**（本包不读时钟）。
	IngestTime time.Time

	// VisibleAt 是按声明的可见锚点导出的 PIT 可见时刻，由写侧闸门填写。
	// 零值 = 尚未填写，**不得**把它当成 TradeTime 的替代。
	VisibleAt time.Time
}

// ToJSON 序列化 Data（JSONB 列 / 审计留痕）。
func (p UnifiedDataPoint) ToJSON() ([]byte, error) {
	return json.Marshal(p.Data)
}

// DeduplicateKey 返回 ETL 去重用的键：Symbol | DataType | TradeTime（UTC，RFC3339Nano）。
//
// Source 刻意不在键里：取数的 fallback 链可能从多个源加载同一条数据，
// 保留哪一条由调用方按源优先级决定（本函数不裁决）。
func (p UnifiedDataPoint) DeduplicateKey() string {
	return p.Symbol + "|" + p.DataType + "|" + p.TradeTime.UTC().Format(time.RFC3339Nano)
}
