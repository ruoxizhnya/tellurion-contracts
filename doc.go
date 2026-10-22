// Package contracts 是 tellurion-contracts 的根包（仅文档，零导出符号）。
//
// 库代码住子包：
//   - fees/          费率模型的形状
//   - domain/market/ 行情值对象 + 频率 / 时间语义词表
//   - domain/        信号 / 订单 / 组合 / 公司行为 / 因子名册
//   - engine/        四契约 / 内核路由 / canonical 身份 / 哨兵 / abi（包名仍是 contracts）
//
// 保留这个文件是为了让 `go build ./...` 始终有一个根包可构建，并让工具能
// `go list ./...`；不要往根包塞代码。
package contracts
