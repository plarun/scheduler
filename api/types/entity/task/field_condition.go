package task

import (
	"github.com/plarun/scheduler/api/types/condition"
)

type Condition struct {
	key   Field
	value string
	empty bool
}

func NewCondition(val string, empty bool) *Condition {
	return &Condition{
		key:   FIELD_CONDITION,
		value: val,
		empty: empty,
	}
}

func (c *Condition) UnimplementedTaskField() {}

func (c *Condition) Key() Field {
	return c.key
}

func (c *Condition) Name() string {
	return string(c.key)
}

func (c *Condition) Empty() bool {
	return c.empty
}

func (c *Condition) DistinctTasks() []string {
	return condition.GetDistinctTasks(c.value)
}

func (c *Condition) Value() string {
	return c.value
}
