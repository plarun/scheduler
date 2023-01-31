package condition

type Wrapper struct {
	Parent     *Wrapper
	Operator   Operator
	Conditions []Expression
	Result     bool
}

// newwrapper creates a new wrapper to enclose or group
// multiple condition. Status and JobName are empty.
func newWrapper() *Wrapper {
	return &Wrapper{
		Operator:   OperatorEmpty,
		Conditions: make([]Expression, 0),
	}
}

func (w *Wrapper) build() string {
	var str string
	str += "("
	for _, child := range w.Conditions {
		str += child.build()
	}
	str += ")" + condOperatorToString(w.Operator)
	return str
}

func (w *Wrapper) String() string {
	str := w.build()
	return str[1 : len(str)-1]
}

func (w *Wrapper) IsWrapper() bool {
	return true
}

func (w *Wrapper) addChild(child Expression) {
	w.Conditions = append(w.Conditions, child)
	child.setParent(w)
}

func (w *Wrapper) setParent(parent *Wrapper) {
	w.Parent = parent
}

func (w *Wrapper) getParent() *Wrapper {
	return w.Parent
}
