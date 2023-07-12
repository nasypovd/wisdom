package quote

import (
	"math/rand"
	"time"
	"wisdom/pkg/domain"
)

type Quotes struct {
	data []string
	Rand *rand.Rand
}

func New(quotes []string) *Quotes {
	if len(quotes) == 0 {
		panic("quotes must not be empty") // panic early at initialization time
	}

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return &Quotes{data: quotes, Rand: r}
}

func (r *Quotes) Get() domain.Quote {
	randomIndex := r.Rand.Intn(len(r.data))
	return domain.Quote(r.data[randomIndex])
}
