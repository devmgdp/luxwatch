package main

import (
	"luxwatch/internal/cron"
	"luxwatch/internal/database"
	"luxwatch/migrations/scraper"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// --- CARREGAMENTO DE .ENV ---
	godotenv.Load()

	database.InitDB()
	cron.StartPriceMonitor()

	go func() {
		scraper.RunDealFinder()
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			scraper.RunDealFinder()
		}
	}()

	r := gin.Default()

	// Configuração de Sessão (Obrigatório para o login funcionar)
	// .env -> Se não existir, usa uma password padrão por segurança.
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "secret_key_luxwatch_2026"
	}
	store := cookie.NewStore([]byte(secret))
	r.Use(sessions.Sessions("luxsession", store))

	r.LoadHTMLGlob("suggestions/templates/*.html")
	r.Static("/static", "./static")

	// --- ROTAS DE AUTENTICAÇÃO ---

	r.GET("/auth", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth.html", nil)
	})

	// Mostra a página de Cadastro (ESSA É A QUE ESTAVA FALTANDO!)
	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})

	// --- ROTAS DE AÇÃO (POST - Processar Dados) ---

	r.POST("/signup", func(c *gin.Context) {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		// Validação de segurança
		if len(password) < 6 {
			c.String(200, `<div style="color: #ff9f43; border: 1px solid #ff9f43; padding: 10px; border-radius: 8px; margin-top: 10px; font-family: monospace; font-size: 0.8rem; background: rgba(255, 159, 67, 0.1);">⚠️ ERRO: A SENHA DEVE TER NO MÍNIMO 6 CARACTERES</div>`)
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		_, err := database.DB.Exec("INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)",
			username, email, string(hash))

		if err != nil {
			c.String(200, `<div style="color: #ff4d4d; border: 1px solid #ff4d4d; padding: 10px; border-radius: 8px; margin-top: 10px; font-family: monospace; font-size: 0.8rem; background: rgba(255, 77, 77, 0.1);">⚠️ ERRO: USUÁRIO OU EMAIL JÁ EXISTENTE</div>`)
			return
		}

		// Se o cadastro for sucesso, o HTMX redireciona para o login
		c.Writer.Header().Set("HX-Redirect", "/auth")
		c.Status(http.StatusOK)
	})

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		var id int
		var hash string
		err := database.DB.QueryRow("SELECT id, password_hash FROM users WHERE username=$1", username).Scan(&id, &hash)

		if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
			c.String(200, `<div style="color: #ff4d4d; border: 1px solid #ff4d4d; padding: 10px; border-radius: 8px; margin-top: 10px; font-family: monospace; font-size: 0.8rem; background: rgba(255, 77, 77, 0.1); animation: glitch-main 0.3s infinite;">⚠️ ACESSO NEGADO: CREDENCIAIS INVÁLIDAS</div>`)
			return
		}

		session := sessions.Default(c)
		session.Set("user_id", id)
		session.Save()

		c.Writer.Header().Set("HX-Redirect", "/")
		c.Status(http.StatusOK)
	})

	// --- ROTAS DO PROJETO ---

	// --- ROTA HOME COM USERNAME ---
	r.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		var username string
		if userID != nil {
			// Busca o nome do usuário logado
			database.DB.QueryRow("SELECT username FROM users WHERE id=$1", userID).Scan(&username)
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Season":      "Promoções Ativas",
			"Suggestions": []string{"Eletrônicos", "Casa", "Saúde"},
			"IsLogged":    userID != nil,
			"Username":    username,
		})
	})

	r.GET("/category/:name", func(c *gin.Context) {
		cat := c.Param("name")

		// Usamos ILIKE para ignorar maiúsculas/minúsculas no PostgreSQL
		// Ou usamos a função LOWER() para garantir a compatibilidade
		rows, err := database.DB.Query("SELECT name, price, image_url, product_url FROM daily_deals WHERE LOWER(category)=LOWER($1)", cat)

		if err != nil {
			c.String(500, "Erro ao buscar categorias")
			return
		}
		defer rows.Close()

		var products []map[string]interface{}
		for rows.Next() {
			var n, img, url string
			var p float64
			rows.Scan(&n, &p, &img, &url)
			products = append(products, map[string]interface{}{
				"Name":     n,
				"Price":    p,
				"ImageURL": img,
				"URL":      url,
			})
		}

		c.HTML(http.StatusOK, "category.html", gin.H{
			"Category": cat,
			"Products": products,
		})
	})

	r.POST("/save-product", func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.String(401, "❌ Acesse sua conta primeiro!")
			return
		}

		url := c.PostForm("url")
		priceStr := c.PostForm("price")
		price, _ := strconv.ParseFloat(priceStr, 64)

		_, err := database.DB.Exec(`
      INSERT INTO watchlist (user_id, product_url, target_price) 
      VALUES ($1, $2, $3) 
      ON CONFLICT (user_id, product_url) DO UPDATE SET target_price = $3`,
			userID, url, price)

		if err != nil {
			c.String(500, "Erro ao salvar")
			return
		}
		c.String(200, "✅ Monitorando!")
	})

	r.GET("/saves", func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.Redirect(http.StatusSeeOther, "/auth")
			return
		}

		// SQL PODEROSO: Busca os dados do produto E o nome do dono (username)
		query := `
        SELECT w.product_url, w.target_price, u.username 
        FROM watchlist w
        JOIN users u ON w.user_id = u.id
        WHERE w.user_id = $1`

		rows, err := database.DB.Query(query, userID)
		if err != nil {
			c.String(500, "Erro ao carregar lista")
			return
		}
		defer rows.Close()

		var savedItems []map[string]interface{}
		for rows.Next() {
			var url, username string
			var price float64
			rows.Scan(&url, &price, &username)

			savedItems = append(savedItems, map[string]interface{}{
				"URL":      url,
				"Price":    price,
				"Username": username, // Agora passamos o nome para o HTML
			})
		}
		c.HTML(http.StatusOK, "saves.html", gin.H{"Items": savedItems})
	})

	r.GET("/api/history/*name", func(c *gin.Context) {
		name := c.Param("name")
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}
		rows, err := database.DB.Query("SELECT price, checked_at FROM price_history WHERE product_name=$1 ORDER BY checked_at ASC", name)
		if err != nil {
			c.JSON(500, gin.H{"error": "Erro no histórico"})
			return
		}
		defer rows.Close()

		var history []map[string]interface{}
		for rows.Next() {
			var p float64
			var d time.Time
			rows.Scan(&p, &d)
			history = append(history, map[string]interface{}{"price": p, "date": d.Format(time.RFC3339)})
		}
		c.JSON(200, history)
	})

	// --- BLOQUEIO NO SCRAPE ---
	r.POST("/scrape", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user_id") == nil {
			// Mensagem pedindo para logar, estilizada para a Index
			c.String(200, `
            <div class="glass-panel" style="border-color: var(--accent); margin-top: 20px; text-align: center;">
                <p style="color: var(--text-main);">🔒 PROTOCOLO BLOQUEADO</p>
                <p style="color: var(--text-muted); font-size: 0.9rem;">Logue ou crie uma conta para pesquisar um produto.</p>
                <a href="/auth" class="btn-buy" style="display: inline-block; margin-top: 10px; text-decoration: none;">ACESSAR TERMINAL</a>
            </div>
        `)
			return
		}

		url := c.PostForm("url")
		product, err := scraper.ExtractData(url)
		if err != nil {
			c.String(500, "Erro ao analisar link")
			return
		}
		c.HTML(http.StatusOK, "product_card.html", product)
	})

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear() // Limpa o ID do usuário da sessão
		session.Save()
		c.Redirect(http.StatusSeeOther, "/") // Manda de volta pra Home
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
