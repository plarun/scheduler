package condition

type wrapper struct {
	Parent     *wrapper
	Operator   Operator
	Conditions []Expression
}

// newwrapper creates a new wrapper to enclose or group
// multiple condition. Status and JobName are empty.
func newWrapper() *wrapper {
	return &wrapper{
		Operator:   OperatorEmpty,
		Conditions: make([]Expression, 0),
	}
}

func (w *wrapper) build() string {
	var str string
	str += "("
	for _, child := range w.Conditions {
		str += child.build()
	}
	str += ")" + condOperatorToString(w.Operator)
	return str
}

func (w *wrapper) String() string {
	str := w.build()
	return str[1 : len(str)-1]
}

func (w *wrapper) isWrapper() bool {
	return true
}

func (w *wrapper) addChild(child Expression) {
	w.Conditions = append(w.Conditions, child)
	child.setParent(w)
}

func (w *wrapper) setParent(parent *wrapper) {
	w.Parent = parent
}

func (w *wrapper) getParent() *wrapper {
	return w.Parent
}
