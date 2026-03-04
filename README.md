# LuxWatch - Monitor de Ofertas Inteligente

O **LuxWatch** é uma plataforma robusta de monitoramento de preços que automatiza a busca por ofertas em grandes marketplaces. O sistema combina a performance do **Go** com a flexibilidade do **Playwright** para oferecer alertas em tempo real diretamente no e-mail do usuário.

🔗 **Acesse o projeto ao vivo:** [https://luxwatch-8ysj.onrender.com/](https://luxwatch-8ysj.onrender.com/)

---

## 🎯 Funcionalidades

*   **Web Scraping Resiliente:** Coleta automática de dados utilizando técnicas de evasão (Stealth) para contornar proteções anti-bot.
*   **Watchlist Dinâmica:** Usuários podem cadastrar produtos e definir um "Preço Alvo".
*   **Notificações Automáticas:** O backend monitora os preços 24/7 e dispara e-mails via SMTP quando o valor desejado é atingido.
*   **Autenticação Segura:** Sistema de login com criptografia de senhas via **Bcrypt**.
*   **Proteção & Segurança:** Implementação de proteção contra ataques básicos e gerenciamento de sessões criptografadas.

---

## 🏗️ Arquitetura Técnica

O projeto utiliza uma arquitetura contida em ambiente **Docker**, dividida em:

1.  **Backend (Golang):** API REST utilizando o framework **Gin**, responsável pela lógica de negócios, autenticação e gerenciamento do banco de dados.
2.  **Scraper Engine (Node.js + Playwright):** Um módulo especializado em navegação headless que supera desafios de JavaScript dinâmico.
3.  **Database (PostgreSQL):** Armazenamento persistente de usuários, produtos e histórico de preços.

---

## 🛠️ Stack Tecnológica

*   **Linguagem Principal:** Go (v1.24)
*   **Web Framework:** Gin Gonic
*   **Automação de Browser:** Playwright Stealth (Node.js)
*   **Banco de Dados:** PostgreSQL
*   **Containerização:** Docker & Docker Compose
*   **Deploy:** Render Cloud

---

## 📸 Demonstração do projeto

<div align="center">
  <h3>Página Principal & Busca</h3>
  <img src="./screenshots/home.png" width="800px" alt="Home LuxWatch">
  <p><i>Interface limpa com listagem de ofertas capturadas pelo scraper em tempo real.</i></p>
  
  <br>

  <table align="center">
    <tr>
      <td align="center">
        <b>Produtos</b><br>
        <img src="./screenshots/category.png" width="380px" alt="Dashboard">
      </td>
      <td align="center">
        <b>Alerta de Preço (E-mail)</b><br>
        <img src="./screenshots/email.png" width="380px" alt="Alerta E-mail">
      </td>
    </tr>
  </table>
  <p><i>Acompanhamento de preços alvos e sistema de notificação automática via SMTP.</i></p>
</div>

---

## 🚀 Como Executar o Projeto

### Pré-requisitos
*   Docker e Docker Compose instalados.

### Passo a Passo

1.  **Clone o repositório:**
    ```bash
    git clone [https://github.com/devmgdp/luxwatch.git](https://github.com/devmgdp/luxwatch.git)
    cd luxwatch
    ```

2.  **Configure as Variáveis de Ambiente:**
    Crie um arquivo `.env` na raiz do projeto conforme o exemplo:
    ```env
    DB_URL=postgres://user:pass@localhost:5432/luxwatch?sslmode=disable
    SMTP_EMAIL=seu-email@gmail.com
    SMTP_PASS=sua-senha-de-app
    SESSION_SECRET=sua-chave-secreta
    ```

3.  **Suba os Containers:**
    ```bash
    docker-compose up --build
    ```
    A aplicação estará disponível em `http://localhost:8080`.

---

## 🧠 Desafios ao fazer o Projeto

*   **Evasão Anti-Bot:** Implementação de técnicas de *stealth browsing* para evitar desafios de verificação de JavaScript (Mercado Livre).
*   **Concorrência em Go:** Uso de Goroutines para rodar o scraper em background sem travar a API.
*   **Otimização de Docker:** Configuração de ambiente Linux headless para suportar binários do Chromium no Render.

---

## 🤝 Contato

Desenvolvido por **Miguel Peixoto** - [Seu LinkedIn](www.linkedin.com/in/devmgdp)