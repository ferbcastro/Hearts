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

# Jogo

- Jogo usa as abstrações implementadas pela biblioteca de rede
- Define sua própria mensagem com tipo e demais campos
