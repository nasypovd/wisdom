package pow

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"wisdom/pkg/domain"
)

type PoW struct {
	Difficulty int
	Rand       *rand.Rand
}

func NewPoW(difficulty int) *PoW {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return &PoW{Difficulty: difficulty, Rand: r}
}

func (pow *PoW) Generate() domain.Challenge {
	value := pow.Rand.Intn(1000000)
	return domain.Challenge{Value: strconv.Itoa(value), Difficulty: pow.Difficulty}
}

func (pow *PoW) Verify(challenge domain.Challenge, solution domain.Solution) bool {
	hash := sha256.Sum256([]byte(challenge.Value + string(solution)))
	hashStr := fmt.Sprintf("%x", hash)
	prefix := strings.Repeat("0", pow.Difficulty)
	return strings.HasPrefix(hashStr, prefix)
}

type Solver struct{}

func NewSolver() *Solver {
	return &Solver{}
}

func (s *Solver) Solve(challenge domain.Challenge) domain.Solution {
	for nonce := 0; ; nonce++ {
		hash := sha256.Sum256([]byte(challenge.Value + strconv.Itoa(nonce)))
		hashStr := fmt.Sprintf("%x", hash)
		prefix := strings.Repeat("0", challenge.Difficulty)

		if strings.HasPrefix(hashStr, prefix) {
			return domain.Solution(strconv.Itoa(nonce))
		}
	}
}
