// Package contracts 的根包夹具：验证「根包导出符号也能登记」（小节标题写 `.`）。
package contracts

// RootValue 挂在根包上：不认 `.` 标题的护栏会让它无路可登记，当场永久红。
var RootValue = "root"
