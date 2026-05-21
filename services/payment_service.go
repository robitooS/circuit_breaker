package services

import "fmt"

func Pay(req string) (any, error) {
    fmt.Println("Teste pagamento ...")
    return "Pagamento Processado com Sucesso!", nil 
}