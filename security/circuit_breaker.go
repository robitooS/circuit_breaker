package security

import (
	"fmt"
	"sync"
	"time"
)

type State string

const (
	Open     State = "open"
	Closed   State = "closed"
	HalfOpen State = "half-open"
)

type CircuitBreaker struct {
	mu               sync.Mutex
	state            State
	failureCount     int
	failureThreshold int // Limiar de erro que deve chegar no máximo, por definição 10
	lastStageChange  time.Time
	halfOpenTested   bool
}

func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{
		state:            "closed",
		failureCount:     0,
		failureThreshold: 10,
		lastStageChange:  time.Now(),
	}
}

func (cb *CircuitBreaker) doPreRequest() error {
	// Se estiver aberto o circuito, verifica se o tempo ultrapassou 10 segundos, assim ele coloca o circuito como meio aberto
	if cb.state == Open {
		if time.Since(cb.lastStageChange) > 10*time.Second {
			cb.state = HalfOpen
			cb.halfOpenTested = false
		} else {
			return fmt.Errorf("circuit is open, wait before retrying")
		}
	}

	if cb.state == HalfOpen {
		if cb.halfOpenTested {
			return fmt.Errorf("already tested, waiting for result")
		}
		cb.halfOpenTested = true
	}

	return nil
}

func (cb *CircuitBreaker) doAfterRequest(errReq error) error {
	// Se a resposta da requisição vier um erro, devemos tratar no @doAfterRequest para incrementar os erros e verificar o fluxo a se seguir
	if cb.state == HalfOpen {
		if errReq != nil {
			cb.state = Open
			cb.lastStageChange = time.Now()
		} else {
			cb.state = Closed
			cb.failureCount = 0
		}
		cb.halfOpenTested = false
		return nil
	}

	if errReq != nil {
		cb.failureCount++
	}

	if cb.failureCount >= cb.failureThreshold {
		cb.state = Open
		cb.failureCount = 0
		cb.lastStageChange = time.Now()
		return nil
	}

	if errReq == nil {
		cb.failureCount = 0
		cb.state = Closed
	}
	return nil
}

func (cb *CircuitBreaker) Execute(req func() (any, error)) (any, error) {
	cb.mu.Lock()
	err := cb.doPreRequest()
	cb.mu.Unlock()

	if err != nil {
		return nil, err
	}

	res, errReq := req()

	cb.mu.Lock()
	_ = cb.doAfterRequest(errReq)
	cb.mu.Unlock()
	
	if errReq != nil {
		return nil, errReq
	}

	return res, nil
}
