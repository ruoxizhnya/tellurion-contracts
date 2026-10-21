---
status: active
last-verified: 2026-10-20
verified-by: 功能边界固化（三份契约：db / api / dataflow）
---

# Dataflow — tellurion-contracts

> 目的：把「数据从哪来、在谁手里变形、到哪去」写死。**任何跨仓数据交换都必须能在本图上找到一条边**；找不到 = 这段代码在偷偷建第二条路。

## 1. 主链路

```
core ──┐                        ┌── tellurion-data  (实现 DataPort)
agent ─┼── require/版本 ──► contracts ◄── require/版本 ──┤   tellurion-broker (实现 BrokerGateway)
服务 ──┘      (types/ports)         └── tellurion-strategy / risk / analytics / orchestration
                                     （都只依赖形状，不依赖彼此实现）
```

## 2. 关键时序

1. 下游 `go.mod` 出现 `github.com/ruoxizhnya/tellurion-contracts vX.Y.Z`
2. 下游代码 `import .../domain`、`.../contracts/engine`
3. 编译期完成校验：形状不符 → 编不过（这就是"契约"的全部强制力）
4. 运行期由实现方保证语义（契约测试在 QT-09 ③）

## 3. 失败模式（每一处 fail-closed 都要能指名）

| 失败点 | 症状 | 兜底 |
|---|---|---|
| 类型漂移（两侧各写一份理解） | 编译过、行为不一致 | 契约测试（QT-09 ③）：同一组用例喂两个实现 |
| 契约里塞实现 | 改一处牵动全部下游 | Code review + AST 扫描（禁 `Impl`/`Config`/`Options`） |
| 单仓专用类型混入 | contracts 膨胀成 common | 「≥2 仓使用」判据 + core-only 清单回收 |

## 4. 不变量（可被测试钉住的）

1. 本仓 `go build ./...` / `go test ./...` 在没有其他仓时也绿（R1/R2）。
2. 反向 import（contracts → core/服务）编译失败。
3. `domain` 之外的包不得 import `domain` 的内部实现细节（只允许用其公开类型）。

## 5. 与契约文档的关系

- 表结构 → `docs/db-schema.md`
- 接口形状 → `docs/api-schema.md`
- 本文件只讲**流向与顺序**，不重复字段定义。
