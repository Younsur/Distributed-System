package main

type Bracha struct {
	Parties []*Party
}

func NewBracha() *Bracha {
	bracha := &Bracha{}
	for i := 0; i < N; i++ {
		node := &Party{
			IsHonesty: true,
			Id:        i,
			Echo:      make(map[string]int),
			Ready:     make(map[string]int),
			Bracha:    bracha,
			OutputCh:  make(chan string),
			IsLocked:  false,
			Content:   "",
		}
		bracha.Parties = append(bracha.Parties, node)
	}
	return bracha
}
