package todo

type TodoID string

func (id TodoID) String() string { return string(id) }

func (id TodoID) Valid() bool { return len(id) > 0 }
