# 本仓边界与依赖方向 — tellurion-contracts

## 定位

跨仓唯一缝隙：**port 接口 · DTO/值对象（含 domain）· 四个引擎契约 · 哨兵错误 · abi 常量**。所有跨边界的东西只能住这里。

状态：**S1 落代码（QT-05）／S2 打 tag（QT-07）——当前为占位仓，尚无实现**

## 依赖方向（编译器强制）

```
github.com/ruoxizhnya/tellurion-kit
        ↑
github.com/ruoxizhnya/tellurion-contracts
        ↑                    ↑
   tellurion (core)   tellurion-agent / 各服务
```

- 箭头 = import 方向。**不允许反向**，也不允许成环 —— 反向 import 应当**编译失败**，这是判据不是口号。
- 跨仓只认版本号：`github.com/ruoxizhnya/tellurion-contracts` 以需要的版本出现在下游 `go.mod`。

## 边界表

| | 内容 |
|---|---|
| 允许 | 被 core / agent / 各服务 `require`；只放**被 ≥2 个仓使用**的跨边界类型与接口 |
| 禁止 | 任何实现、业务规则、DB 访问；`Config` / `Options` / `Impl` 命名（契约泄漏实现细节的形状）；只有单个仓使用的类型（那条属于泄漏，应留在使用方） |

## 关键约束

1. **单仓可构建**：`go build ./...` / `go vet ./...` / `go test ./...` 不依赖任何兄弟目录。
2. **无 `replace`**：指向仓外的 `replace` 一律不接受（CI 会查）。
3. **无 `go.work` 入库**：本地联调可以建，提交即错。
4. **契约只增不破**：0.x 期间允许 breaking，但必须 bump MINOR + 写 CHANGELOG（版本规则见 QT-19）。

## 未决 / 待定

- 包名策略：`engine/` 内保持 `package contracts`、`domain/` 内 `package domain`（QT-05 ⑧）。
- 版本规则表与兼容矩阵：QT-19（未落地前不要凭印象发版）。
