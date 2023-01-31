package condition

type Clause struct {
	Parent     *Wrapper
	Status     string
	TaskName   string
	Operator   Operator
	Conditions []Expression
	Result     bool
}

// newClause creates a new task condition which represents
// status, job_name and boolean join operator '&' or '|'.
func newClause(status, task, op string) *Clause {
	return &Clause{
		Status:     status,
		TaskName:   task,
		Operator:   stringToCondOperator(op),
		Conditions: make([]Expression, 0),
	}
}

func (c *Clause) build() string {
	return c.Status + "(" + c.TaskName + ")" + condOperatorToString(c.Operator)
}

func (c *Clause) IsWrapper() bool {
	return false
}

func (c *Clause) addChild(_ Expression) {
}

func (c *Clause) setParent(parent *Wrapper) {
	c.Parent = parent
}

func (c *Clause) getParent() *Wrapper {
	return c.Parent
}

func (c *Clause) String() string {
	return c.build()
}
