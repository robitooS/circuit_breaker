package services

import (
	"log"

	"github.com/robitooS/circuit_breaker/security"
)

func CallPaymentService(cb *security.CircuitBreaker, req string) {
	r, err := cb.Execute(func() (any, error) {
		return Pay(req)
	})
	if err != nil {
		log.Printf("got error: %s", err.Error())
		return
	}

	resp, ok := r.(string)
	if !ok {
		log.Println("type assertion failed")
		return
	}

	log.Printf("got response: %s", resp)

}