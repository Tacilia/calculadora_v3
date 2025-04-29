package main

import (
	"database/sql" //pacote para conectar e trabalhar com bancos de dados.
	"fmt"          //imprime menssagens no terminal

	_ "github.com/mattn/go-sqlite3" // importa o driver do banco SQLite, o _ é porque só queremos registrar o driver, sem usar diretamente.
	//Tambbém instalei a dependencia via terminal: go get github.com/mattn/go-sqlite3
	//permite que a linguagem Go consiga se conectar e fazer consultas em um banco de dados SQLite.
)

func main() {
	// Abre (ou cria) o banco de dados SQLite com nome "calculadora.db"
	// A função sql.Open retorna dois valores: um ponteiro para o banco (db) e um erro (err)
	db, err := sql.Open("sqlite3", "calculadora.db")
	if err != nil {
		// Se houver erro ao abrir o banco, o programa é encerrado com panic
		panic(err)
	}
	// Garante que o banco de dados será fechado ao final da função main
	defer db.Close()

	// Comando SQL para criar a tabela "operacoes" caso ela ainda não exista
	sqlTabela := `CREATE TABLE IF NOT EXISTS operacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID único, gerado automaticamente
		operando1 REAL,                        -- Primeiro número da operação (float)
		operando2 REAL,                        -- Segundo número da operação (float)
		operacao TEXT,                         -- Tipo da operação: "+", "-", "*", "/"
		resultado REAL                         -- Resultado do cálculo
	);`

	// Executa o comando SQL para criar a tabela
	resultadoTabela, err := db.Exec(sqlTabela)
	if err != nil {
		// Se ocorrer erro, o programa para imediatamente
		panic(err)
	}
	// Mostra no terminal que a tabela foi criada (ou já existia)
	fmt.Println("Tabela criada ou já existia:", resultadoTabela)

	// Insere uma operação na tabela: 2 + 3 = 5
	// Os "?" são placeholders que evitam SQL injection
	resultadoInsert, err := db.Exec(
		"INSERT INTO operacoes (operando1, operando2, operacao, resultado) VALUES (?, ?, ?, ?)",
		2, 3, "+", 5,
	)
	if err != nil {
		// Se houver erro ao inserir, o programa para
		panic(err)
	}

	// Verifica quantas linhas foram modificadas (inseridas)
	linhasModificadas, err := resultadoInsert.RowsAffected()
	if err != nil {
		// Se não conseguir contar as linhas, mostra erro e sai da função
		fmt.Println("Erro ao obter linhas modificadas:", err)
		return
	}
	// Mostra o número de linhas modificadass (deve ser 1)
	fmt.Println("Inserção feita. Linhas modificadas:", linhasModificadas)

	// Consulta todos os dados da tabela "operacoes"
	linhas, err := db.Query("SELECT id, operando1, operando2, operacao, resultado FROM operacoes")
	if err != nil {
		// Se der erro na consulta, o programa para
		panic(err)
	}
	// Garante que os resultados da consulta serão fechados no final
	defer linhas.Close()

	// Cabeçalho para os resultados
	fmt.Println("\n📋 Resultados:")

	// Percorre cada linha retornada pelo SELECT
	for linhas.Next() {
		var id int          // Armazena o ID da operação
		var a, b, r float64 // Armazena os operandos e o resultado
		var op string       // Armazena o tipo da operação

		// Pega os valores da linha atual e coloca nas variáveis acima
		linhas.Scan(&id, &a, &b, &op, &r)

		// Imprime a operação formatada, exemplo: "1 -> 2.0 + 3.0 = 5.0"
		fmt.Printf("%d -> %.1f %s %.1f = %.1f\n", id, a, op, b, r)
	}
}
