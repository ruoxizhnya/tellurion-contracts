package contracts

// EngineRunner 是端口：方法集就是它的形状，逐个方法都要登记。
type EngineRunner interface {
	RunBacktest(name string) error
	Close()
}

// Status 是结果状态。
type Status string

const (
	StatusCompleted Status = "completed"
	statusInternal  Status = "internal"
)

// ErrFault 是哨兵。
var ErrFault = Status("fault")

// New 返回一个执行器。
func New() EngineRunner { return nil }

// Runner 是导出类型：导出方法进导出面，未导出的不进。
type Runner struct{}

func (r *Runner) Run() {}

func (r *Runner) run() {}

type internal struct{}

func (i internal) DoIt() {}
