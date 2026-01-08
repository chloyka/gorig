package configTypes

type ConfigSaver interface {
	SetSaveChan(chan struct{})
	Save()
}

type configSaver struct {
	saveChan chan struct{}
}

func (cs *configSaver) SetSaveChan(ch chan struct{}) {
	cs.saveChan = ch
}

func (cs *configSaver) Save() {
	if cs.saveChan != nil {
		select {
		case cs.saveChan <- struct{}{}:
		default:

		}
	}
}
