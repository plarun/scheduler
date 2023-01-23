package task

type Profile struct {
	key   Field
	value string
	empty bool
}

func NewProfile(val string, empty bool) *Profile {
	return &Profile{
		key:   FIELD_PROFILE,
		value: val,
		empty: empty,
	}
}

func (p *Profile) UnimplementedTaskField() {}

func (p *Profile) Key() Field {
	return p.key
}

func (p *Profile) Name() string {
	return string(p.key)
}

func (p *Profile) Empty() bool {
	return p.empty
}

func (p *Profile) Value() string {
	return p.value
}
