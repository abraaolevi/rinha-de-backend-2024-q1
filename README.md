# Rinha de Backend

- Go
- Nginx
- Postgres

## Evolução do projeto

- v1: 

Implementação simples, utilizando o novo Router nativo do Go 1.22, e o template do `docker compose` default da rinha, e transactions usando a lib `pgxpool` (v5), apenas para fazer funcionar. Porém ao rodar o teste de carga pela primeira vez fui surpreendido com vários KOs de `j.i.IOException: Premature close`

- v2: 

Com o objetivo de eliminar os KOs, comecei a brincar com os valores de CPU e memória sem muito sucesso. Alterar o modo de rede de `bridge` para `host` no macOS pareceu ter tido uma melhora, mas não muito significativo.

Quando eu esbarrei [neste artigo falando sobre TABLE UNLOGGED](https://www.crunchydata.com/blog/postgresl-unlogged-tables) a performance melhorou bastante. O `unlogged` basicamente remove o Write-Ahead Logging (WAL) que é um recurso do PostgreSQL que permite a recuperação de dados em caso de falhas ou interrupções do sistema.

Seguindo o objetivo da rinha que é performance (ou eu acho que é.. além do aprendizado é claro), acho que não seria um problema desligar o WAL.

Quando uma transação é executada no PostgreSQL, as mudanças são gravadas inicialmente em um arquivo de log, e só então são escritas no disco no `checkpoint`. Isso garante que as alterações sejam registradas no log antes de serem confirmadas como permanentes, e evita a perda de dados em caso de falhas do sistema.

Então um outro "hack" seria aumentar o `checkpoint_timeout` e o `max_wal_size` para diminuir os acessos em disco.

- v3:

Ainda tentando achar otimizações simples, adicionei indices nas tabela, mesmo nas chaves primarias. (Nesse caso também não senti diferença, porém mantive essa alteração)

Nessa versão, resolvi alterar a estratégia, pois até então eu estava: abrindo uma transaction, lendo do banco, validando no código, alterando no código, escrevendo no banco e finalizando a transaction. A nova estratégia foi jogar toda essa regra para o banco usando uma _function_. No geral todos os indicadores melhoraram e não teve KO. Então mantive essa estratégia.