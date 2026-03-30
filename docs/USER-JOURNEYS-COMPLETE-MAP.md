# Mapeamento Completo de Jornadas de Usuário - BGC Platform

**Document Version:** 1.0.0
**Last Updated:** 2026-01-20
**Owner:** Product Management Team
**Status:** Active

---

## Executive Summary

Este documento mapeia **TODAS as jornadas de usuário necessárias** para transformar o BGC de um MVP funcional (simulador de destinos) em um **produto completo de trade intelligence e operações internacionais**.

**Status Atual:**
- ✅ **1 jornada implementada:** Simulação de Destinos (Anonymous User)
- ❌ **47 jornadas faltantes:** Mapeadas neste documento

**Organização:**
- **Framework:** AARRR (Pirate Metrics)
- **Priorização:** North Star Metric = Time-to-First-Export-Match
- **Fases:** MVP Beta → V1 Public Launch → V2 Platform Maturity

---

## Table of Contents

1. [Acquisition - Como Usuários Chegam](#1-acquisition---como-usuários-chegam)
2. [Activation - Primeira Experiência de Valor](#2-activation---primeira-experiência-de-valor)
3. [Retention - Usuários Retornam](#3-retention---usuários-retornam)
4. [Revenue - Monetização](#4-revenue---monetização)
5. [Referral - Network Effects](#5-referral---network-effects)
6. [Administrative Journeys - Operações Internas](#6-administrative-journeys---operações-internas)
7. [Roadmap Visual - Sequência de Implementação](#roadmap-visual---sequência-de-implementação)

---

# 1. ACQUISITION - Como Usuários Chegam

## J-A01: Descoberta Orgânica via SEO/Blog
**Persona:** PME Exportador (iniciante)
**User Story:** Como empresário interessado em exportar, quero encontrar informações sobre mercados internacionais via Google, para avaliar viabilidade sem me cadastrar ainda
**Prioridade:** P1 - High
**Status:** ❌ Não iniciado
**Dependências:** Nenhuma

### Fluxo:
1. Usuário pesquisa "para qual país exportar açúcar" no Google
2. Encontra artigo do BGC Blog: "Top 10 Destinos para Açúcar Brasileiro em 2026"
3. Lê conteúdo educacional (SEO-optimized)
4. Vê CTA: "Descubra os melhores destinos para SEU produto" → redireciona para /simulator
5. Realiza primeira simulação (anonymous)
6. Vê modal: "Criar conta grátis para salvar resultados"

### Páginas/Rotas Necessárias:
- `/blog` (listing de artigos)
- `/blog/[slug]` (artigo individual)
- `/blog/category/[category]` (categorização por tema)
- SEO meta tags e structured data (JSON-LD)

### Componentes/Features:
- Blog CMS integration (Contentful, Strapi, ou Markdown-based)
- SEO optimization engine
- Content analytics (Google Analytics 4)
- Lead capture forms (newsletter, gated content)
- Related articles recommendation
- Social sharing buttons

### Success Metrics:
- Organic traffic: 5k visitors/month (6 meses após launch)
- Blog → Simulator conversion: >15%
- Top 10 ranking para 20+ keywords de export intelligence
- Avg. time on page: >2 min

### Effort Estimate:
12 dias de desenvolvimento (blog engine + 20 artigos iniciais + SEO setup)

---

## J-A02: Campanhas Pagas (Google Ads / LinkedIn)
**Persona:** PME Exportador (ativo)
**User Story:** Como exportador que já vende internacionalmente, quero anúncios relevantes sobre ferramentas de export intelligence, para otimizar minhas operações
**Prioridade:** P2 - Medium
**Status:** ❌ Não iniciado
**Dependências:** J-A04 (Landing Pages Segmentadas)

### Fluxo:
1. Usuário vê anúncio no LinkedIn: "Descubra os destinos mais lucrativos para sua exportação"
2. Clica e é direcionado para landing page específica: `/lp/simulator-sugar` (se anúncio é para setor açucareiro)
3. Landing page mostra value prop focado no setor
4. CTA: "Simular agora grátis" (sem necessidade de cadastro inicial)
5. Após simulação, modal de conversão: "Receba análise completa por email"
6. Captura email + nome + NCM principal

### Páginas/Rotas Necessárias:
- `/lp/[campaign-slug]` (landing pages dinâmicas)
- `/lp/simulator-sugar`, `/lp/simulator-meat`, `/lp/simulator-soy`, etc.
- UTM parameter tracking (Google Analytics)

### Componentes/Features:
- Landing page builder (com A/B testing)
- Form validation e lead capture
- UTM parameter extraction e storage
- Conversion pixel (Google Ads, LinkedIn Insight Tag)
- Dynamic content injection (setor-specific messaging)
- Exit-intent modal (quando usuário tenta sair)

### Success Metrics:
- Cost per Lead (CPL): <R$50
- Landing page conversion rate: >8%
- CAC payback period: <6 meses
- ROAS (Return on Ad Spend): >3x

### Effort Estimate:
8 dias de desenvolvimento (landing page engine + tracking + A/B testing framework)

---

## J-A03: Indicação Direta (Boca-a-Boca)
**Persona:** PME Exportador (referenciado por colega)
**User Story:** Como exportador que ouviu falar do BGC por um colega, quero acessar rapidamente e ver se funciona para meu produto, para validar a recomendação
**Prioridade:** P1 - High
**Status:** ❌ Não iniciado
**Dependências:** J-R01 (Programa de Referral)

### Fluxo:
1. Usuário recebe link de indicação: `bgc.com.br/simulator?ref=ABC123`
2. Acessa homepage com banner: "Você foi convidado por [Nome do Amigo]. Ganhe 10 simulações grátis!"
3. Cria conta (ou usa como anônimo)
4. Se criar conta via referral, tanto indicador quanto indicado ganham bônus
5. Realiza primeiras simulações com limite expandido (10 ao invés de 5)

### Páginas/Rotas Necessárias:
- `/simulator?ref=[code]` (com lógica de referral tracking)
- `/invite/[friend-code]` (página específica de convite)

### Componentes/Features:
- Referral code generation e tracking
- Bonus credit system (extra simulações)
- Referral attribution logic (cookies + URL params)
- Social proof display ("João Silva convidou você")
- Gamification (badges para quem indica 5+, 10+, 50+ pessoas)

### Success Metrics:
- Virality coefficient (K-factor): >0.5 (cada usuário traz 0.5 novos)
- Referral conversion rate: >25%
- Referred user retention: >1.3x vs organic

### Effort Estimate:
6 dias de desenvolvimento (referral engine + tracking + rewards)

---

## J-A04: Parcerias com Associações de Exportadores
**Persona:** PME Exportador (membro de associação)
**User Story:** Como membro da APEX ou ABIEC, quero acessar ferramenta recomendada pela associação, para confiar que é legítima e útil
**Prioridade:** P2 - Medium
**Status:** ❌ Não iniciado
**Dependências:** Nenhuma (mas requer BD externo com parcerias)

### Fluxo:
1. Associação envia email blast para membros: "Parceria exclusiva: BGC oferece 30 dias premium grátis"
2. Usuário clica em link: `bgc.com.br/partners/apex`
3. Landing page co-branded (logo da APEX + BGC)
4. CTA: "Ativar acesso premium gratuito"
5. Cadastro simplificado (email corporativo obrigatório @associacao.com.br)
6. Onboarding customizado para membros da associação

### Páginas/Rotas Necessárias:
- `/partners/[partner-slug]` (landing pages co-branded)
- `/partners/apex`, `/partners/abiec`, `/partners/abimaq`, etc.

### Componentes/Features:
- Partner portal (dashboard para associações verem métricas de adoção)
- SSO integration (login via credenciais da associação)
- Co-branded email templates
- Partner-specific promotions engine
- Bulk user provisioning (importar lista de membros)

### Success Metrics:
- Partner accounts activated: 5 parcerias no Q1 2026
- Members onboarded: >20% dos membros da associação
- Premium conversion (após trial): >12%

### Effort Estimate:
10 dias de desenvolvimento (partner portal + SSO + co-branding system)

---

## J-A05: Webinars e Eventos Presenciais
**Persona:** PME Exportador (participa de eventos do setor)
**User Story:** Como participante de feira de exportação, quero testar ferramenta apresentada em palestra, para aplicar o que aprendi
**Prioridade:** P3 - Low
**Status:** ❌ Não iniciado
**Dependências:** J-A04 (parcerias), marketing collateral

### Fluxo:
1. BGC apresenta palestra em feira APAS Show ou ExpoFoodBR
2. Palestrante faz demo ao vivo do simulador
3. QR Code exibido no slide: "Teste agora e ganhe acesso premium 7 dias"
4. Participante escaneia QR Code → redireciona para `/events/apas2026`
5. Landing page com vídeo do palestrante
6. Cadastro rápido (email + telefone)
7. Acesso premium temporário ativado automaticamente

### Páginas/Rotas Necessárias:
- `/events/[event-slug]` (landing pages específicas por evento)
- `/events/apas2026`, `/events/expofood2026`, etc.

### Componentes/Features:
- Event tracking system (QR code generation)
- Time-limited premium access (expire após 7 dias)
- Event-specific onboarding flow
- Webinar replay hosting (Vimeo/YouTube embed)
- Lead scoring (participantes de eventos têm score +10)

### Success Metrics:
- Events attended: 4/trimestre
- Leads captured per event: >50
- Event → Paid conversion: >8%

### Effort Estimate:
5 dias de desenvolvimento (event landing pages + QR tracking + temporary access)

---

# 2. ACTIVATION - Primeira Experiência de Valor

## J-AC01: Onboarding de Novo Usuário (First-Time Experience)
**Persona:** PME Exportador (primeiro acesso)
**User Story:** Como novo usuário, quero entender rapidamente como usar a plataforma, para ter minha primeira experiência de valor em <5 minutos
**Prioridade:** **P0 - Critical (MVP Beta Blocker)**
**Status:** ❌ Não iniciado
**Dependências:** J-AC02 (Sistema de Autenticação)

### Fluxo:
1. Usuário cria conta (email + senha ou Google OAuth)
2. Tela de boas-vindas: "Bem-vindo ao BGC! Vamos configurar sua conta em 3 passos"
3. **Passo 1/3:** "Qual produto você exporta?" → autocomplete de NCMs
4. **Passo 2/3:** "Qual volume médio mensal?" → input numérico (kg)
5. **Passo 3/3:** "Para quais regiões já exporta?" → multi-select de continentes/países
6. Ao finalizar: "Preparamos uma análise inicial para você" → redireciona para dashboard com simulação pré-carregada
7. Overlay tutorial (tooltips interativos) destacando features principais

### Páginas/Rotas Necessárias:
- `/onboarding/welcome` (página de boas-vindas)
- `/onboarding/step-1`, `/onboarding/step-2`, `/onboarding/step-3`
- `/dashboard` (primeira tela após onboarding)

### Componentes/Features:
- Multi-step wizard (stepper component MUI)
- NCM autocomplete com fuzzy search
- Progress indicator (3 steps visualizados)
- Skip option (usuário pode pular e ir direto ao dashboard)
- Interactive tooltips (Joyride ou Shepherd.js)
- Onboarding analytics (track step abandonment)
- Pre-filled demo simulation (baseada em NCM escolhido)

### Success Metrics:
- Onboarding completion rate: >70%
- Time-to-first-value: <5 minutos
- First simulation within 24h: >60%
- Tooltip engagement: >40%

### Effort Estimate:
8 dias de desenvolvimento (wizard + tooltips + analytics + UX refinement)

---

## J-AC02: Sistema de Autenticação e Perfil
**Persona:** Qualquer usuário
**User Story:** Como usuário, quero criar conta e fazer login seguro, para salvar minhas simulações e preferências
**Prioridade:** **P0 - Critical (MVP Beta Blocker)**
**Status:** ❌ Não iniciado
**Dependências:** Nenhuma (base do sistema)

### Fluxo:
1. **Signup:**
   - Usuário clica "Criar Conta" na homepage
   - Formulário: Email, Senha (min 8 chars), Nome da Empresa, CNPJ (opcional)
   - Validação em tempo real (email único, senha forte)
   - Envia email de confirmação com link de verificação
   - Após verificação, redireciona para J-AC01 (onboarding)

2. **Login:**
   - Email + senha OU Google OAuth
   - "Esqueci minha senha" → reset via email
   - 2FA opcional (SMS ou Authenticator app)
   - Sessão expira após 7 dias de inatividade

3. **Perfil:**
   - Editar dados da empresa (nome, CNPJ, endereço, setor)
   - Avatar upload (logo da empresa)
   - Preferências (idioma, moeda, notificações)

### Páginas/Rotas Necessárias:
- `/signup` (cadastro)
- `/login` (autenticação)
- `/forgot-password` (recuperação)
- `/reset-password/[token]` (nova senha)
- `/verify-email/[token]` (confirmação de email)
- `/profile` (edição de perfil)
- `/profile/security` (senha, 2FA)

### Componentes/Features:
- JWT-based authentication (httpOnly cookies)
- OAuth2 integration (Google, LinkedIn)
- Email service (SendGrid ou AWS SES)
- Password hashing (bcrypt, Argon2)
- Rate limiting (prevenir brute force)
- CAPTCHA (signup e login após 3 tentativas falhas)
- Session management (Redis store)
- User roles (admin, premium, free)

### Success Metrics:
- Signup conversion (homepage → account): >12%
- Email verification rate: >85%
- OAuth adoption: >30% dos signups
- Login success rate: >98%

### Effort Estimate:
12 dias de desenvolvimento (auth engine + OAuth + email flows + security hardening)

---

## J-AC03: Primeiro Uso do Simulador (Logged-In User)
**Persona:** PME Exportador (acabou de criar conta)
**User Story:** Como usuário logado, quero simular destinos sem limite do freemium, para explorar livremente a plataforma
**Prioridade:** **P0 - Critical**
**Status:** 🔨 Parcial (simulador existe, mas integração com auth falta)
**Dependências:** J-AC02 (autenticação)

### Fluxo:
1. Usuário acessa `/simulator` (já logado)
2. Banner no topo: "Você tem simulações ilimitadas! (Plano Free)" ← não há limite para usuários logados no free tier (ajuste estratégico)
3. Formulário pré-preenchido com NCM do onboarding (se disponível)
4. Ao submeter, vê resultados + botão "Salvar esta simulação"
5. Modal: "Simulação salva! Acesse a qualquer momento em Histórico"
6. Tooltip: "Dica: Clique em qualquer destino para ver detalhes completos"

### Páginas/Rotas Necessárias:
- `/simulator` (mesma página, mas comportamento diferente se autenticado)
- `/simulator/history` (lista de simulações salvas)
- `/simulator/saved/[id]` (detalhes de simulação salva)

### Componentes/Features:
- Authentication check (redirect se não logado)
- Save simulation button (POST /v1/user/simulations)
- Simulation history persistence (PostgreSQL user_simulations table)
- Pre-fill form from onboarding data
- Enhanced results view (gráficos, comparação lado-a-lado)

### Success Metrics:
- Simulations per user (first 7 days): >3
- Save rate: >40%
- Return to history: >25%

### Effort Estimate:
6 dias de desenvolvimento (auth integration + save functionality + history UI)

---

## J-AC04: Exploração do Dashboard de Market Intelligence
**Persona:** PME Exportador (ativado)
**User Story:** Como exportador, quero ver insights sobre meu mercado (TAM, tendências, preços), para tomar decisões estratégicas
**Prioridade:** P1 - High
**Status:** ❌ Não iniciado
**Dependências:** J-AC02, J-AC03

### Fluxo:
1. Usuário logado acessa `/dashboard`
2. Dashboard exibe 4 cards principais:
   - **Card 1:** "Seu Mercado Total" → TAM/SAM/SOM baseado no NCM principal
   - **Card 2:** "Tendências de Preço" → gráfico de linha dos últimos 12 meses
   - **Card 3:** "Destinos Recomendados" → top 3 países do simulador
   - **Card 4:** "Próximas Ações" → tasks (ex: "Complete seu perfil para unlock insights avançados")
3. Cada card é clicável e expande para página detalhada

### Páginas/Rotas Necessárias:
- `/dashboard` (visão geral)
- `/intelligence/market-size` (TAM/SAM/SOM detalhado)
- `/intelligence/price-trends` (análise temporal de preços)
- `/intelligence/competitors` (quem mais exporta seu NCM)

### Componentes/Features:
- Dashboard grid layout (responsive)
- Chart components (Recharts: line, bar, pie)
- Data fetching hooks (SWR ou React Query)
- Skeleton loading states
- Empty states (quando usuário não tem NCM configurado)
- Export to PDF/Excel (dashboard completo)

### Success Metrics:
- Dashboard visit rate (first week): >50%
- Time on dashboard: >3 min
- Chart interaction rate: >30%
- Export report usage: >10%

### Effort Estimate:
10 dias de desenvolvimento (dashboard layout + charts + API integration + export feature)

---

## J-AC05: Configuração de Alertas e Notificações
**Persona:** PME Exportador (engajado)
**User Story:** Como exportador, quero receber alertas quando houver mudanças no mercado (preços, demanda, novos destinos), para não perder oportunidades
**Prioridade:** P2 - Medium
**Status:** ❌ Não iniciado
**Dependências:** J-AC02

### Fluxo:
1. Usuário acessa `/profile/notifications`
2. Configurações disponíveis:
   - **Email Alerts:** Diário, Semanal, Mensal, Nunca
   - **Alert Types:** Mudanças de preço (>10%), Novos destinos emergentes, Mudanças regulatórias
   - **Channels:** Email, SMS, Push (mobile app futuro)
3. Usuário seleciona preferências e salva
4. Sistema envia email de confirmação: "Alertas configurados com sucesso"
5. Primeiro alerta é enviado no próximo ciclo (ex: segunda-feira se escolheu semanal)

### Páginas/Rotas Necessárias:
- `/profile/notifications` (configuração)
- `/alerts/history` (histórico de alertas enviados)

### Componentes/Features:
- Notification preferences UI (toggles, frequency selectors)
- Email templating engine (Handlebars ou React Email)
- Cron jobs (alertas diários: 8am, semanais: segunda 9am)
- Alert generation logic (detect price changes, new opportunities)
- Unsubscribe mechanism (one-click unsubscribe link)
- Alert analytics (open rate, click rate)

### Success Metrics:
- Alert opt-in rate: >40%
- Email open rate: >25%
- Click-through rate: >8%
- Unsubscribe rate: <5%

### Effort Estimate:
8 dias de desenvolvimento (notification system + email engine + cron jobs + analytics)

---

# 3. RETENTION - Usuários Retornam

## J-R01: Histórico de Simulações e Comparação
**Persona:** PME Exportador (retornando)
**User Story:** Como exportador que já fez várias simulações, quero comparar resultados lado-a-lado, para decidir qual mercado priorizar
**Prioridade:** P1 - High
**Status:** ❌ Não iniciado
**Dependências:** J-AC03 (simulações salvas)

### Fluxo:
1. Usuário acessa `/simulator/history`
2. Vê lista de todas simulações salvas (ordenadas por data, mais recente primeiro)
3. Cada item mostra: NCM, Volume, Data, Top 3 Destinos, Score Médio
4. Checkbox em cada item: "Selecionar para comparar"
5. Usuário seleciona 2-4 simulações e clica "Comparar Selecionadas"
6. Redireciona para `/simulator/compare?ids=1,2,3`
7. Página mostra tabela comparativa com métricas lado-a-lado
8. Gráficos de radar comparando scores por destino

### Páginas/Rotas Necessárias:
- `/simulator/history` (lista)
- `/simulator/compare` (comparação lado-a-lado)
- `/simulator/saved/[id]` (detalhes de simulação individual)

### Componentes/Features:
- Simulation list (tabela ou cards)
- Multi-select checkbox system
- Comparison table (sortable columns)
- Radar chart (overlaying múltiplas simulações)
- Export comparison (PDF com todas simulações)
- Notes/annotations (usuário pode adicionar observações em cada simulação)

### Success Metrics:
- Return to history (D7): >40%
- Comparison usage: >15% dos usuários
- Avg simulations compared: 2.5
- Export comparison report: >8%

### Effort Estimate:
7 dias de desenvolvimento (history UI + comparison logic + charts + export)

---

## J-R02: Watchlist de Destinos Favoritos
**Persona:** PME Exportador (focado)
**User Story:** Como exportador focando em 3-5 mercados estratégicos, quero monitorar apenas esses destinos, para receber atualizações relevantes
**Prioridade:** P2 - Medium
**Status:** ❌ Não iniciado
**Dependências:** J-AC03

### Fluxo:
1. Usuário vê resultados do simulador
2. Em cada card de destino, botão "Adicionar ao Watchlist" (estrela)
3. Ao clicar, destino é salvo em `/watchlist`
4. Página `/watchlist` mostra dashboard personalizado:
   - Gráficos de tendência para cada país no watchlist
   - Alertas específicos (ex: "China aumentou importações de açúcar em 12% este mês")
   - Quick actions (simular novamente, ver histórico, remover)
5. Email semanal: "Atualizações dos seus mercados favoritos"

### Páginas/Rotas Necessárias:
- `/watchlist` (dashboard de destinos favoritados)
- `/watchlist/[country-code]` (detalhes de país específico)

### Componentes/Features:
- Favorite/star toggle button
- Watchlist persistence (user_watchlist table)
- Country-specific trends (time series data)
- Weekly digest email (curated insights)
- Remove from watchlist (undo mechanism)
- Watchlist limit (free: 3 países, premium: ilimitado)

### Success Metrics:
- Watchlist adoption: >30%
- Avg destinations in watchlist: 2.8
- Weekly email engagement: >18%
- Return to watchlist (D30): >50%

### Effort Estimate:
6 dias de desenvolvimento (watchlist system + country trends + email digest)

---

## J-R03: Exploração de Novos NCMs (Cross-Sell)
**Persona:** PME Exportador (diversificando portfólio)
**User Story:** Como exportador que quer diversificar produtos, quero descobrir NCMs relacionados ao meu atual, para expandir operações
**Prioridade:** P2 - Medium
**Status:** ❌ Não iniciado
**Dependências:** J-AC04

### Fluxo:
1. Usuário acessa dashboard ou simulador
2. Widget lateral: "Produtos Relacionados ao Seu NCM"
3. Sistema recomenda NCMs similares (mesma seção/capítulo ou mesmo setor)
4. Ex: Se usuário exporta 17011400 (açúcar bruto), recomenda 17019900 (outros açúcares)
5. Ao clicar em NCM recomendado, pre-fill simulador com novo NCM
6. Usuário pode simular e comparar com NCM principal

### Páginas/Rotas Necessárias:
- `/intelligence/related-products` (exploração de NCMs similares)
- `/intelligence/ncm/[code]` (deep-dive em NCM específico)

### Componentes/Features:
- NCM recommendation algorithm (collaborative filtering ou rule-based)
- Related products widget (sidebar ou modal)
- NCM deep-dive page (história, mercados principais, concorrentes)
- Add to comparison (comparar NCM atual vs recomendado)

### Success Metrics:
- Related NCM exploration: >12%
- New NCM simulation rate: >8%
- Portfolio diversification (users simulating 2+ NCMs): >20%

### Effort Estimate:
8 dias de desenvolvimento (recommendation engine + NCM deep-dive pages + comparison)

---

## J-R04: Análise de Concorrentes (Quem Mais Exporta)
**Persona:** PME Exportador (competitivo)
**User Story:** Como exportador, quero saber quem são meus concorrentes e quanto eles exportam, para benchmarking e estratégia
**Prioridade:** P2 - Medium
**Status:** ❌ Não iniciado
**Dependências:** J-AC02, dados de exportadores (requer integração MDIC ou scraping)

### Fluxo:
1. Usuário acessa `/intelligence/competitors`
2. Sistema mostra top exportadores do mesmo NCM (dados agregados, sem CNPJ)
3. Exibe: Ranking, Estado de origem, Volume médio, Destinos principais
4. Se usuário configurou CNPJ no perfil, mostra "Sua posição estimada: Top 15%"
5. Gráficos: Market share por estado, Evolução temporal de concorrência

### Páginas/Rotas Necessárias:
- `/intelligence/competitors` (visão geral)
- `/intelligence/competitors/[ncm]` (concorrentes por NCM)

### Componentes/Features:
- Competitor ranking table (sortable)
- Market share pie chart (por estado/região)
- Time series (evolução de concorrência)
- User positioning (se CNPJ fornecido, cross-reference com dados MDIC)
- Privacy-preserving aggregation (não expor CNPJs de concorrentes)

### Success Metrics:
- Competitor page visit: >20%
- Time on page: >2 min
- CNPJ provision (to unlock positioning): >15%

### Effort Estimate:
10 dias de desenvolvimento (data pipeline MDIC + competitor analysis + privacy layer)

---

## J-R05: Participação em Comunidade/Fórum
**Persona:** PME Exportador (quer network)
**User Story:** Como exportador, quero trocar experiências com outros exportadores, para aprender melhores práticas e resolver dúvidas
**Prioridade:** P3 - Low
**Status:** ❌ Não iniciado
**Dependências:** J-AC02 (autenticação)

### Fluxo:
1. Usuário acessa `/community`
2. Vê fórum estilo Stack Overflow:
   - Perguntas recentes
   - Topics: Logística, Documentação, Negociação, Mercados Específicos
3. Pode fazer pergunta, responder, upvote, salvar tópicos
4. Sistema de reputação (gamification): Pontos por respostas úteis, badges
5. Moderação por admins BGC + usuários premium (trusted members)

### Páginas/Rotas Necessárias:
- `/community` (homepage do fórum)
- `/community/[topic]` (tópico específico)
- `/community/ask` (criar nova pergunta)
- `/community/users/[id]` (perfil de usuário no fórum)

### Componentes/Features:
- Forum engine (Discourse integration ou custom build)
- Voting system (upvote/downvote)
- Moderation tools (flag, hide, ban)
- Reputation/badge system
- Email notifications (new answers, mentions)
- Search (full-text search em perguntas/respostas)

### Success Metrics:
- Community activation: >10% dos usuários
- Questions asked per month: >50
- Response rate: >60%
- Avg time to first answer: <24h

### Effort Estimate:
15 dias de desenvolvimento (forum engine + moderation + gamification + integration)

---

## J-R06: Email Marketing e Re-Engagement
**Persona:** PME Exportador (inativo há 30 dias)
**User Story:** Como usuário que não voltou à plataforma, quero receber lembretes personalizados, para re-engajar se houver valor
**Prioridade:** P1 - High
**Status:** ❌ Não iniciado
**Dependências:** J-AC05 (sistema de notificações)

### Fluxo:
1. Sistema detecta usuário inativo há 15 dias (sem login)
2. **Email 1 (D15):** "Sentimos sua falta! Veja o que mudou no mercado [NCM]"
   - Destaca insights novos (ex: preço de açúcar subiu 8%)
   - CTA: "Ver análise completa"
3. **Email 2 (D30):** "Seus concorrentes estão se movendo. E você?"
   - Social proof (X empresas usaram simulador esta semana)
   - CTA: "Simular novo destino"
4. **Email 3 (D45):** "Última chance: Oferta especial 20% off premium"
   - Desconto por tempo limitado
   - CTA: "Ativar desconto"
5. Se não responde a nenhum email, pausa sequência

### Páginas/Rotas Necessárias:
- Nenhuma nova (links direcionam para dashboard/simulator)

### Componentes/Features:
- Drip email campaign engine (Mailchimp, SendGrid, ou custom)
- User activity tracking (last_login_at)
- Personalized content injection (NCM, nome, empresa)
- A/B testing de subject lines e CTAs
- Unsubscribe/preference center
- Re-engagement scoring (opens, clicks → re-activated user)

### Success Metrics:
- Email open rate: >20%
- Click-through rate: >6%
- Re-activation rate (D45): >8%
- Unsubscribe rate: <3%

### Effort Estimate:
7 dias de desenvolvimento (drip campaign setup + personalization + A/B testing)

---

# 4. REVENUE - Monetização

## J-REV01: Descoberta do Plano Premium (Paywall Soft)
**Persona:** PME Exportador (free user, high engagement)
**User Story:** Como usuário free que usa muito a plataforma, quero entender benefícios do premium, para avaliar se vale upgrade
**Prioridade:** **P0 - Critical (Monetization Enabler)**
**Status:** ❌ Não iniciado
**Dependências:** J-AC02

### Fluxo:
1. Usuário free tenta acessar feature premium (ex: "Análise de Tarifas Detalhada")
2. Modal aparece: "Recurso Premium"
   - Lista benefícios: Simulações ilimitadas, Relatórios em PDF, Alertas avançados, Suporte prioritário
   - Preço destacado: R$ 197/mês ou R$ 1.970/ano (economize 17%)
   - CTA: "Iniciar trial 14 dias grátis"
3. Se usuário fecha modal, banner fixo no topo: "Upgrade para Premium e desbloqueie insights avançados"
4. Botão "Ver Planos" sempre visível no menu principal

### Páginas/Rotas Necessárias:
- `/pricing` (tabela comparativa Free vs Premium vs Enterprise)
- `/upgrade` (página dedicada de upgrade com benefícios detalhados)

### Componentes/Features:
- Paywall modal (soft block em features premium)
- Pricing table (comparison: Free vs Premium vs Enterprise)
- Feature toggles (backend flag `user.tier`)
- Trial activation flow (auto-downgrade após 14 dias se não converter)
- Pricing calculator (anual vs mensal, economia visualizada)
- Social proof (depoimentos de clientes premium)

### Success Metrics:
- Paywall modal impression → trial start: >5%
- Pricing page visit rate: >25% dos free users
- Trial activation rate: >10%
- Trial → Paid conversion: >25%

### Effort Estimate:
6 dias de desenvolvimento (paywall system + pricing page + trial logic)

---

## J-REV02: Processo de Checkout e Pagamento
**Persona:** PME Exportador (decidiu assinar premium)
**User Story:** Como exportador pronto para pagar, quero checkout rápido e seguro, para começar a usar premium imediatamente
**Prioridade:** **P0 - Critical**
**Status:** ❌ Não iniciado
**Dependências:** J-REV01

### Fluxo:
1. Usuário clica "Assinar Premium" na pricing page
2. Redireciona para `/checkout`
3. **Step 1:** Escolher plano (Mensal R$197 ou Anual R$1.970)
4. **Step 2:** Informações de faturamento (Nome, CNPJ, Endereço)
5. **Step 3:** Método de pagamento:
   - Cartão de crédito (Stripe/iugu)
   - Boleto bancário (processamento em 1-3 dias)
   - PIX (aprovação instantânea)
6. **Step 4:** Confirmação e emissão de nota fiscal (NFSe automática)
7. Após pagamento, email de boas-vindas + badge "Premium" no perfil
8. Acesso imediato a todas features premium

### Páginas/Rotas Necessárias:
- `/checkout` (processo de checkout)
- `/checkout/success` (confirmação pós-pagamento)
- `/checkout/failed` (erro no pagamento, retry)
- `/billing` (gerenciar assinatura, faturas, método de pagamento)

### Componentes/Features:
- Stripe integration (ou iugu para Brasil-first)
- Payment method selection (card, boleto, PIX)
- Subscription management (create, upgrade, downgrade, cancel)
- Invoice generation (NFSe via API de prefeituras)
- Receipt emails (transacional)
- Failed payment retry logic (dunning management)
- PCI compliance (never store card data, use Stripe.js)

### Success Metrics:
- Checkout abandonment rate: <30%
- Payment success rate: >95%
- PIX adoption: >40% dos pagamentos
- Avg time to complete checkout: <3 min

### Effort Estimate:
12 dias de desenvolvimento (Stripe integration + checkout flow + billing dashboard + NFSe)

---

## J-REV03: Gerenciamento de Assinatura e Billing
**Persona:** PME Exportador (cliente premium)
**User Story:** Como assinante premium, quero gerenciar minha assinatura (cancelar, trocar plano, atualizar cartão), para ter controle total
**Prioridade:** P1 - High
**Status:** ❌ Não iniciado
**Dependências:** J-REV02

### Fluxo:
1. Usuário premium acessa `/billing`
2. Dashboard mostra:
   - Plano atual (Premium Mensal)
   - Próxima cobrança (15/02/2026 - R$197)
   - Método de pagamento (Mastercard •••• 1234)
   - Histórico de faturas (download PDF/NFSe)
3. Ações disponíveis:
   - "Trocar para plano anual" (oferta: economize 17%)
   - "Atualizar método de pagamento"
   - "Cancelar assinatura" (retention flow)
4. Se clicar cancelar:
   - Modal: "Sentiremos sua falta. Podemos oferecer 20% off por 3 meses?"
   - Opções: Aceitar oferta / Pausar 30 dias / Confirmar cancelamento
5. Se confirmar cancelamento, downgrade para free ao final do período pago

### Páginas/Rotas Necessárias:
- `/billing` (dashboard)
- `/billing/invoices` (histórico)
- `/billing/payment-method` (atualizar cartão)
- `/billing/cancel` (fluxo de cancelamento com retention offers)

### Componentes/Features:
- Subscription dashboard (current plan, next billing, usage stats)
- Plan upgrade/downgrade (prorated billing)
- Payment method update (Stripe Customer Portal)
- Invoice history (list, download PDF)
- Cancellation flow (multi-step retention)
- Pause subscription (freeze account por 30 dias)
- Reactivation (win-back campaigns)

### Success Metrics:
- Churn rate: <5%/mês
- Retention offer acceptance: >30%
- Plan upgrade (monthly → annual): >20%
- Payment method update success: >90%

### Effort Estimate:
8 dias de desenvolvimento (billing dashboard + retention flow + Stripe Customer Portal integration)

---

## J-REV04: Plano Enterprise (Custom Pricing)
**Persona:** Grande Exportador ou Trading Company
**User Story:** Como empresa com >R$50M de faturamento, quero plano customizado com API access, white-label, e suporte dedicado
**Prioridade:** P2 - Medium
**Status:** ❌ Não iniciado
**Dependências:** J-REV01, J-REV02

### Fluxo:
1. Usuário vê pricing page, último tier: "Enterprise - Sob Consulta"
2. Clica "Falar com Vendas"
3. Formulário de contato:
   - Nome, Email, Empresa, Faturamento anual, Necessidades específicas
4. Lead é enviado para CRM (Salesforce/Pipedrive)
5. SDR (Sales Development Rep) entra em contato em <24h
6. Demo personalizada agendada
7. Proposta customizada gerada (pricing, SLA, features exclusivas)
8. Contrato assinado (DocuSign)
9. Onboarding white-glove (Customer Success Manager dedicado)

### Páginas/Rotas Necessárias:
- `/pricing` (tier Enterprise visível)
- `/enterprise/contact` (formulário de contato)
- `/enterprise/demo` (agendamento de demo)

### Componentes/Features:
- Lead capture form (qualificação de leads)
- CRM integration (Salesforce, Pipedrive, HubSpot)
- Calendar integration (Calendly para agendar demos)
- Contract generation (DocuSign API)
- White-label customization (logo, cores, domínio customizado)
- API access (docs, keys, rate limits customizados)
- Dedicated support (Slack Connect, CSM assignment)

### Success Metrics:
- Enterprise leads per month: >5
- Demo booking rate: >60%
- Enterprise conversion rate: >15%
- Avg deal size: R$15k/mês

### Effort Estimate:
10 dias de desenvolvimento (enterprise flow + CRM integration + white-label engine + API docs)

---

## J-REV05: Upsell de Add-Ons (Relatórios Custom, API)
**Persona:** PME Exportador (premium, power user)
**User Story:** Como usuário premium que precisa de mais, quero contratar add-ons específicos, para customizar minha experiência
**Prioridade:** P3 - Low
**Status:** ❌ Não iniciado
**Dependências:** J-REV02

### Fluxo:
1. Usuário premium acessa `/billing/add-ons`
2. Vê marketplace de add-ons:
   - **Relatórios Custom:** R$99/relatório (análise sob-medida por analista BGC)
   - **API Access:** R$297/mês (10k requests/mês)
   - **White-Label Dashboard:** R$497/mês (embbed BGC no seu site)
   - **Priority Support:** R$197/mês (SLA 2h response time)
3. Seleciona add-on e adiciona à assinatura
4. Cobrança adicional aparece na próxima fatura (prorated)

### Páginas/Rotas Necessárias:
- `/billing/add-ons` (marketplace)
- `/billing/add-ons/[addon-id]` (detalhes de add-on específico)

### Componentes/Features:
- Add-on marketplace (catalog)
- Add-on purchase flow (add to subscription)
- Prorated billing (adjust next invoice)
- Add-on management (enable/disable)
- Usage tracking (API requests, reports generated)

### Success Metrics:
- Add-on adoption: >15% dos premium users
- Most popular add-on: API Access (hypothesis)
- Avg add-on revenue per premium user: R$150/mês

### Effort Estimate:
7 dias de desenvolvimento (add-on marketplace + billing logic + usage tracking)

---

# 5. REFERRAL - Network Effects

## J-REF01: Programa de Indicação (Refer-a-Friend)
**Persona:** PME Exportador (satisfeito)
**User Story:** Como usuário satisfeito, quero indicar colegas e ganhar benefícios, para ajudá-los e ser recompensado
**Prioridade:** P1 - High
**Status:** ❌ Não iniciado
**Dependências:** J-AC02

### Fluxo:
1. Usuário acessa `/referral`
2. Vê seu link único: `bgc.com.br/invite/ABC123`
3. Dashboard mostra:
   - Indicações enviadas: 5
   - Indicações convertidas: 2
   - Recompensas ganhas: R$200 em créditos
4. Opções de compartilhamento:
   - WhatsApp (pre-filled message)
   - Email (email template)
   - LinkedIn (social post)
   - Copiar link
5. Quando amigo se cadastra via link:
   - Indicador ganha: 1 mês grátis de premium
   - Indicado ganha: 30 simulações extras no free tier
6. Se indicado assina premium, indicador ganha R$100 de crédito (50% da primeira mensalidade)

### Páginas/Rotas Necessárias:
- `/referral` (dashboard de indicações)
- `/invite/[code]` (landing page para indicados)

### Componentes/Features:
- Referral code generation (unique per user)
- Referral tracking (cookies + URL params + database)
- Reward distribution (credits, free months)
- Social sharing buttons (WhatsApp, email, LinkedIn)
- Referral leaderboard (gamification: top referrers)
- Fraud detection (prevent self-referral, fake accounts)

### Success Metrics:
- Referral program activation: >25% dos usuários
- Virality coefficient (K-factor): >0.6
- Referral → signup conversion: >30%
- Referral → paid conversion: >15%

### Effort Estimate:
8 dias de desenvolvimento (referral engine + tracking + rewards + social sharing)

---

## J-REF02: Integração com CRM/ERP (Network Lock-In)
**Persona:** Empresa Premium (quer automatizar)
**User Story:** Como empresa que usa Salesforce/SAP, quero integrar BGC ao meu CRM, para sincronizar leads e oportunidades automaticamente
**Prioridade:** P3 - Low
**Status:** ❌ Não iniciado
**Dependências:** J-REV04 (enterprise tier)

### Fluxo:
1. Cliente Enterprise solicita integração com Salesforce
2. BGC fornece app na Salesforce AppExchange (ou API credentials)
3. Admin instala app e configura mapeamento de campos
4. Sincronização automática:
   - Leads do simulador → Salesforce Leads
   - Simulações salvas → Salesforce Opportunities
   - Buyer matches → Salesforce Contacts
5. Webhooks bidirecionais (mudanças no SF refletem no BGC e vice-versa)

### Páginas/Rotas Necessárias:
- `/integrations` (marketplace de integrações)
- `/integrations/salesforce/setup` (wizard de configuração)

### Componentes/Features:
- OAuth2 integration (Salesforce, HubSpot, SAP)
- Field mapping UI (drag-and-drop)
- Webhook engine (bidirectional sync)
- Sync logs (audit trail de sincronizações)
- Error handling e retry logic
- API documentation (para custom integrations)

### Success Metrics:
- Enterprise clients using integrations: >50%
- Churn reduction (integrated clients): -30% vs non-integrated
- API uptime: 99.9%

### Effort Estimate:
15 dias de desenvolvimento (OAuth + field mapping + webhook engine + docs)

---

## J-REF03: Co-Marketing com Parceiros (Associações, Freight Forwarders)
**Persona:** Parceiro (Associação de Exportadores)
**User Story:** Como associação parceira, quero co-marketing com BGC, para oferecer valor aos membros e ter visibilidade
**Prioridade:** P2 - Medium
**Status:** ❌ Não iniciado
**Dependências:** J-A04 (parcerias)

### Fluxo:
1. BGC firma parceria com ABIEC (Associação Brasileira das Indústrias Exportadoras de Carnes)
2. Landing page co-branded: `/partners/abiec`
3. Membros ABIEC ganham 50% off premium por 6 meses
4. BGC fornece conteúdo exclusivo para newsletter da ABIEC
5. Webinar conjunto: "Como Diversificar Destinos de Exportação de Carne"
6. Métricas compartilhadas (dashboard para ABIEC ver adoção dos membros)

### Páginas/Rotas Necessárias:
- `/partners/[partner-slug]` (landing pages co-branded)
- `/partners/dashboard` (analytics dashboard para parceiros)

### Componentes/Features:
- Partner portal (login para parceiros verem métricas)
- Co-branded landing pages (logo, cores, messaging customizados)
- Discount code generation (partner-specific)
- Content syndication (BGC content → partner newsletter)
- Webinar platform integration (Zoom, Google Meet)
- Lead attribution (track conversions from partner)

### Success Metrics:
- Active partnerships: >10 (ano 1)
- Partner-referred users: >30% de todos signups
- Partner satisfaction (NPS): >70

### Effort Estimate:
10 dias de desenvolvimento (partner portal + co-branding + analytics dashboard)

---

## J-REF04: User-Generated Content (Case Studies, Reviews)
**Persona:** PME Exportador (case de sucesso)
**User Story:** Como exportador que teve sucesso usando BGC, quero compartilhar minha história, para ajudar outros e ter visibilidade
**Prioridade:** P3 - Low
**Status:** ❌ Não iniciado
**Dependências:** J-AC02

### Fluxo:
1. BGC identifica usuário premium com alta atividade e bons resultados
2. Outreach: "Gostaria de ser case study? Publicaremos no site e LinkedIn"
3. Usuário aceita e preenche formulário: Desafio, Solução com BGC, Resultados
4. BGC formata como case study profissional (com fotos, logo da empresa)
5. Publicação em `/case-studies/[empresa-slug]`
6. Divulgação em LinkedIn, newsletter, homepage
7. Empresa ganha badge "Featured Customer" no perfil
8. Opcional: vídeo-depoimento gravado e editado por BGC

### Páginas/Rotas Necessárias:
- `/case-studies` (lista de cases)
- `/case-studies/[slug]` (case individual)
- `/reviews` (página de reviews/depoimentos)

### Componentes/Features:
- Case study submission form
- Case study template (Notion ou Markdown-based)
- Video testimonial hosting (YouTube embed)
- Review/rating system (5 stars, texto aberto)
- Badge system (Featured Customer, Top Reviewer)
- Social sharing (LinkedIn, Twitter)

### Success Metrics:
- Case studies published: >12/ano
- Case study page visits: >1k/mês
- Conversion rate (visitor reads case → signup): >8%
- Review count: >50 (primeiro ano)

### Effort Estimate:
7 dias de desenvolvimento (case study system + review platform + badges)

---

## J-REF05: Programa de Afiliados (Influencers, Consultores)
**Persona:** Consultor de Comex (quer monetizar audiência)
**User Story:** Como consultor de comércio exterior com audiência, quero ser afiliado BGC, para gerar receita recorrente indicando a plataforma
**Prioridade:** P3 - Low
**Status:** ❌ Não iniciado
**Dependências:** J-REV02 (billing), J-REF01 (referral base)

### Fluxo:
1. Consultor se cadastra em `/affiliates`
2. Aprovação manual por equipe BGC (verificação de audiência/credibilidade)
3. Afiliado recebe:
   - Link rastreável: `bgc.com.br?aff=consultor123`
   - Banner kit (imagens, CTAs)
   - Comissão: 30% recorrente por 12 meses
4. Dashboard mostra: Cliques, Conversões, Comissões ganhas, Próximo pagamento
5. Pagamentos mensais via PIX/transferência (threshold: R$100 mínimo)
6. Afiliados top ganham bônus (ex: MacBook para quem trouxer 50+ clientes)

### Páginas/Rotas Necessárias:
- `/affiliates` (signup e info do programa)
- `/affiliates/dashboard` (métricas e comissões)
- `/affiliates/resources` (banner kit, email templates)

### Componentes/Features:
- Affiliate program management (applicant review, approval)
- Affiliate tracking (cookies, URL params, database attribution)
- Commission calculation engine (30% recorrente)
- Payout system (PIX integration, invoice generation)
- Affiliate dashboard (clicks, conversions, earnings)
- Marketing resources repository (banners, copy, email templates)

### Success Metrics:
- Active affiliates: >50 (ano 1)
- Affiliate-driven revenue: >20% de MRR
- Avg commission per affiliate: R$800/mês
- Top affiliate earnings: >R$5k/mês

### Effort Estimate:
12 dias de desenvolvimento (affiliate platform + tracking + commission engine + payouts)

---

# 6. ADMINISTRATIVE JOURNEYS - Operações Internas

## J-ADM01: Onboarding de Parceiros (Freight Forwarders, Despachantes)
**Persona:** Admin BGC
**User Story:** Como admin, quero cadastrar parceiros logísticos, para conectá-los com exportadores e gerar receita por transação
**Prioridade:** P2 - Medium
**Status:** ❌ Não iniciado
**Dependências:** Nenhuma (admin-facing)

### Fluxo:
1. Admin acessa `/admin/partners`
2. Clica "Adicionar Novo Parceiro"
3. Formulário:
   - Tipo (Freight Forwarder, Despachante, Seguradora, Banco)
   - Nome, CNPJ, Email, Telefone
   - Áreas de atuação (países/portos cobertos)
   - Termos de comissão (ex: 10% por transação fechada)
4. Parceiro recebe email de boas-vindas com credenciais de acesso
5. Parceiro completa perfil em `/partner-portal`
6. Admin aprova perfil → parceiro aparece em marketplace para exportadores

### Páginas/Rotas Necessárias:
- `/admin/partners` (CRUD de parceiros)
- `/partner-portal` (área do parceiro)
- `/partner-portal/leads` (leads recebidos de exportadores)

### Componentes/Features:
- Admin panel (CRUD: Create, Read, Update, Delete partners)
- Partner onboarding wizard
- Partner profile (serviços, pricing, avaliações)
- Lead distribution engine (exportador solicita cotação → distribui para 3 parceiros)
- Commission tracking (transações fechadas → comissão BGC)
- Partner ratings (exportadores avaliam parceiros)

### Success Metrics:
- Partners onboarded: >20 (ano 1)
- Partner response rate: >80%
- Avg partner rating: >4.2/5
- Commission revenue: >R$50k/ano

### Effort Estimate:
10 dias de desenvolvimento (admin panel + partner portal + lead distribution)

---

## J-ADM02: Moderação de Conteúdo e Usuários
**Persona:** Admin BGC
**User Story:** Como admin, quero moderar conteúdo (fórum, reviews), para manter qualidade e prevenir spam/abuso
**Prioridade:** P2 - Medium
**Status:** ❌ Não iniciado
**Dependências:** J-R05 (fórum), J-REF04 (reviews)

### Fluxo:
1. Usuário posta pergunta no fórum com conteúdo inapropriado
2. Sistema detecta palavras-chave suspeitas (auto-flagging)
3. Post entra em fila de moderação (`/admin/moderation`)
4. Admin revisa e decide:
   - Aprovar (publicar)
   - Editar (remover parte inapropriada)
   - Rejeitar (deletar)
   - Banir usuário (se recorrente)
5. Usuário recebe notificação da decisão
6. Dashboard mostra métricas: Posts flagged, False positives, Response time

### Páginas/Rotas Necessárias:
- `/admin/moderation` (fila de moderação)
- `/admin/moderation/reports` (relatórios de abuso)
- `/admin/users` (gestão de usuários: ban, suspend, reset password)

### Componentes/Features:
- Auto-flagging engine (keyword detection, spam patterns)
- Moderation queue (priority sorting: high-risk first)
- Moderation actions (approve, edit, reject, ban)
- User reporting system (flag button em posts/reviews)
- Ban/suspend logic (temporary vs permanent)
- Audit log (todas ações de moderação registradas)

### Success Metrics:
- Avg moderation response time: <2h
- False positive rate: <10%
- User reports resolved: >95%
- Banned users: <1% da base

### Effort Estimate:
8 dias de desenvolvimento (moderation queue + auto-flagging + actions + audit log)

---

## J-ADM03: Analytics e Business Intelligence
**Persona:** Admin BGC / Product Manager
**User Story:** Como PM, quero dashboards de métricas-chave, para tomar decisões data-driven sobre produto e negócio
**Prioridade:** P1 - High
**Status:** ❌ Não iniciado
**Dependências:** Nenhuma (mas requer instrumentação de eventos)

### Fluxo:
1. PM acessa `/admin/analytics`
2. Dashboards disponíveis:
   - **Growth:** MAU, DAU, Signups, Churn, Retention cohorts
   - **Engagement:** Simulations per user, Features usage, Time on platform
   - **Revenue:** MRR, ARR, ARPU, LTV, CAC, Payback period
   - **Product:** Funnel conversion (homepage → signup → premium), A/B tests
3. Filtros: Date range, User segment (free/premium), Acquisition channel
4. Export relatórios (CSV, PDF)
5. Scheduled reports (enviar por email toda segunda 9am)

### Páginas/Rotas Necessárias:
- `/admin/analytics` (dashboards overview)
- `/admin/analytics/growth`, `/admin/analytics/revenue`, etc.

### Componentes/Features:
- Analytics dashboards (Metabase, Superset, ou custom Recharts)
- Event tracking (Segment, Mixpanel, ou custom)
- Funnel analysis (visualização de drop-off)
- Cohort analysis (retention por signup month)
- A/B testing framework (Optimizely, GrowthBook, ou custom)
- Scheduled reports (cron jobs + email)

### Success Metrics:
- Dashboards used by team: >80% semanalmente
- Time to insight: <5 min
- Data-driven decisions: >50% de features priorizadas por dados

### Effort Estimate:
12 dias de desenvolvimento (dashboard setup + event tracking + scheduled reports)

---

## J-ADM04: Gestão de Configurações de Sistema
**Persona:** Admin BGC
**User Story:** Como admin, quero gerenciar configurações globais (rate limits, preços, features flags), para operar a plataforma sem deploy
**Prioridade:** P2 - Medium
**Status:** ❌ Não iniciado
**Dependências:** Nenhuma

### Fluxo:
1. Admin acessa `/admin/settings`
2. Configurações disponíveis:
   - **Rate Limits:** Free tier (5 → 10 simulations/day)
   - **Pricing:** Premium monthly (R$197 → R$247)
   - **Feature Flags:** Enable/disable features (ex: community forum, beta features)
   - **Email Templates:** Editar subject/body de emails transacionais
   - **Maintenance Mode:** Ativar página de manutenção
3. Mudanças têm efeito imediato (sem necessidade de deploy)
4. Audit log registra quem mudou o quê e quando

### Páginas/Rotas Necessárias:
- `/admin/settings` (configurações globais)
- `/admin/feature-flags` (toggles de features)
- `/admin/email-templates` (editor de templates)

### Componentes/Features:
- Settings management (key-value store em PostgreSQL)
- Feature flags system (LaunchDarkly-like, ou custom)
- Email template editor (WYSIWYG: TinyMCE ou Quill)
- Maintenance mode (static page served quando ativo)
- Audit log (track all config changes)
- Rollback mechanism (desfazer mudança se quebrou algo)

### Success Metrics:
- Config changes per month: ~10
- Incidents caused by config change: <1/trimestre
- Rollback usage: <5% das mudanças

### Effort Estimate:
8 dias de desenvolvimento (settings system + feature flags + editor + audit)

---

## J-ADM05: Suporte ao Cliente (Ticketing System)
**Persona:** Admin BGC (Customer Support)
**User Story:** Como agente de suporte, quero gerenciar tickets de clientes, para resolver problemas rapidamente e manter satisfação alta
**Prioridade:** P1 - High
**Status:** ❌ Não iniciado
**Dependências:** J-AC02 (usuários autenticados)

### Fluxo:
1. Usuário encontra problema e clica "Suporte" no menu
2. Formulário de ticket:
   - Categoria (Bug, Dúvida, Feedback, Outro)
   - Prioridade (Baixa, Média, Alta, Urgente)
   - Descrição
   - Anexos (screenshots)
3. Ticket criado e enviado para `/admin/support`
4. Agente de suporte pega ticket, responde via interface
5. Usuário recebe email com resposta
6. Thread continua (usuário pode responder, tudo fica no ticket)
7. Quando resolvido, agente marca "Fechado"
8. Usuário recebe CSAT survey: "Como foi o atendimento? 1-5 estrelas"

### Páginas/Rotas Necessárias:
- `/support` (abrir novo ticket)
- `/support/tickets` (meus tickets)
- `/support/tickets/[id]` (detalhes de ticket específico)
- `/admin/support` (fila de tickets para agentes)

### Componentes/Features:
- Ticketing system (Zendesk-like, ou custom)
- Email integration (replies via email viram comments no ticket)
- SLA tracking (response time, resolution time)
- Canned responses (templates de respostas comuns)
- Internal notes (agentes podem deixar notas invisíveis para usuário)
- CSAT survey (post-resolution)
- Knowledge base integration (sugerir artigos automaticamente)

### Success Metrics:
- Avg first response time: <4h (business hours)
- Avg resolution time: <24h
- CSAT score: >4.5/5
- Ticket volume: <5% dos MAU abrem ticket

### Effort Estimate:
10 dias de desenvolvimento (ticketing system + email integration + SLA tracking + CSAT)

---

## J-ADM06: Gestão de Dados (ETL, Backups, Data Quality)
**Persona:** Admin BGC (Data Engineer)
**User Story:** Como data engineer, quero monitorar pipelines de dados, para garantir que simulador sempre tem dados frescos e corretos
**Prioridade:** P1 - High
**Status:** ❌ Não iniciado
**Dependências:** Nenhuma (infraestrutura)

### Fluxo:
1. Engineer acessa `/admin/data-pipelines`
2. Vê status de jobs:
   - **ComexStat Ingestion:** Última execução 2026-01-20 03:00 ✅ Success (10M records)
   - **Countries Metadata Refresh:** Última execução 2026-01-19 04:00 ❌ Failed (API timeout)
   - **Materialized Views Refresh:** Última execução 2026-01-20 05:00 ✅ Success
3. Clica em job failed → vê logs de erro
4. Trigger manual re-run
5. Configurar alertas (Slack/Email se job falhar)
6. Data quality checks:
   - Verificar se não há gaps temporais (missing months)
   - Validar consistência (somas batem com totais oficiais)
   - Detectar outliers (preço 10x acima da média)

### Páginas/Rotas Necessárias:
- `/admin/data-pipelines` (status de jobs)
- `/admin/data-quality` (relatórios de qualidade)

### Componentes/Features:
- Job scheduler (Kubernetes CronJobs)
- Job monitoring dashboard (Airflow-like UI)
- Manual trigger (run job on-demand)
- Alerting (Slack, email, PagerDuty)
- Data quality framework (Great Expectations ou custom)
- Automated backups (PostgreSQL pg_dump daily)
- Disaster recovery plan (restore from backup)

### Success Metrics:
- Job success rate: >99%
- Data freshness: <24h lag
- Data quality issues detected: >90% (before impacting users)
- Backup restore time: <2h (RTO)

### Effort Estimate:
10 dias de desenvolvimento (job monitoring + data quality + alerting + backups)

---

# ROADMAP VISUAL - Sequência de Implementação

## Fases de Desenvolvimento

### FASE 0: MVP Beta (Atual → +4 semanas)
**Objetivo:** Validar core value prop com 50 early adopters

**Jornadas Críticas (P0 Blockers):**
1. ✅ J-AC03 (Parcial): Simulador Anonymous (DONE)
2. ❌ J-AC02: Sistema de Autenticação
3. ❌ J-AC01: Onboarding First-Time User
4. ❌ J-AC03 (Completo): Simulador Logged-In + Save
5. ❌ J-REV01: Descoberta Premium (Paywall)
6. ❌ J-REV02: Checkout e Pagamento

**Deliverables:**
- Usuários podem criar conta, simular, salvar resultados
- Paywall soft implementado (features premium bloqueadas)
- Checkout funcional (PIX + Cartão)
- 50 beta users onboarded

**Timeline:** 4 semanas (Semanas 2-5 de 2026)

**Success Criteria:**
- 50 signups
- 30% activation (fazem primeira simulação)
- 3 conversões para premium
- NPS beta: >40

---

### FASE 1: V1 Public Launch (Semanas 6-12 de 2026)
**Objetivo:** Lançamento público, acquisition channels, retention básica

**Jornadas Prioritárias (P1):**
1. J-A01: SEO/Blog (acquisition orgânica)
2. J-A03: Programa de Referral (viral growth)
3. J-AC04: Dashboard de Market Intelligence
4. J-AC05: Alertas e Notificações
5. J-R01: Histórico e Comparação de Simulações
6. J-R02: Watchlist de Destinos
7. J-R06: Email Marketing (re-engagement)
8. J-REV03: Billing Dashboard
9. J-ADM03: Analytics BI
10. J-ADM05: Suporte ao Cliente

**Deliverables:**
- Blog com 20 artigos SEO-optimized
- Referral program ativo (K-factor >0.5)
- Dashboard de insights lançado
- Email drip campaigns configuradas
- Analytics dashboards operacionais
- Ticketing system para suporte

**Timeline:** 6 semanas

**Success Criteria:**
- 500 signups totais
- 5k organic visits/month
- Referral: 30% dos signups
- MRR: R$10k
- Churn: <8%

---

### FASE 2: V2 Platform Maturity (Meses 4-6 de 2026)
**Objetivo:** Marketplace, matchmaking, integrações, enterprise

**Jornadas Prioritárias (P2):**
1. J-A02: Campanhas Pagas (Google Ads)
2. J-A04: Parcerias com Associações
3. J-R03: Exploração de NCMs Relacionados
4. J-R04: Análise de Concorrentes
5. J-REV04: Plano Enterprise
6. J-REV05: Upsell de Add-Ons
7. J-REF01: Programa de Indicação (enhanciar)
8. J-REF03: Co-Marketing com Parceiros
9. J-ADM01: Onboarding de Parceiros
10. J-ADM02: Moderação de Conteúdo

**Deliverables:**
- Google Ads + LinkedIn Ads rodando (CAC <R$150)
- 5 parcerias ativas (APEX, ABIEC, etc.)
- Enterprise tier com 2 clientes
- Add-ons marketplace lançado
- Partner portal operacional

**Timeline:** 3 meses

**Success Criteria:**
- 2k signups totais
- MRR: R$50k
- Enterprise deals: 2 (avg R$10k/mês cada)
- Partners onboarded: 20

---

### FASE 3: Scale & Optimization (Meses 7-12 de 2026)
**Objetivo:** Escalar, otimizar conversão, network effects

**Jornadas Prioritárias (P3 + Advanced):**
1. J-A05: Webinars e Eventos
2. J-R05: Comunidade/Fórum
3. J-REF02: Integrações CRM/ERP
4. J-REF04: User-Generated Content
5. J-REF05: Programa de Afiliados
6. J-ADM04: Feature Flags e Configs
7. J-ADM06: Data Quality e Pipelines

**Deliverables:**
- Fórum ativo (50+ perguntas/mês)
- Salesforce/HubSpot integrations
- 50 afiliados ativos
- 12 case studies publicados
- Data quality framework completo

**Timeline:** 6 meses

**Success Criteria:**
- 5k signups totais
- MRR: R$100k
- Community: 500 membros ativos
- Affiliate revenue: 20% de MRR
- NPS: >60

---

## Priorização RICE - Top 15 Jornadas

| ID | Jornada | Reach | Impact | Confidence | Effort (days) | RICE Score | Fase |
|----|---------|-------|--------|------------|---------------|------------|------|
| J-AC02 | Autenticação e Perfil | 1000 | 3 | 0.9 | 12 | 225 | 0 (MVP) |
| J-AC01 | Onboarding First-Time | 1000 | 3 | 0.8 | 8 | 300 | 0 (MVP) |
| J-REV01 | Descoberta Premium | 800 | 3 | 0.8 | 6 | 320 | 0 (MVP) |
| J-REV02 | Checkout Pagamento | 500 | 3 | 0.7 | 12 | 88 | 0 (MVP) |
| J-A01 | SEO/Blog | 1200 | 2 | 0.7 | 12 | 140 | 1 (V1) |
| J-A03 | Referral Program | 800 | 3 | 0.6 | 8 | 180 | 1 (V1) |
| J-AC04 | Dashboard Intelligence | 700 | 2 | 0.7 | 10 | 98 | 1 (V1) |
| J-R01 | Histórico Comparação | 600 | 2 | 0.8 | 7 | 137 | 1 (V1) |
| J-R06 | Email Re-Engagement | 900 | 2 | 0.7 | 7 | 180 | 1 (V1) |
| J-ADM03 | Analytics BI | 50 | 3 | 0.8 | 12 | 10 | 1 (V1) |
| J-A04 | Parcerias Associações | 500 | 3 | 0.5 | 10 | 75 | 2 (V2) |
| J-REV04 | Enterprise Tier | 20 | 3 | 0.6 | 10 | 3.6 | 2 (V2) |
| J-ADM01 | Onboard Parceiros | 200 | 2 | 0.7 | 10 | 28 | 2 (V2) |
| J-R05 | Fórum Comunidade | 400 | 2 | 0.4 | 15 | 21 | 3 (Scale) |
| J-REF05 | Programa Afiliados | 600 | 2 | 0.5 | 12 | 50 | 3 (Scale) |

---

## Dependencies Graph (Críticas)

```
J-AC02 (Auth)
  ├── J-AC01 (Onboarding) [BLOCKS TUDO]
  ├── J-AC03 (Simulator Logged)
  ├── J-AC04 (Dashboard)
  ├── J-AC05 (Alertas)
  ├── J-R01 (Histórico)
  ├── J-R02 (Watchlist)
  ├── J-REV01 (Premium Discovery)
  └── J-ADM05 (Support Tickets)

J-REV01 (Premium Discovery)
  └── J-REV02 (Checkout) [MONETIZATION BLOCKER]

J-REV02 (Checkout)
  ├── J-REV03 (Billing Management)
  ├── J-REV04 (Enterprise)
  └── J-REV05 (Add-Ons)

J-AC03 (Simulator)
  ├── J-R01 (Histórico)
  ├── J-R02 (Watchlist)
  └── J-R03 (NCMs Relacionados)

J-A04 (Parcerias)
  ├── J-A05 (Eventos)
  └── J-REF03 (Co-Marketing)

J-REF01 (Referral Basic)
  └── J-REF05 (Affiliate Program - enhanced)

J-R05 (Fórum)
  └── J-ADM02 (Moderação)
```

---

## Quick Wins vs. Long-Term Bets

### Quick Wins (High Impact, Low Effort)
1. **J-AC05: Alertas Email** (7 dias, RICE 180) - retention imediata
2. **J-R02: Watchlist** (6 dias, retention) - feature simples, alto valor
3. **J-REV01: Paywall Soft** (6 dias, RICE 320) - monetization enabler

### Long-Term Bets (Strategic, High Effort)
1. **J-R05: Fórum Comunidade** (15 dias) - network effect, mas lento
2. **J-REF02: Integrações CRM** (15 dias) - lock-in enterprise, mas nicho
3. **J-REF05: Afiliados** (12 dias) - scalable acquisition, mas setup complexo

---

## North Star Metric Alignment

**North Star:** Time-to-First-Export-Match

**Jornadas que mais impactam NSM:**
1. J-AC01: Onboarding (reduce time-to-first-value)
2. J-AC03: Simulador (core value delivery)
3. J-AC04: Dashboard (insights → decision → action)
4. J-R01: Histórico (comparação → decisão mais rápida)
5. Futuro: J-MATCH01 (Buyer Matching - não mapeado aqui, será Epic 8)

---

## Gaps Identificados (Jornadas Futuras)

**Marketplace & Matchmaking (Q2 2026):**
- J-MATCH01: Buyer Discovery (exportador encontra compradores)
- J-MATCH02: RFQ Flow (Request for Quote)
- J-MATCH03: Negociação In-Platform
- J-MATCH04: Deal Closure Tracking

**Operational Enablement (Q3 2026):**
- J-OPS01: Document Automation (Invoice, Packing List)
- J-OPS02: Siscomex Integration
- J-OPS03: Logistics Quote Comparison
- J-OPS04: Shipment Tracking

**Financial Services (Q3-Q4 2026):**
- J-FIN01: FX Simulation
- J-FIN02: Trade Finance Application
- J-FIN03: Credit Scoring
- J-FIN04: Insurance Quotes

---

## Conclusão

**Status Atual:**
- ✅ 1 jornada completa (Simulator Anonymous)
- ❌ 47 jornadas mapeadas e pendentes
- 🎯 6 jornadas críticas (P0) para MVP Beta
- 📅 Roadmap de 12 meses para produto maduro

**Próximos Passos Imediatos (Semana 2-5/2026):**
1. Implementar J-AC02 (Auth) - 12 dias
2. Implementar J-AC01 (Onboarding) - 8 dias
3. Integrar J-AC03 (Simulator + Auth) - 6 dias
4. Implementar J-REV01 + J-REV02 (Monetization) - 18 dias

**Total Effort Estimado:** 44 dias = ~2 meses com 1 dev full-time (ou 1 mês com 2 devs)

---

**Versão:** 1.0.0
**Última Atualização:** 2026-01-20
**Responsável:** BGC Product Management Team
**Próxima Revisão:** 2026-02-01 (Sprint Planning Fase 0)
