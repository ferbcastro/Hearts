# Copas

# Protocolo

- Os pacotes enviados não tem marcador de início
- O tamanho máximo que pode ser enviado por pacote é de 512 bytes
- Não há detecção/correção de erros
- Timeout de 400ms configurado
- Não lida com sequencialização das mensagens

# Rede em anel

- O conceito do bastão é representado pela variável do pacote TokenBusy
- Não há distinção entre o bastão e a mensagem de dados
- Para enviar na rede:
  - Ao total são feitos duas operações de recv e dois send
  - Espera por mensagem com TokenBusy igual a `false` chegar
  - Ao invés de fazer forward, escreve os dados e destino e define TokenBusy `true`
  - Envia e espera a mensagem percorrer a rede
  - Assim que receber a mensagem de volta, define TokenBusy como `false`
  - Passa bastão para frente
  - Um caso especial ocorre no primeiro envio quando o bastão não está circulando:
    cada cliente tem uma variável `waitForToken` e o criador da rede a inicia com
    `false`, enquanto que demais em `true`. Ao fazer primeiro envio, não espera
    por bastão para mandar pacote e define a variável como `true`
- Para escutar na rede:
  - Caso receba um pacote com outro destino, faz forward do mesmo
  - Se o destino for para máquina que chama `Recv`, faz forward e retorna os dados
  - Se sinal de `broadcast` estiver ativo, repete passo anterior
- Para criar/entrar na rede:
  - `CreateRing` e `EnterRing` foram criadas para esses objetivos
  - `CreateRing` é usado pela máquina que cria a rede, logo recebe um array de IPs de todos
  - O método tenta montar a rede até que um teste da rede resulte em sucesso
  - `EnterRing` é usado pelas demais máquinas e precisa apenas do IP de quem chama
  - O método espera até que a rede esteja completa


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

### Pacote 

Estrutura do pacote de dados que trafega entre os nós. Contém:

- ID do remetente
- ID do destinatário.

- PkgType: Tipo do pacote (DATA, BOOT, FORWARD, BROADCAST, RING_COMPLETE).

- Número de série (Serial)
- flag de ACK para confirmação.

- Campo Data, com os dados codificados via gob.

- TokenBusy: Indica se o token está em uso.

- Ack: Usado para confirmar o recebimento de pacotes.

#### Tipos de Pacotes:

DATA - Dados
BOOT - Pacotes usados durante o boot da rede, utiliza campos de dados como endereço e dados
FORWARD - Pacote usado durante o boot para testar integridade da rede, pacote so é repassado 
BROADCAST - Todos pegam o este pacote
RING_COMPLETE - Pacote usado durante o boot para comunicar que a rede esta pronta

# Jogo

- Jogo usa as abstrações implementadas pela biblioteca de rede
- Define sua própria mensagem com tipo e demais campos


