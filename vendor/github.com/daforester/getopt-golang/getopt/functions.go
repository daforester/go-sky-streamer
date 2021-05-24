package getopt

type addOperandFunc func(string) error
type getValueFunc func() string
type handlerFunc func(opt *GetOpt) error
type setOptionFunc func(string, getValueFunc) error
type setCommandFunc func(CommandInterface) error
type validationFunc func(...string) bool
