package services

import "fmt"

var callCount int

func Pay(req string) (any, error) {
    callCount++
    
    // Simula falha nas primeiras 12 chamadas
    if callCount <= 10 {
        return nil, fmt.Errorf("payment service unavailable (call %d)", callCount)
    }

    fmt.Println("Processando pagamento...")
    return "Pagamento Processado com Sucesso!", nil
}