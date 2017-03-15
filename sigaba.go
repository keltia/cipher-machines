package machine

import (

)

type Sigaba struct {
}

func (m *Sigaba) Step() {

}

func (m *Sigaba) Index() {

}

func (m *Sigaba) Encrypt(text string) (cipher string) {

	return
}

func (m *Sigaba) Decrypt(text string) (clear string) {

	return
}

func NewSigaba() (m *Sigaba, err error) {
	return &Sigaba{}, nil
}


