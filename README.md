# Circuit Breaker em Go

Este projeto é uma implementação educacional e funcional do padrão de design **Circuit Breaker** utilizando a linguagem Go. O objetivo é proteger sistemas distribuídos contra falhas em cascata, interrompendo chamadas a serviços que estão instáveis e permitindo uma recuperação controlada.

## Objetivo

O Circuit Breaker atua como um interruptor de segurança para chamadas de rede. Quando um serviço remoto começa a falhar repetidamente, o "circuito" se abre, impedindo que novas requisições sobrecarreguem o serviço problemático ou consumam recursos desnecessários do sistema chamador.

## Como Funciona

A implementação segue a máquina de estados clássica:

1. **Closed (Fechado):** O fluxo normal. As requisições passam livremente. Erros são contabilizados. Se o número de erros atingir o limiar (`failureThreshold`), o circuito muda para **Open**.
2. **Open (Aberto):** O circuito bloqueia todas as requisições imediatamente, retornando um erro de "circuito aberto". Após um tempo de espera (10 segundos nesta implementação), ele muda para **Half-Open**.
3. **Half-Open (Meio-Aberto):** Uma requisição de "sonda" é permitida para testar a saúde do serviço.
   * Se a sonda **falhar**, o circuito volta para **Open**.
   * Se a sonda tiver **sucesso**, o circuito volta para **Closed** e o contador de falhas é resetado.

## Estrutura do Projeto

* `main.go`: Ponto de entrada que simula um ciclo de vida completo: falhas iniciais, abertura do circuito, tempo de espera e recuperação.
* `security/circuit_breaker.go`: Core do sistema. Contém a lógica de gerenciamento de estados e a função `Execute` que envolve as chamadas.
* `services/payment_service.go`: Mock de um serviço de pagamento que simula falhas temporárias para demonstrar o comportamento do breaker.
* `services/service_teste.go`: Camada de serviço que integra a lógica de negócio com o Circuit Breaker.

## Como Executar

Certifique-se de ter o Go instalado em sua máquina.

1. Clone o repositório ou navegue até a pasta do projeto.
2. Execute o comando:
   ```bash
   go run main.go
   ```

Você verá no console as diferentes fases do circuito, desde as falhas iniciais até a recuperação automática.

## Diferenciais desta Implementação

* **Thread-Safety:** Uso de `sync.Mutex` para garantir que o estado do circuito seja gerenciado corretamente em ambientes concorrentes.
* **Encapsulamento:** A lógica do breaker está separada da lógica de negócio, permitindo fácil reutilização em outros serviços.
* **Simulação Realista:** O `main.go` demonstra não apenas o erro, mas o processo de auto-recuperação após o timeout.

## Sugestões de Melhorias (Roadmap)

Para tornar este projeto pronto para produção, poderiam ser adicionados:

- [ ] **Configuração Dinâmica:** Permitir configurar o limiar de falhas e o tempo de timeout via variáveis de ambiente ou arquivo de config.
- [ ] **Métricas e Logs:** Integração com Prometheus/Grafana para monitorar o estado do circuito em tempo real.
- [ ] **Testes de Unidade:** Cobertura de testes para validar cada transição de estado.
- [ ] **Suporte a Context:** Adicionar suporte a `context.Context` para lidar com cancelamentos e timeouts de requisição.
- [ ] **Exponential Backoff:** Implementar um tempo de espera que aumenta progressivamente se o serviço continuar falhando.

---

Desenvolvido para fins de estudo sobre Resiliência de Sistemas.
