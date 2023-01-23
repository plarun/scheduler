package task

type Command struct {
	key   Field
	value string
	empty bool
}

func NewCommand(val string, empty bool) *Command {
	return &Command{
		key:   FIELD_COMMAND,
		value: val,
		empty: empty,
	}
}

func (c *Command) UnimplementedTaskField() {}

func (c *Command) Key() Field {
	return c.key
}

func (c *Command) Name() string {
	return string(c.key)
}

func (c *Command) Empty() bool {
	return c.empty
}

func (c *Command) Value() string {
	return c.value
}
