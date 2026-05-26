package main

import (
    "fmt"
    "time"

    "github.com/robitooS/circuit_breaker/security"
    "github.com/robitooS/circuit_breaker/services"
)

func main() {
    fmt.Println("=== Iniciando Circuit Breaker ===")
    cb := security.NewCircuitBreaker()

    // Fase 1: dispara 10 falhas → circuito abre
    fmt.Println("\n--- Fase 1: forçando falhas ---")
    for i := 0; i < 12; i++ {
        services.CallPaymentService(cb, "req")
    }

    // Fase 2: circuito aberto → todas bloqueadas
    fmt.Println("\n--- Fase 2: circuito aberto ---")
    services.CallPaymentService(cb, "req")
    services.CallPaymentService(cb, "req")

    // Fase 3: espera 10s → entra em HalfOpen
    fmt.Println("\n--- Fase 3: aguardando 10s ---")
    time.Sleep(11 * time.Second)

    // Fase 4: primeira req passa (sonda) → serviço ainda falha (calls 13-14 <= 12? não, já passou)
    // callCount já está em 12, próxima vai ser 13 → sucesso
    fmt.Println("\n--- Fase 4: testando HalfOpen ---")
    services.CallPaymentService(cb, "req") // sonda → sucesso → fecha circuito

    // Fase 5: circuito fechado, tudo funcionando
    fmt.Println("\n--- Fase 5: circuito fechado ---")
    services.CallPaymentService(cb, "req")
    services.CallPaymentService(cb, "req")
}