## Historico

1) Deixei tudo default... estava dando muitos KOs
file:///Users/abraaolevi/dev/rinha-de-backend-2024-q1/load-test/user-files/results/rinhabackendcrebitossimulation-20240227020542526/index.html#requests

file:///Users/abraaolevi/dev/rinha-de-backend-2024-q1/load-test/user-files/results/rinhabackendcrebitossimulation-20240227022037165/index.html#requests

file:///Users/abraaolevi/dev/rinha-de-backend-2024-q1/load-test/user-files/results/rinhabackendcrebitossimulation-20240228174840030/index.html

2) Na primeira versão do código eu estava fazendo a versão classica de "ler, atulizar, gravar" em transactions. Alterei para fazer tudo em um update... teve um boa melhora, caiu pela metade os KOs
3) Alterei a rede do Docker para host... deu uma melhora, mas ainda estava dando muitos KOs
4) Remanejei os recursos, para aumentar o banco (mais CPU e mais memoria)
5) Adicionei TABLE UNLOGGED para melhorar a escrita (https://www.crunchydata.com/blog/postgresl-unlogged-tables)
Realmente fez a diferença: apenas 45 j.i.IOException: Premature close
file:///Users/abraaolevi/dev/rinha-de-backend-2024-q1/load-test/user-files/results/rinhabackendcrebitossimulation-20240228204730977/index.html#stats

segundo teste: Sem KO!!!
file:///Users/abraaolevi/dev/rinha-de-backend-2024-q1/load-test/user-files/results/rinhabackendcrebitossimulation-20240228205402117/index.html

terceira
file:///Users/abraaolevi/dev/rinha-de-backend-2024-q1/load-test/user-files/results/rinhabackendcrebitossimulation-20240228211359978/index.html

6) Adicionei index, mesmo na primary key, para verificar se faz alguma diferença na leitura
create index if not exists idx_account_id on accounts(id);
create index if not exists idx_transactions_id on transactions(id desc);

Pareceu ter uma leve melhora, mas nada significativo

7) resolvi mexer um pouco nas configurações do Nginx e brincar com o numero de workers. Diminui de 1000 para 500 o numero de workers para tentar dar um "respiro" para o servers e pool de conexeções.

Não deu muito certo, tomei 137 j.i.IOException: Premature close
Mas os percentis deram uma melhorada
file:///Users/abraaolevi/dev/rinha-de-backend-2024-q1/load-test/user-files/results/rinhabackendcrebitossimulation-20240228214835088/index.html#ranges

Aumentar para 2000 parece ter tido um resultado ainda melhor: 0 KOs e os percentis firam ainda melhores
file:///Users/abraaolevi/dev/rinha-de-backend-2024-q1/load-test/user-files/results/rinhabackendcrebitossimulation-20240228220304050/index.html

Pode ter sido o fator "sorte"

8) Distribuindo os recursos para ter mais recursos para o banco
app 50mb 1.5 cpu
nginx 50bm 0.2

9) Reimplementei jogando todo o processamento da transacion pro banco usando uma function do banco de dados, melhorou ainda mais os percentis..

file:///Users/abraaolevi/dev/rinha-de-backend-2024-q1/load-test/user-files/results/rinhabackendcrebitossimulation-20240229115249261/index.html

10) Aprendi um novo comando: `docker stats` pra ver a utilização dos recursos dos containers

file:///Users/abraaolevi/dev/rinha-de-backend-2024-q1/load-test/user-files/results/rinhabackendcrebitossimulation-20240229200333624/index.html

11) Aprendei sobre Lock otimista e Lock pessimista, mas não apliquei lock otimista

12) Fiz uma versão utilizando traefik e fiber, mas os resultados não deram KO, não parei para investigar mais a funco. Por enquanto a versão `cmd/poc/main.go` foi a melhor -> jogando a transaction direto para o banco

## TODO

- [] Corrigir a versão `cmd/api/main.go` e resolver o problema dos KO de validação
- [] Corrigir a versão `cmd/fiber/main.go` e resolver o problema dos KO de validação
- [] Testar um banco de dados não relacional (talvez mongoDb)

### Referencias

- https://sematext.com/blog/postgresql-performance-tuning/
- https://shijuvar.medium.com/building-rest-apis-with-go-1-22-http-servemux-2115f242f02b
- https://donchev.is/post/working-with-postgresql-in-go-using-pgx/
- https://www.akitaonrails.com/2023/09/20/akitando-145-16-linguagens-em-16-dias-minha-saga-da-rinha-de-backend
- https://gist.github.com/rgreenjr/3637525
- https://hackernoon.com/postgresql-transaction-isolation-levels-with-go-examples-lt5g3yh5
- https://www.crunchydata.com/blog/prepared-statements-in-transaction-mode-for-pgbouncer