package condition

type clause struct {
	Parent     *wrapper
	Status     string
	JobName    string
	Operator   Operator
	Conditions []Expression
}

// newClause creates a new job condition which represents
// status, job_name and boolean join operator '&' or '|'.
func newClause(status, job, op string) *clause {
	return &clause{
		Status:     status,
		JobName:    job,
		Operator:   stringToCondOperator(op),
		Conditions: make([]Expression, 0),
	}
}

func (c *clause) build() string {
	return c.Status + "(" + c.JobName + ")" + condOperatorToString(c.Operator)
}

func (c *clause) isWrapper() bool {
	return false
}

func (c *clause) addChild(_ Expression) {
}

func (c *clause) setParent(parent *wrapper) {
	c.Parent = parent
}

func (c *clause) getParent() *wrapper {
	return c.Parent
}

func (c *clause) String() string {
	return c.build()
}
