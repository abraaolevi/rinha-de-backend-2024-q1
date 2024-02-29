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

TODO:
- tunar nginx
- tunar banco