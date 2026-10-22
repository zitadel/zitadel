# Account Manager

Aplicação Angular de autosserviço isolada. O layout usa as páginas **01 · Proposta autoral · Desktop** e **02 · Proposta autoral · Mobile** do [Figma aprovado](https://www.figma.com/design/CtI7rVbSery3hs7V316EcO).

Em 22/10/2026 o operador priorizou **a interface e os contratos**, adiando configuração e testes de uso. A auditoria autenticada posterior confirmou carregamento de identidade, navegação e diálogos; operações que gravam credenciais/dados e cenários adversariais continuam com aceite pendente (ver REVIEW.md). Não houve publicação do serviço.

## Fontes de verdade

No Obsidian:

- `Plataformas/Identidade/Produto/login-portal.md`, seção `PLANO DETALHADO — Implementação do Account Manager (21/10/2026)`;
- `Plataformas/Identidade/Produto/branding-ailton-dev-account.md`;
- `Plataformas/Identidade/Arquitetura/arquitetura.md`;
- `Plataformas/Identidade/Arquitetura/modelo-acesso.md`;
- `Plataformas/Identidade/Operacao/tutorial-universal-integracao-aplicacoes-zitadel.md`;
- `Plataformas/Identidade/Operacao/runbook-deploy.md`.

O operador retomou o code review e a conferência de cobertura; configuração e homologação real continuam pendentes. Consulte [REVIEW.md](REVIEW.md) para achados, correções e matriz completa de fluxos. A direção anterior adiou QA funcional; ela também prevalece sobre a sugestão histórica de novos testes unitários no plano. Nenhum teste unitário foi criado.

## Apresentação

- Desktop: lateral editorial de 247 px, identidade visual dinâmica, grupo de navegação pessoal, breadcrumb e menu da conta com avatar, perfil, aplicativos, atividade, idioma, aparência e logout.
- Tablet: lateral compacta e colunas fluidas. Mobile: cabeçalho de marca, abas textuais com indicador de seleção, conteúdo em uma coluna e controles acessíveis.
- Visão geral: identidade real, status de verificação real e linhas de acesso ao perfil e à segurança.
- Perfil: formulário nativo em duas colunas, seção de contatos nativa, foto/visibilidade na coluna lateral no desktop e antes do formulário no mobile e ações que abrem os diálogos nativos.
- Segurança: linhas de senha/passkeys/MFA, painel editorial de proteção e provedores vinculados em leitura. Subpáginas compõem as interfaces nativas.
- PT/EN/ES com `@ngx-translate/core`, mantendo os catálogos nativos intactos. Foco visível, atalho para conteúdo, landmarks e redução de movimento.
- Fonte, imagens e cores da Organization via `getMyLabelPolicy`, respeitando a herança do servidor. Não há logo local, nome de marca fixo, fonte de marca fixa ou cores do Figma usadas como configuração de produto.

Os campos adicionais do formulário nativo (nome de exibição, apelido e gênero) foram preservados; o Account altera sua disposição visual, sem copiar o formulário. Os dados ilustrativos do design não foram copiados. Aplicativos integram os grants reais da própria conta; organizações, assinaturas comerciais e revogação de dispositivos continuam sem integração. O método legado `listMyUserSessions` foi identificado, mas não oferece por si só a experiência completa de dispositivos/revogação.

## Isolamento e reuso

O `console/angular.json`, o `console/package.json`, os targets Nx preexistentes, `console/src/app`, `apps/login`, Go/Core e protobufs não foram alterados.

`build.mjs` usa o Angular Architect já instalado para acrescentar **em memória** o projeto de `account/angular.json`. Isso preserva `ng build` e `ng serve` do workspace administrativo original. Saída independente: `console/account/dist`.

Reutilização direta, sem cópias dos componentes/serviços:

| Capacidade              | Implementação nativa                                                                                                                                                 |
| ----------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Formulário e avatar     | `DetailFormModule`, `DetailFormComponent`, `ProfilePictureComponent`, `AvatarModule`, `AssetService`                                                                 |
| Contatos e verificações | `ContactComponent`, `PhoneDetailComponent`, `EditDialogComponent`, `CodeDialogComponent`, `GrpcAuthService`                                                          |
| Senha                   | `PasswordComponent`, `UserService`, `NewAuthService`, validadores/política de complexidade nativos                                                                   |
| MFA                     | `AuthUserMfaComponent`, `AuthFactorDialogComponent`, `DialogU2FComponent` e serviços nativos                                                                         |
| Passkeys                | Tabela/remoção de `AuthPasswordlessComponent`; adaptador Account usando Auth API `AddMyPasswordlessLink` / `SendMyPasswordlessLink`, com cadastro no Login do issuer |
| IdPs vinculados         | `ExternalIdpsComponent` com `GrpcAuthService`, colunas somente de leitura                                                                                            |
| Transporte              | `GrpcService` e seus clientes/interceptadores existentes                                                                                                             |
| Tradução e apresentação | `TranslateModule`, `InputModule`, `CardModule`, `WarnDialogModule`, `DetailLayoutModule` e módulos Material                                                          |

**Origem WebAuthn:** cadastro de passkeys oferece continuar no dispositivo, via link self nativo adaptado à rota oficial Login V2, ou envio por e-mail via Auth API self. A chamada V2 `CreatePasskeyRegistrationLink` foi removida por exigir `user.passkey.write` no interceptor; não conceder papéis administrativos. Entrega e cadastro final ainda exigem homologação pelo operador. Cadastro U2F só é oferecido na origem do issuer, pois o backend nativo restringe RP/origens. Em domínio separado há aviso explícito e esse fluxo continua pendente; listagem/remoção são nativas. Ver [revisão e gates](REVIEW.md).

Não importar `AppModule`, `AppRoutingModule`, `UsersModule` ou `UserDetailModule` administrativos.

**Fronteira de build necessária:** `AssetService` importa o enum público `Theme` de dentro de um componente administrativo. O build Account substitui exclusivamente esse módulo por `native-asset-theme.ts` (os dois valores do enum), com verificação do contrato no build. Isso impede que o compilador traga a tela administrativa de branding para o Account. O arquivo/serviço original e o build do Console permanecem intactos. Reavaliar essa fronteira se o enum nativo mudar.

## Rotas

Protegidas: `/`, `/perfil`, `/seguranca`, `/seguranca/senha`, `/seguranca/passkeys`, `/seguranca/mfa`.

Infraestrutura: `/auth/callback`, `/auth/error`, `/signedout`. O wildcard produz estado seguro de não encontrado. Não existem parâmetros de ID de usuário nas rotas Account. Nenhuma rota administrativa é registrada, inclusive para pessoas com permissões IAM.

## Contrato OIDC preparado

Biblioteca existente `angular-oauth2-oidc`; Code + PKCE S256, state/nonce, issuer/audience/expiração, callback fixo, transação de uso único com validade de dez minutos e retorno por allowlist exata. Queries e URLs externas não são destinos de retorno. A configuração é exclusiva da aplicação e o storage usa o prefixo `account:v1:`.

A assinatura RS256 do ID token usa WebCrypto por um `ValidationHandler` da biblioteca: a implementação padrão dessa biblioteca não valida assinaturas. Não foi introduzido outro cliente HTTP/OIDC. Algoritmos diferentes de RS256 falham fechados e exigem revisão explícita.

Scopes solicitados:

```text
openid profile email
urn:zitadel:iam:user:resourceowner
urn:zitadel:iam:org:id:<LOGIN_ORG_ID>
urn:zitadel:iam:org:project:id:zitadel:aud
```

O último scope indica a audiência da API ZITADEL; não concede papéis administrativos. O contexto exige a claim `urn:zitadel:iam:user:resourceowner:id` e confirma o sujeito/owner por `GetMyUser` antes da abertura das páginas. Portanto, a aplicação deve disponibilizar as claims solicitadas no ID token; a configuração real ainda não foi homologada. Políticas e autorização continuam no ZITADEL.

Logout usa o mecanismo da biblioteca/provedor; expiração fecha os diálogos, limpa o estado de autenticação e o cache de consultas e retorna ao estado seguro. Isso ainda exige teste real. Não compartilhar tokens, secrets ou callbacks com o Console ou com o PasteMax.

### Registro criado antes do adiamento da configuração

Com autorização explícita do operador, no issuer `https://auth.ailton.dev.br`:

- Organization informada: `395747685065297410` — Ailton Dev Account;
- novo Project: `396191836894399748` — Account Manager;
- nova Application: `396191976497614084` — Account Manager Web;
- Client ID público: `396191976497679620`;
- tipo User Agent/SPA, Code + PKCE, autenticação None, sem segredo;
- callback registrado: `http://localhost:8088/auth/callback`;
- pós-logout registrado: `http://localhost:8088/signedout`;
- modo de desenvolvimento apenas nessa aplicação, para os callbacks HTTP locais;
- Login V2 selecionado somente nessa aplicação, URL `https://auth.ailton.dev.br/ui/v2/login`.

O Project fornecido inicialmente (`395747873473433090`) foi conferido como **Paste Max**; não foi reutilizado nem modificado. `pastemax-service` não foi usado no Account.

A opção **Include user's profile info in the ID Token** permaneceu desativada: a revisão automática rejeitou sua alteração por entender que amplia a exposição de claims além do registro autorizado. Em seguida o operador adiou configuração e testes. Não tentar contornar essa rejeição. Retomar a decisão explícita sobre claims antes de homologar OIDC. Nenhum domínio/callback público do Account foi publicado ou registrado nesta execução.

### Variáveis e publicação futura

`ACCOUNT_ISSUER`, `ACCOUNT_API`, `ACCOUNT_CLIENT_ID`, `ACCOUNT_LOGIN_ORG_ID` são configurações públicas. Não há PAT, secret OIDC, token M2M ou credencial de administração no frontend. O Project ID não é substituto do Organization ID nem é necessário para o scope de login.

Para execução estática local, fornecer esses campos como `api`, `issuer`, `clientid`, `loginOrgId` em `assets/environment.json` **do artefato**, sem versionar os valores no código. Há um exemplo deliberadamente inválido que falha fechado até receber configuração real.

`Dockerfile`, `compose.yml`, `10-account-config.sh` e `nginx.conf.template` preparam um serviço independente. O entrypoint valida origens HTTPS/ID organizacional e gera o JSON sem interpolação insegura. O Nginx aceita apenas os hosts `conta.ailton.dev.br`, `localhost` e `127.0.0.1`, e as rotas Account/OIDC explicitamente permitidas; caminhos desconhecidos e outros hosts recebem HTTP 404 com página própria. Assets e health check `/healthz` têm locations dedicados; o Compose verifica essa rota antes de considerar o serviço saudável. O Compose usa `ghcr.io/ailtonacr/ailton-account-zitadel:latest` como imagem padrão. Também aplica CSP, no-store, nosniff e bloqueio de framing. O Compose publica somente `127.0.0.1:8088`, sem modificar serviços existentes.

Domínios confirmados pelo operador: Login/API `auth.ailton.dev.br`; Account `conta.ailton.dev.br`. São origens diferentes: WebAuthn do Login não pode ser iniciado diretamente na origem Account.

Publicação futura: configurar hostname/HTTPS exclusivamente no serviço Account (`conta.ailton.dev.br`), registrar callback HTTPS exato e pós-logout, revisar CORS/CSP/claims e homologar. Não apontar o domínio Account para o bundle do Console. Rollback: restaurar somente a imagem anterior do Account ou parar seu serviço; não alterar API, Login, proxy do IdP ou banco. Produção depende de autorização específica.

## Comandos (da raiz)

```sh
pnpm install --frozen-lockfile
pnpm nx run @zitadel/console:account-build
pnpm nx run @zitadel/console:account-typecheck
pnpm nx run @zitadel/console:account-lint
node console/account/build.mjs development
pnpm nx run @zitadel/console:account-dev
```

Os targets dependem da geração nativa e da biblioteca compartilhada. `go.mod` controla o Go usado pelo target preexistente de geração; não alterar targets existentes para contornar ambiente local. O dev server suporta a compilação; não configura cliente, issuer ou credenciais automaticamente.

Sem Nx, depois da geração/build dos clientes nativos, `node console/account/build.mjs production` produz o mesmo artefato.

## Estado de validação e checkpoints

| Checkpoint | Evidência e limite                                                                                                                                                                                                                                                       |
| ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| CP-00      | Checkout inicial limpo, branch `feat/ailton-account`, base `123c7079b96b30398ccde5a548dac22eb9471b5c`. Código, Obsidian, AI Memory e frames Figma consultados.                                                                                                           |
| CP-01      | Bootstrap, configuração e artefato próprios compilados. Targets existentes preservados.                                                                                                                                                                                  |
| CP-02      | Contrato OIDC e guard implementados; novo cliente registrado. Sessão autenticada e contexto self observados posteriormente na aba local; matriz completa de login/callback/logout/SSO/outra Organization ainda pendente.                                                 |
| CP-03      | Composição de perfil, contatos, verificações, username e avatar com componentes e serviços nativos. Operações reais pendentes.                                                                                                                                           |
| CP-04      | Composição nativa de senha/MFA/passkeys e IdPs somente leitura. Operações reais pendentes.                                                                                                                                                                               |
| CP-05      | Apresentação desktop/tablet/mobile, branding dinâmico e PT/EN/ES implementados a partir do Figma. Telas principais conferidas posteriormente com branding real em desktop/mobile e temas; estados de escrita/credenciais e matriz completa de contraste ainda pendentes. |
| CP-06      | Sem páginas fictícias ou dados comerciais. Capacidades opcionais sem fonte integrada foram omitidas.                                                                                                                                                                     |
| CP-07      | Typecheck, lint e builds; nenhum teste unitário novo. Usabilidade, permissões reais e fluxos OIDC adiados explicitamente.                                                                                                                                                |
| CP-08      | Artefato/Docker/proxy/variáveis/rollback preparados; sem deploy. A imagem local foi construída antes dos últimos refinamentos visuais e precisa de rebuild na retomada.                                                                                                  |

Resultados registrados nesta execução:

A auditoria autenticada posterior está em [REVIEW.md](REVIEW.md): leituras e rotas principais reais, diálogos abertos/cancelados, validação de perfil inválido sem gravação, PT/EN/ES e desktop/mobile. Criação de passkeys, envio de mensagens e alterações de credenciais permanecem para o operador.

- Account: TypeScript sem erros; ESLint passou; builds production e development passaram nas composições verificadas. O build de produção inclui aproximadamente 9,1 MB brutos, com alertas herdados de CommonJS/Sass e budget de aviso; não há erro de budget.
- Console original: `pnpm nx run @zitadel/console:build --excludeTaskDependencies` passou após geração dos clientes. Seus arquivos, configuração Angular e definições preexistentes dos targets estão intactos.
- Login V2 original: `pnpm nx run @zitadel/login:build` passou, sem alterações em `apps/login`.
- Inspeção estática do bundle: sem roteador/páginas administrativas identificados (`AppRoutingModule`, `OrgCreateComponent`, `PrivateLabelingPolicyComponent` e registros de rotas administrativas ausentes). Isso não substitui teste adversarial autenticado.
- Antes do adiamento, o navegador carregou o shell local e o estado de erro seguro sem configuração. Nenhum cenário autenticado, alteração de perfil, senha, MFA, passkey ou cadastro foi certificado.
- O endpoint de discovery de desenvolvimento retornou HTTP 502 na tentativa realizada. Os IDs fornecidos e o cadastro autorizado foram conferidos em produção; isso não prova homologação de dev.
- O Docker multi-stage completou um build local preparatório. Não houve `docker compose up`, deploy no Coolify, alteração do proxy ou publicação.

Os logs de trabalho ficaram em `/tmp/account-*.log`; são evidência transitória, não dependência do projeto. Não registrar tokens ou dados pessoais nos relatórios.

## Arquivos da entrega

```text
console/
├── project.json                     # apenas quatro targets Account adicionais
├── tsconfig.account.json
├── src/account/
│   ├── account.module.ts
│   ├── main.ts
│   ├── index.html
│   ├── routes.ts
│   ├── shell.ts
│   ├── shell.html
│   ├── page.html
│   ├── styles.scss
│   ├── branding.ts
│   ├── auth.service.ts
│   ├── signature.ts
│   ├── storage.ts
│   ├── native-security.ts
│   ├── native-presentation.ts
│   ├── identity.ts
│   ├── native-asset-theme.ts
│   └── assets/
│       ├── environment.example.json
│       └── i18n/{pt,en,es}.json
└── account/
    ├── .gitignore
    ├── .env.example
    ├── angular.json
    ├── build.mjs
    ├── Dockerfile
    ├── Dockerfile.dockerignore
    ├── 10-account-config.sh
    ├── nginx.conf.template
    ├── compose.yml
    ├── REVIEW.md
    └── README.md
```

Somente `console/project.json` foi modificado entre os arquivos preexistentes; os demais itens são novos. Sem push ou deploy autorizado nesta etapa.

## Contratos de aplicativos, atividade e preferências (22/10/2026)

- `/aplicativos`: `GrpcAuthService.listMyUserGrants`, somente grants self, paginação por offset, agrupamento por Project/Organization e estado real. Tags seguem o padrão do operador `{app}.{group}.{tier}`: `pastemax.plan.free` aparece como **Free**, com chave original apenas no tooltip. Não há seção “Níveis de acesso”. Uma tag de grant não consulta nem comprova cobrança/assinatura comercial.
- `/atividade`: `GrpcAuthService.listMyUserChanges`, o mesmo endpoint da timeline MYUSER do Console. Eventos em ordem decrescente, responsável e data/hora no fuso do navegador; paginação por sequência, deduplicação, fim da lista, timeout, erro/retry e atualização. Não aceita userId arbitrário. A troca de idioma refaz a consulta para obter `localizedMessage` da API via interceptor nativo `Accept-Language`. Retenção e eventos disponíveis dependem do servidor. Não é uma lista de sessões ou revogação de dispositivos.
- Textos estruturais: catálogos oficiais Console PT/EN/ES + extensões Account. Nomes de idiomas usam `Intl.DisplayNames` na língua ativa (Português/Inglês/Espanhol). Dados, políticas, marca e mensagens de eventos vêm da API. A API de traduções do Hosted Login é específica do Login V2, com contrato diferente e leitura `iam.policy.read`; não foi exposta ao SPA nem foram concedidas permissões administrativas para consumi-la. Configurar traduções do Login não altera automaticamente os rótulos Account.
- Idioma/aparência do menu persistem como preferências locais, sem gravar o perfil; salvar idioma do perfil usa a operação nativa. Foto passa por compressão JPEG com orientação e proporção preservadas, até 1024 px e no máximo 480.000 bytes; metadados removidos e transparência sobre branco. Arquivos inválidos/acima de 30 MB falham localmente. Upload/remoção e atualização continuam pelo fluxo nativo; não há upload se a compressão falhar.
- Passkeys: Fechar + Cadastrar no dispositivo são as ações principais. E-mail é alternativa secundária porque existe no autosserviço nativo; nenhum envio automático.
- Logout: a biblioteca OIDC encerra a sessão central correspondente ao ID token pelo `end_session_endpoint` e limpa o armazenamento local. A desconexão imediata dos cookies próprios de cada produto depende de integração desse produto com back-channel logout; não equivale a encerrar todas as sessões em todos os dispositivos. Referência: https://zitadel.com/docs/guides/integrate/login/oidc/logout.
