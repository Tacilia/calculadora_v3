Anotacoes:

func Handler(w http.ResponseWriter, r *http.Request) --> Aqui dentro dos parênteses (w, r), estamos recebendo dois parâmetros. w --> É quem escreve a resposta que o servidor vai mandar de volta para quem pediu (o navegador, o Insomnia, o Postman, etc). Serve para enviar respostas (texto, JSON, erro, etc). r --> É quem contém todos os dados da requisição que o servidor recebeu (como o corpo da requisição, os parâmetros, o tipo de método GET, POST, etc). Ou seja, a funcao somaHandler é o intermediador que recebe o pedido, r é serviro que lê o pedido que chegou ("quero somar 2 + 2"). E w é a resposta que O servidor responde ("o resultado é 4").

O cliente manda uma requisição (tipo /soma).
O Go vê que /soma está ligado ao somaHandler.
O Go executa o somaHandler, que pega o pedido (r) e manda uma resposta (w).
Por que precisamos instalar dependencia com go get github.com/mattn/go-sqlite3? O Go sozinho não sabe como conversar com o SQLite. Então precisamos instalar essa biblioteca externa que ensina o Go a:

abrir um banco .db,
criar tabelas,
inserir dados,
buscar informações, etc.
Essa biblioteca faz a "ponte" entre o código em Go e o banco SQLite.

go get: baixa e instala uma dependência externa no seu projeto.

github.com/mattn/go-sqlite3: o endereço do pacote do driver SQLite para Go.

Depois de rodar isso:

O Go baixa o código dessa biblioteca.
E adiciona no go.mod o registro que estamos usando está dependencia.

Utilizei o (_ "modernc.org/sqlite") para nao precisar baixar o compilador C e assim criar/registrar a tabela do banco de dados. 
a biblioteca github.com/mattn/go-sqlite3 é um "wrapper" em torno da biblioteca SQLite original, que é escrita em linguagem C. Para que o código Go possa interagir com essa biblioteca C, ele usa um recurso chamado CGO. O CGO permite que você escreva código Go que chama funções escritas em C. Para que isso funcione, você precisa ter um compilador C instalado no seu sistema para que o CGO possa compilar o código C da biblioteca SQLite e criar as "pontes" necessárias para que o seu código Go possa usá-lo.

É por isso que, quando o CGO não está habilitado ou um compilador C não é encontrado, a biblioteca github.com/mattn/go-sqlite3 não consegue funcionar corretamente e vê aquela mensagem de erro.

Go Modules (go mod)

O Go Modules é o sistema oficial de gerenciamento de dependências do Go desde a versão 1.11. Antes do Modules, o Go utilizava o GOPATH, que tinha algumas limitações em relação ao versionamento e ao gerenciamento de projetos com diferentes dependências.

Pense no go mod como o "gerente de projetos" das suas dependências em Go. Ele acompanha quais bibliotecas externas o seu projeto precisa e em quais versões.

Essa biblioteca é um código externo ao seu projeto. Para que o seu projeto funcione corretamente, o Go precisa saber:

Que essa biblioteca é uma dependência.
Qual versão dessa biblioteca usar.

É aqui que o go mod entra em ação. Ao executar o comando go mod init seu_projeto, você inicializa um módulo Go no seu projeto. Isso cria um arquivo chamado go.mod na raiz do seu projeto. Esse arquivo serve como um "arquivo de projeto" para as suas dependências.

Quando você roda comandos como go get github.com/mattn/go-sqlite3, o Go Modules:

Baixa a biblioteca especificada.
Registra essa biblioteca como uma dependência no seu arquivo go.mod, juntamente com a versão utilizada.

Arquivo go.sum

O arquivo go.sum é um arquivo que o Go Modules mantém para garantir a integridade e a segurança das suas dependências. Ele contém os hashes criptográficos (somas de verificação) das versões específicas de cada dependência listada no seu go.mod.

Pense no go.sum como uma "lista de verificação de segurança" para as suas dependências.

Quando o Go baixa uma dependência, ele calcula o hash do conteúdo baixado e compara com o hash esperado que está no go.sum. Se os hashes não coincidirem, isso indica que a dependência pode ter sido alterada de forma inesperada (potencialmente por alguém malicioso), e o Go se recusará a usá-la, protegendo o seu projeto contra vulnerabilidades ou alterações indesejadas.

O go.sum garante que você está usando exatamente as mesmas versões das bibliotecas que foram usadas em builds anteriores ou por outros membros da sua equipe, evitando problemas causados por alterações inesperadas nas dependências.

Em resumo:

go mod (arquivo go.mod): Gerencia as dependências do seu projeto, rastreando quais bibliotecas externas são necessárias e em quais versões. Foi necessário para que o Go soubesse que você precisava da biblioteca SQLite e qual versão usar.
go sum (arquivo go.sum): Garante a integridade e a segurança das suas dependências, verificando se os arquivos baixados correspondem aos hashes esperados. Foi necessário para garantir que a biblioteca SQLite que você baixou não foi adulterada.
Ambos os arquivos (go.mod e go.sum) trabalham juntos para fornecer um sistema robusto e confiável de gerenciamento de dependências no Go, tornando os projetos mais fáceis de construir, compartilhar e manter.