package configTypes

type ConfigPath string

func (c ConfigPath) String() string {
	return string(c)
}
