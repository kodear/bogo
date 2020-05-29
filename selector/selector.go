package selector

type Selector struct {
	body  []byte
}

func (cls *Selector)String()string{
	return string(cls.body)
}


