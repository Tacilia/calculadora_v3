package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Estrutura do Tipo OperacaoRequest para receber os dados da operação
type OperacaoRequest struct {
	Operando1 float64 `json:"operando1"`
	Operando2 float64 `json:"operando2"`
	Operacao  string  `json:"operacao"`
}

// Estrutura do Tipo ResultadoResponse para enviar o resultado
type ResultadoResponse struct {
	Resultado float64 `json:"resultado"`
}

// E Estrutura do Tipo HistoricoItem para guardar uma operação completa com resultado no histórico.
type HistoricoItem struct {
	OperacaoRequest         // Ela inclui todos os campos de OperacaoRequest.
	Resultado       float64 `json:"resultado"` // adiciona o campo Resultado
}

// Slice global para armazenar o histórico
var historico []HistoricoItem

func main() {
	http.HandleFunc("/soma", somaHandler)
	http.HandleFunc("/subtracao", subtracaoHandler)
	http.HandleFunc("/multiplicacao", multiplicacaoHandler)
	http.HandleFunc("/divisao", divisaoHandler)
	http.HandleFunc("/historico", historicoHandler)

	fmt.Println("Servidor rodando na porta 8080...")
	http.ListenAndServe(":8080", nil)
}

// -------------------- HANDLERS --------------------

// Função que responde a um pedido (requisição) que chega no seu servidor.
func somaHandler(w http.ResponseWriter, r *http.Request) { //Aqui dentro dos parênteses (w, r), estamos recebendo dois parâmetros.
	// w --> É quem escreve a resposta que o servidor vai mandar de volta para quem pediu (o navegador, o Insomnia, o Postman, etc). Serve para enviar respostas (texto, JSON, erro, etc).
	// r --> É quem contém todos os dados da requisição que o servidor recebeu (como o corpo da requisição, os parâmetros, o tipo de método GET, POST, etc).
	//Ou seja, a funcao somaHandler é o intermediador que recebe o pedido,  r é serviro que lê o pedido que chegou ("quero somar 2 + 2"). E w é a resposta que O servidor responde ("o resultado é 4").
	//1. O cliente manda uma requisição (tipo /soma).
	//2. O Go vê que /soma está ligado ao somaHandler.
	// // 3. O Go executa o somaHandler, que pega o pedido (r) e manda uma resposta (w).

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	} // Checando se a requiscao é metodo POST, se nao for envia a menssagem metodo nao permitido.

	var req OperacaoRequest //req é uma variável do tipo OperacaoRequest. Serve para ler informações que chegaram na requisição.
	json.NewDecoder(r.Body).Decode(&req)
	//json.NewDecoder(r.Body) --> Cria um "leitor" para os dados que chegaram na requisição HTTP.
	//.Decode(&req) --> Lê os dados JSON e preenche a variável 'req' usando o endereço http dela.

	resultado := req.Operando1 + req.Operando2 //Crie uma variável chamada resultado e guarde dentro dela o valor do primeiro número e segundo número que chegou como resposta na requisição.
	// := significa "criar e guardar" (chamamos isso de declaração e atribuição).

	adicionarAoHistorico(req, resultado) //Chama a função adicionarAoHistorico para registrar essa operação.

	json.NewEncoder(w).Encode(ResultadoResponse{Resultado: resultado}) // Serve para mandar uma resposta no formato JSON para quem chamou a sua API.
}

func subtracaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use o método POST", http.StatusMethodNotAllowed)
		return
	}

	var req OperacaoRequest
	json.NewDecoder(r.Body).Decode(&req)

	resultado := req.Operando1 - req.Operando2

	adicionarAoHistorico(req, resultado)

	json.NewEncoder(w).Encode(ResultadoResponse{Resultado: resultado})
}

func multiplicacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use o método POST", http.StatusMethodNotAllowed)
		return
	}

	var req OperacaoRequest
	json.NewDecoder(r.Body).Decode(&req)

	resultado := req.Operando1 * req.Operando2

	adicionarAoHistorico(req, resultado)

	json.NewEncoder(w).Encode(ResultadoResponse{Resultado: resultado})
}

func divisaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use o método POST", http.StatusMethodNotAllowed)
		return
	}

	var req OperacaoRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.Operando2 == 0 {
		http.Error(w, "Erro: nenhum número pode ser dividido por zero", http.StatusBadRequest)
		return
	}

	resultado := req.Operando1 / req.Operando2

	adicionarAoHistorico(req, resultado)

	json.NewEncoder(w).Encode(ResultadoResponse{Resultado: resultado})
}

func historicoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Use o método GET", http.StatusMethodNotAllowed)
		return
	}

	json.NewEncoder(w).Encode(historico)
}

// -------------------- FUNÇÃO AUXILIAR --------------------

func adicionarAoHistorico(req OperacaoRequest, resultado float64) {
	item := HistoricoItem{
		OperacaoRequest: req,
		Resultado:       resultado,
	}
	historico = append(historico, item)
}
