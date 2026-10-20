# tellurion-contracts

> module: `github.com/ruoxizhnya/tellurion-contracts` ｜ 状态：**待建（S2 / QT-07）**——当前为占位仓，**尚无实现**。迁入动作与验收见 `tellurion:docs/TASKS.md` 的 QT-05（S1 仓内抽取）与 QT-07（S2 拆仓）。

跨仓唯一缝隙：**port 接口 · DTO/值对象（含 domain）· 四个引擎契约 · 哨兵错误 · abi 常量 · 数据 schema**。

本仓由 core 的 `pkg/domain`（17 文件）+ `pkg/backtest/contracts`（6 文件，含 RunSpec / SignalSpec / ResultSpec / EngineRunner）+ core 顶层 `contracts/*.schema.json` 迁入。

## 边界

**可以**
- 被 core / agent / 各服务 `require`
- 只放跨边界的东西（被 ≥2 个仓使用）

**不可以**
- 反向依赖任何实现仓（core / agent / data / risk …）
- 出现 `Config` / `Options` / `Impl` 命名或 DB 访问
- 放只有单一仓使用的类型（那条属于泄漏）

## 依赖方向

```
github.com/ruoxizhnya/tellurion-contracts
```

- 契约与工具一律以 **tag** 依赖（`github.com/ruoxizhnya/tellurion-contracts` / `-kit`），
  **禁止 `replace` 指向仓外路径**（单仓不可构建）。
- `go.work` 一律不入库。

## 施工入口

本仓所有任务登记在 **`tellurion` 仓的 `docs/TASKS.md`**（单一台账，避免双写）。
在台账出现本仓的单号之前，**不要在此写实现**。

## 权威文档

| 主题 | 位置 |
|---|---|
| 仓清单与命名 | `tellurion:docs/design/qt-decoupling/repo-inventory.md` |
| 分仓拓扑与阶梯 | `tellurion:docs/design/qt-decoupling/repo-topology-and-naming.md` |
| P0 修复基线 | `tellurion:docs/design/qt-decoupling/p0-round2-review.md` |
