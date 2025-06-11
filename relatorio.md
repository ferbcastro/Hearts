## Criação do Anel

Um unico computador é reponsável por criar a rede, enquanto os outros esperam pelas intruções

### CreateRing

O primeiro nó (starter) inicializa o anel.

Envia pacotes do tipo BOOT para todos os outros nós,
informando o ID, IP do destinatario e o IP do próximo.
A rede é endereçada por meio do Id adquirido durante 
o boot, portanto durante essa parte inicial essas
mensagens são endereçadas usando o IP no campo de dados.

O starter apos tentar criar todos os links, manda um pacote
FORWARD para percorrer o anel, para testar que todos os links
foram estabelecidos corretamente.

Caso o pacote volte para o starter, envia RING_COMPLETE para informar
que o anel está pronto. Caso não volte ele tentará reestabelecer todos 
os links e repetirá o mesmo processo.

### EnterRing

Um nó que deseja entrar aguarda receber um pacote BOOT, endereçado a ele,
contendo seu ID e o IP do próximo nó.

Após conectar-se ao próximo, o nó aguarda o sinal de RING_COMPLETE para começar
a operar normalmente no anel.

### Comunicação

#### Send

Envia dados para um nó específico. O envio só é permitido se o cliente estiver com o token.


#### Broadcast

Envia dados para todos os nós do anel, independentemente de seus IDs.

#### Recv

Espera por um pacote:

    Se recebe um ponteiro nulo, espera apenas o token.

    Caso contrário, espera um pacote destinado ao nó ou um broadcast e decodifica os dados.

### Pacote 

Estrutura do pacote de dados que trafega entre os nós. Contém:

-ID do remetente
-ID do destinatário.

-PkgType: Tipo do pacote (DATA, BOOT, FORWARD, BROADCAST, RING_COMPLETE).

-Número de série (Serial)
-flag de ACK para confirmação.

-Campo Data, com os dados codificados via gob.

-TokenBusy: Indica se o token está em uso.

-Ack: Usado para confirmar o recebimento de pacotes.

#### Tipos de Pacotes:

DATA - Dados
BOOT - Pacotes usados durante o boot da rede, utiliza campos de dados como endereço e dados
FORWARD - Pacote usado durante o boot para testar integridade da rede, pacote so é repassado 
BROADCAST - Todos pegam o este pacote
RING_COMPLETE - Pacote usado durante o boot para comunicar que a rede esta pronta

### Funcionamento

- Espera o token para falar
- Cada um fala uma vez, espera mensagem retornar e libera o token
- Caso mensagem nao volte ou nao esteja com Ack, ela é reenviada

- Mensagem que nao sejam de broadcast ou não sejam para aquele computador, só são repassadas

- Starter é sempre o primeiro a falar na rede, aquele que montou a rede é considerado como portador inicial do token
