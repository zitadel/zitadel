# Code review — Account Manager — 22/10/2026

## Correção posterior — indicador de carregamento preso

O operador relatou que a visão geral mantinha “Carregando sua conta…”. O runtime já continha as quatro configurações públicas; não faltava variável. A dependência instalada é Angular 21, cujo bootstrap Account não habilitava `provideZoneChangeDetection()`, ao contrário do `console/src/main.ts` nativo. Os componentes reutilizados e o shell atualizam campos comuns após RPCs, dependendo desse mecanismo para refletir a conclusão na tela. O bootstrap Account foi alinhado ao nativo, sem alterar o Console. Typecheck, ESLint, build de produção e diff-check passaram. Após recarregar a aba local autenticada, a visão geral e a identidade na lateral apareceram automaticamente, sem clique adicional para disparar atualização. Esse cenário específico foi confirmado no navegador; os demais gates funcionais continuam com o estado descrito abaixo.

Conclusão: a composição principal está implementada, mas **o aceite integral dos fluxos continua aberto**. Typecheck e build verificam contratos de compilação; não comprovam operações autenticadas, entrega de mensagens ou paridade visual no navegador.

## Auditoria autenticada e refinamento visual

Com a conta autenticada pelo operador, a leitura real da identidade e as rotas de autosserviço foram inspecionadas no Chrome local. Referências Figma reconferidas: visão geral `14:11`, perfil desktop `14:136`, segurança `14:239` e perfil mobile `17:141`. Branding e identidade continuam reais, sem copiar os exemplos do design.

Correções desta rodada:

- Estrutura de avatar circular, imagem com recorte correto e cores da Organization. Contatos com rótulos, espaçamento, linhas e ações legíveis. Cards/tabelas/estados vazios alinhados à hierarquia editorial, preservando templates nativos.
- Registro dos locales Angular PT/EN e Moment PT, eliminando a causa dos erros `NG02100` observados ao exibir datas do paginador.
- Idioma real `und`/vazio representado por “Conforme navegador”, sem perder o valor persistido; nomes acessíveis dos selects de gênero/idioma.
- Falhas de consulta de métodos/políticas em passkeys/MFA agora têm aviso explícito na página, além do estado vazio nativo.
- Bloqueio local de envio de perfil inválido com mensagens nativas; limite de 15 segundos para leitura da página e ação de retry já existente.
- Controles de ícone com nome acessível, verificação de contatos acessível por teclado, expandir/recolher com semântica e PT/EN; nome acessível da navegação sem o marcador decorativo.
- Identidade lateral atualizada após refresh de perfil/avatar por um adaptador Account, sem modificar `GrpcAuthService`.
- Dependência nativa `TypeSafeCellDefModule` incluída para linhas de MFA; flags nativas no seletor de país; links de volta de passkeys/MFA; footer sem duplicação do texto contextual.

### Correção de autorização de passkeys relatada pelo operador

O erro real `membership not found (AUTHZ-cdgFk)` revelou que a chamada V2 `CreatePasskeyRegistrationLink` era inadequada para autosserviço: sua annotation exige `user.passkey.write` e o interceptor consulta membership antes do comando. A permissão self interna no comando não elimina esse gate. A conclusão estática anterior era incompleta.

A chamada foi removida do Account. O adaptador usa exclusivamente os métodos nativos Auth API `AddMyPasswordlessLink` e `SendMyPasswordlessLink`, cujas annotations exigem `authenticated` e cujo usuário é obtido do contexto do token. Não se concede papel administrativo.

Para cadastro no dispositivo em origem separada, o link nativo retornado é validado por origem/path/subject/Organization e seus parâmetros documentados são adaptados aos nomes esperados pela página oficial Login V2 `/ui/v2/login/passkey/set`. A compatibilidade foi conferida no gerador de links, write model de códigos e página/componente nativos. O código de inscrição temporário permanece somente na memória do diálogo e no link de uso único exigido pelo fluxo nativo: não é persistido ou registrado em logs, nem é um access/ID token OIDC. O operador abre o Login em nova aba com `noopener noreferrer`, conclui o cadastro e retorna para atualizar a lista. O envio opcional usa o link gerado pelo servidor Auth API, sem template customizado.

Domínios confirmados: Login/API **auth.ailton.dev.br**, Account **conta.ailton.dev.br**. São origens distintas, assim como localhost/issuer. Não alterar RP ID, allowlist WebAuthn ou permissões para forçar cadastro direto na origem Account. O domínio Account continua sem publicação nesta execução.

## Pendências de aceite

1. **OIDC integral:** a sessão autenticada e leitura self funcionaram na aba local. Não foi repetida toda a matriz de callback/logout/SSO/expiração/outra Organization; não habilitamos configurações antes rejeitadas nem ampliamos claims/permissões.
2. **Cadastro U2F:** permanece indisponível na origem separada; listagem/remoção são nativas. Resolver cadastro exige caminho nativo suportado, sem relaxar origem.
3. **Passkeys:** chamadas self e handoff Login V2 corrigidos; criação efetiva da credencial e entrega do e-mail dependem da validação do operador. A auditoria não criou credenciais nem enviou mensagens. A geração/consumo do link no ambiente real não é certificada por build.
4. **Operações com dados:** senha/perfil/e-mail/telefone/avatar/MFA não foram alterados durante a auditoria. A validação de formulário inválido não salva dados. Estados com fatores/provedores já cadastrados ainda precisam de homologação com conta apropriada.

## Problemas corrigidos nesta revisão

| Prioridade | Problema encontrado                                                                     | Correção                                                                                                                              |
| ---------- | --------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| Alta       | Entrada direta em `/signedout` podia manter tokens locais válidos                       | Limpeza local explícita antes do discovery; nova entrada reinicializa bootstrap/discovery                                             |
| Alta       | Pedido nativo de reautenticação ignorava `force`                                        | Encaminha `prompt=login` e `max_age=0`, com nova transação PKCE                                                                       |
| Alta       | Cadastro nativo de passkey executava WebAuthn no domínio separado                       | Adaptador usa operações Auth API self e handoff ao Login no issuer; chamada V2 administrativa removida após o erro real de membership |
| Alta       | Cadastro U2F oferecido em origem incompatível                                           | Restringe disponibilidade à mesma origem e expõe limitação; gate permanece aberto                                                     |
| Média      | Retry podia não fazer nada após falha de configuração/discovery e perdia destino        | Novo bootstrap, preservação da rota permitida e transação de uso único                                                                |
| Média      | Senha aberta por deep link voltava ao IdP após salvar                                   | Predecessor interno `/seguranca` para o `history.back()` do componente nativo                                                         |
| Média      | Expiração antecipada pelo fator padrão da biblioteca                                    | Timer no instante de expiração; encerra sessão local, diálogos e cache                                                                |
| Média      | Erro de política MFA sem tratamento e campos de contato sem formato nativo de validação | Tratamento seguro de política; validadores nativos de e-mail/telefone/obrigatoriedade                                                 |
| Baixa      | Idioma salvo no perfil não refletido na interface; bootstrap sem indicação              | Aplica PT/EN após salvar e acrescenta status inicial de carregamento                                                                  |

Origens de configuração também exigem HTTPS e valor canônico sem path/query/credenciais. O contrato estrito de Organization foi preservado.

## Rastreabilidade dos fluxos

| Fluxo / rota                                               | Implementação                                                                               | Estado                                                            |
| ---------------------------------------------------------- | ------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| Login, cadastro, recuperação                               | Redirecionamento OIDC para Login V2 nativo; nenhuma tela de autenticação duplicada          | Integração preparada; execução pendente                           |
| Callback, retry, step-up, expiração, logout                | `AccountAuthService`, biblioteca OIDC existente, assinatura RS256 e guard self/Organization | Corrigidos estaticamente; execução pendente                       |
| Visão geral `/`                                            | Identidade e verificação reais via `GetMyUser`; acessos perfil/segurança                    | Implementado; execução pendente                                   |
| Perfil `/perfil`                                           | `DetailFormComponent`, `GrpcAuthService.updateMyProfile`, validadores nativos               | Implementado; execução pendente                                   |
| Username                                                   | Diálogo nativo e `updateMyUserName` self                                                    | Implementado; execução pendente                                   |
| Avatar                                                     | Upload/remoção nativos `DetailFormComponent`/`ProfilePictureComponent`/`AssetService`       | Implementado; execução pendente                                   |
| E-mail                                                     | Diálogo nativo, alteração e reenvio nativos; verificação pelo link do IdP                   | Implementado; mensagem/link pendentes                             |
| Telefone                                                   | Alteração, reenvio, confirmação por código e remoção nativos                                | Implementado; mensagem/operações pendentes                        |
| Senha `/seguranca/senha`                                   | `PasswordComponent`, senha atual e política nativa, self `UserService`                      | Implementado; execução pendente                                   |
| Passkeys `/seguranca/passkeys`                             | Tabela/remoção nativas; adaptador Auth API self com handoff no Login V2                     | Integração preparada; cadastro/uso pendentes                      |
| MFA `/seguranca/mfa`                                       | Listagem/remoção nativas; diálogos OTP, SMS, e-mail conforme política/contato               | Implementado; execução pendente; cadastro U2F separado incompleto |
| IdPs `/seguranca`                                          | `ExternalIdpsComponent`, colunas somente leitura                                            | Implementado; execução pendente                                   |
| Marca/tema/PT/EN/ES                                        | Política efetiva da Organization, assets nativos, tradução nativa + estrutural              | Implementado; aceite visual pendente                              |
| Ajuda/privacidade                                          | Links HTTPS da política efetiva, omitidos se ausentes                                       | Implementado; execução pendente                                   |
| Aplicativos `/aplicativos`                                 | Grants self paginados, estado real e tags legíveis                                          | Leitura real validada; nenhum dado ilustrativo                    |
| Atividade `/atividade`                                     | Histórico MYUSER nativo, paginação por sequência e mensagens localizadas pela API           | Leitura real validada; sem logs administrativos                   |
| Organizações, sessões/dispositivos, assinaturas comerciais | Sem fonte/fluxo oficial completo integrado                                                  | Permanecem omitidos                                               |

## Alinhamento e limites

- Plano: Obsidian `Plataformas/Identidade/Produto/login-portal.md`, seção `PLANO DETALHADO — Implementação do Account Manager (21/10/2026)`, CP-00 a CP-08. Gates funcionais continuam abertos.
- Design: [Figma aprovado](https://www.figma.com/design/CtI7rVbSery3hs7V316EcO), páginas Desktop `14:5` e Mobile `17:5`; segurança `14:239` reconferida nesta revisão. Dados ilustrativos não viraram dados funcionais.
- Componentes/templates/serviços nativos reutilizados. Adaptadores Account não modificam `console/src/app`, `apps/login`, Go/Core/protos nem targets preexistentes.
- Sem testes unitários novos, acesso de escrita a produção, push ou deploy. A imagem Docker anterior precisa de rebuild.
- Referências do contrato de passkeys: [conceitos oficiais](https://zitadel.com/docs/concepts/features/passkeys), [integração Login UI](https://zitadel.com/docs/guides/integrate/login-ui/passkey); implementação local em `proto/zitadel/user/v2/auth.proto`, `internal/command/user_v2_passkey.go`, `internal/command/permission_checks.go` e `apps/login/src/app/(login)/passkey/set/page.tsx`.

## Verificação executada

TypeScript Account, ESLint Account, Prettier, build de produção Account e `git diff --check`. Resultado final de cada comando registrado após a última alteração. Os builds anteriores de Console e Login não foram repetidos porque suas fontes não mudaram; isso não constitui homologação dos fluxos acima.

## Evidência de navegação desta auditoria

- Visão geral, perfil, segurança, senha, passkeys e MFA carregaram na sessão local autenticada.
- Perfil comparado em 1440 × 960 e 390 × 844; hierarquia, colunas, avatar, contatos e ação de salvar conferidos com os frames aprovados. Os campos nativos adicionais são intencionais.
- Temas claro/escuro e troca estrutural PT/EN conferidos; a escolha do idioma da interface não grava o perfil.
- Campo obrigatório vazio seguido de Salvar exibiu validação nativa e aviso Account, sem chamar o método de gravação; o valor local foi restaurado e a página recarregada.
- Diálogos de e-mail, telefone/países e foto abriram e foram cancelados/fechados; nenhum dado foi gravado ou arquivo enviado.
- Seleção de MFA abriu apenas o método TOTP permitido pela política e foi fechada antes de iniciar registro.
- Novo diálogo de passkey exibiu o issuer real e as escolhas de dispositivo/e-mail; foi fechado antes de gerar link, enviar mensagem ou criar credencial.
- Paginador de provedores exibiu data PT corretamente após registro do locale; não houve nova ocorrência do erro de locale na navegação final.
- Servidor local reiniciado quando perdeu o processo; runtime reaplicado a partir de `console/account/.env` ignorado pelo Git.

Essas observações não substituem teste de cadastro/remoção real de fatores, entrega transacional, contas com tabelas populadas nem a matriz adversarial de autorização. Não houve push/deploy ou concessão de papéis.

## Revisão de UX, aplicativos e atividade (22/10/2026)

Solicitações posteriores do operador implementadas: aplicações por grants self; tags `{app}.{group}.{tier}` (Free em vez de pastemax.plan.free), sem título “Níveis de acesso”; histórico da conta com eventos reais; menu único com avatar; PT/EN/ES e nomes de idiomas localizados; foto antes do formulário mobile e compressão automática; ações principais de passkey e e-mail como alternativa nativa. O frame Figma Aplicativos `14:339` foi consultado; as alterações explícitas do operador para menu/foto têm precedência sobre a composição anterior.

O problema dos ícones cortados era uma regra tipográfica do menu que sobrescrevia a fonte Material Icons. A regra agora aplica a fonte de marca somente ao texto. Fonte computada e captura mobile confirmam ícones gráficos de 24 px, sem ligaturas como “person_outline” cortadas. O seletor da marca também foi limitado ao logo direto para não recortar o avatar aninhado.

Compressão exercitada no navegador com imagem sintética de ruído: PNG de 6.631.183 bytes convertido em JPEG de 360.194 bytes, 1024 × 768. Tipo inválido rejeitado; nenhum arquivo pessoal enviado à API nessa verificação. Isso verifica a preparação local, sem certificar upload/remoção reais.

Histórico real carregado em português com datas e responsáveis, inclusive autenticações e alterações de perfil. A conta disponível possui menos de 30 eventos: fim de lista observado; paginação de histórico maior permanece sem aceite com dados reais. Controle de erro não transforma indisponibilidade em lista vazia. Mensagens de eventos são API; textos estruturais são catálogos, conforme contrato no README. Não houve concessão de papéis ou alteração em traduções/políticas do IdP.

Verificação final desta revisão: TypeScript, ESLint Account, build de produção, Prettier e diff-check passaram. Apps conferidos em 320/390/768/1440 px sem overflow horizontal do documento; navegação e rótulos PT/EN/ES conferidos e restaurados para português. Diálogo de passkey conferido em 390 px com Fechar/Cadastrar no dispositivo e e-mail secundário, fechado sem executar registro/envio. Menu mobile inspecionado com fonte Material Icons e captura salva; viewport restaurado ao final. Obsidian atualizado na seção 12 e relido. Fontes nativas protegidas continuam sem alterações desde a base.

O Nginx foi restringido aos hosts do Account (produção e localhost) e allowlist de rotas funcionais, callbacks e logout. Hosts/caminhos fora da lista retornam 404 e página própria; assets não existentes também 404. Revisão estática da configuração; a imagem Nginx não foi executada nesta etapa.
