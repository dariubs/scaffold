package index

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/dariubs/scaffold/app/config"
	"github.com/dariubs/scaffold/app/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"gorm.io/gorm"
)

func getGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.C.GoogleOAuth.ClientID,
		ClientSecret: config.C.GoogleOAuth.ClientSecret,
		RedirectURL:  config.C.GoogleOAuth.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.C.OAuthGoogleEnabled() {
			c.Redirect(http.StatusFound, "/login?error=google_disabled")
			return
		}
		googleOauthConfig := getGoogleOAuthConfig()

		// Generate state token for CSRF protection
		state := uuid.New().String()
		session := sessions.Default(c)
		session.Set("oauth_state", state)
		if err := session.Save(); err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{
				"Title": "Login",
				"Error": "Failed to save session",
			})
			return
		}

		url := googleOauthConfig.AuthCodeURL(state)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func GoogleCallback(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.C.OAuthGoogleEnabled() {
			c.Redirect(http.StatusFound, "/login?error=google_disabled")
			return
		}
		googleOauthConfig := getGoogleOAuthConfig()

		code := c.Query("code")
		state := c.Query("state")

		if code == "" {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Title": "Login",
				"Error": "Authorization code not received",
			})
			return
		}

		// Validate state token (CSRF protection)
		session := sessions.Default(c)
		savedState := session.Get("oauth_state")
		if savedState == nil || savedState.(string) != state {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Title": "Login",
				"Error": "Invalid state token",
			})
			return
		}

		// Clear state from session
		session.Delete("oauth_state")
		session.Save()

		// Exchange code for token
		token, err := googleOauthConfig.Exchange(context.Background(), code)
		if err != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Title": "Login",
				"Error": "Failed to exchange token",
			})
			return
		}

		// Get user info from Google
		client := googleOauthConfig.Client(context.Background(), token)
		oauth2Service, err := googleoauth2.New(client)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{
				"Title": "Login",
				"Error": "Failed to create OAuth service",
			})
			return
		}

		userInfo, err := oauth2Service.Userinfo.Get().Do()
		if err != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Title": "Login",
				"Error": "Failed to get user info",
			})
			return
		}

		// Check if user exists
		var user model.User
		err = db.Where("google_id = ? OR email = ?", userInfo.Id, userInfo.Email).First(&user).Error

		if err != nil {
			// User doesn't exist, create new one
			user = model.User{
				Username:    userInfo.Email, // Use email as username for OAuth users
				Email:       userInfo.Email,
				Password:    "", // No password for OAuth users
				Name:        userInfo.Name,
				AvatarURL:   userInfo.Picture,
				GoogleID:    userInfo.Id,
				LoginMethod: "google",
			}

			if err := db.Create(&user).Error; err != nil {
				c.HTML(http.StatusInternalServerError, "login.html", gin.H{
					"Title": "Login",
					"Error": "Failed to create user account",
				})
				return
			}
		} else {
			// User exists, update Google ID if not set
			if user.GoogleID == "" {
				user.GoogleID = userInfo.Id
				user.LoginMethod = "google"
				user.AvatarURL = userInfo.Picture
				user.Name = userInfo.Name
				db.Save(&user)
			}
		}

		// Set session
		session = sessions.Default(c)
		session.Set("user_id", user.Model.ID)
		session.Save()

		c.Redirect(http.StatusFound, "/")
	}
}

// GitHub OAuth

func getGitHubOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.C.GitHubOAuth.ClientID,
		ClientSecret: config.C.GitHubOAuth.ClientSecret,
		RedirectURL:  config.C.GitHubOAuth.RedirectURL,
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}
}

func GitHubLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.C.OAuthGitHubEnabled() {
			c.Redirect(http.StatusFound, "/login?error=github_disabled")
			return
		}
		cfg := getGitHubOAuthConfig()
		state := uuid.New().String()
		session := sessions.Default(c)
		session.Set("oauth_state", state)
		if err := session.Save(); err != nil {
			c.Redirect(http.StatusFound, "/login?error=session")
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, cfg.AuthCodeURL(state))
	}
}

type githubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func GitHubCallback(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.C.OAuthGitHubEnabled() {
			c.Redirect(http.StatusFound, "/login?error=github_disabled")
			return
		}
		code := c.Query("code")
		state := c.Query("state")
		if code == "" {
			c.Redirect(http.StatusFound, "/login?error=code")
			return
		}
		session := sessions.Default(c)
		savedState, _ := session.Get("oauth_state").(string)
		if savedState == "" || savedState != state {
			c.Redirect(http.StatusFound, "/login?error=state")
			return
		}
		session.Delete("oauth_state")
		session.Save()

		cfg := getGitHubOAuthConfig()
		token, err := cfg.Exchange(context.Background(), code)
		if err != nil {
			c.Redirect(http.StatusFound, "/login?error=exchange")
			return
		}
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "https://api.github.com/user", nil)
		token.SetAuthHeader(req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.Redirect(http.StatusFound, "/login?error=userinfo")
			return
		}
		defer resp.Body.Close()
		var gu githubUser
		if json.NewDecoder(resp.Body).Decode(&gu) != nil {
			c.Redirect(http.StatusFound, "/login?error=userinfo")
			return
		}
		email := gu.Email
		if email == "" {
			// try /user/emails
			req2, _ := http.NewRequestWithContext(context.Background(), "GET", "https://api.github.com/user/emails", nil)
			token.SetAuthHeader(req2)
			resp2, err := http.DefaultClient.Do(req2)
			if err == nil && resp2.StatusCode == http.StatusOK {
				var emails []struct {
					Email   string `json:"email"`
					Primary bool   `json:"primary"`
				}
				if json.NewDecoder(resp2.Body).Decode(&emails) == nil {
					for _, e := range emails {
						if e.Primary {
							email = e.Email
							break
						}
					}
					if email == "" && len(emails) > 0 {
						email = emails[0].Email
					}
				}
				resp2.Body.Close()
			}
		}
		if email == "" {
			email = gu.Login + "@github.user"
		}

		var user model.User
		githubIDStr := fmt.Sprintf("%d", gu.ID)
		err = db.Where("github_id = ? OR email = ?", githubIDStr, email).First(&user).Error
		if err != nil {
			user = model.User{
				Username:    gu.Login,
				Email:       email,
				Password:    "",
				Name:        gu.Name,
				AvatarURL:   gu.AvatarURL,
				GitHubID:    githubIDStr,
				LoginMethod: "github",
			}
			if db.Create(&user).Error != nil {
				c.Redirect(http.StatusFound, "/login?error=create")
				return
			}
		} else {
			if user.GitHubID == "" {
				user.GitHubID = githubIDStr
				user.LoginMethod = "github"
				user.AvatarURL = gu.AvatarURL
				user.Name = gu.Name
				db.Save(&user)
			}
		}
		session = sessions.Default(c)
		session.Set("user_id", user.Model.ID)
		session.Save()
		c.Redirect(http.StatusFound, "/")
	}
}

// LinkedIn OAuth

var linkedInEndpoint = oauth2.Endpoint{
	AuthURL:  "https://www.linkedin.com/oauth/v2/authorization",
	TokenURL: "https://www.linkedin.com/oauth/v2/accessToken",
}

func getLinkedInOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.C.LinkedInOAuth.ClientID,
		ClientSecret: config.C.LinkedInOAuth.ClientSecret,
		RedirectURL:  config.C.LinkedInOAuth.RedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     linkedInEndpoint,
	}
}

func LinkedInLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.C.OAuthLinkedInEnabled() {
			c.Redirect(http.StatusFound, "/login?error=linkedin_disabled")
			return
		}
		cfg := getLinkedInOAuthConfig()
		state := uuid.New().String()
		session := sessions.Default(c)
		session.Set("oauth_state", state)
		if err := session.Save(); err != nil {
			c.Redirect(http.StatusFound, "/login?error=session")
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, cfg.AuthCodeURL(state))
	}
}

type linkedInUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func LinkedInCallback(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.C.OAuthLinkedInEnabled() {
			c.Redirect(http.StatusFound, "/login?error=linkedin_disabled")
			return
		}
		code := c.Query("code")
		state := c.Query("state")
		if code == "" {
			c.Redirect(http.StatusFound, "/login?error=code")
			return
		}
		session := sessions.Default(c)
		savedState, _ := session.Get("oauth_state").(string)
		if savedState == "" || savedState != state {
			c.Redirect(http.StatusFound, "/login?error=state")
			return
		}
		session.Delete("oauth_state")
		session.Save()

		cfg := getLinkedInOAuthConfig()
		token, err := cfg.Exchange(context.Background(), code)
		if err != nil {
			c.Redirect(http.StatusFound, "/login?error=exchange")
			return
		}
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "https://api.linkedin.com/v2/userinfo", nil)
		token.SetAuthHeader(req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.Redirect(http.StatusFound, "/login?error=userinfo")
			return
		}
		defer resp.Body.Close()
		var lu linkedInUser
		if json.NewDecoder(resp.Body).Decode(&lu) != nil {
			c.Redirect(http.StatusFound, "/login?error=userinfo")
			return
		}
		email := lu.Email
		if email == "" {
			email = lu.Sub + "@linkedin.user"
		}

		var user model.User
		err = db.Where("linkedin_id = ? OR email = ?", lu.Sub, email).First(&user).Error
		if err != nil {
			user = model.User{
				Username:    lu.Email,
				Email:       email,
				Password:    "",
				Name:        lu.Name,
				AvatarURL:   lu.Picture,
				LinkedInID:  lu.Sub,
				LoginMethod: "linkedin",
			}
			if user.Username == "" {
				user.Username = "linkedin_" + lu.Sub
			}
			if db.Create(&user).Error != nil {
				c.Redirect(http.StatusFound, "/login?error=create")
				return
			}
		} else {
			if user.LinkedInID == "" {
				user.LinkedInID = lu.Sub
				user.LoginMethod = "linkedin"
				user.AvatarURL = lu.Picture
				user.Name = lu.Name
				db.Save(&user)
			}
		}
		session = sessions.Default(c)
		session.Set("user_id", user.Model.ID)
		session.Save()
		c.Redirect(http.StatusFound, "/")
	}
}

// X (Twitter) OAuth 2.0 with PKCE

func pkceCodeVerifierAndChallenge() (verifier, challenge string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	verifier = base64.RawURLEncoding.EncodeToString(b)
	if len(verifier) < 43 {
		verifier = verifier + strings.Repeat("a", 43-len(verifier))
	} else if len(verifier) > 128 {
		verifier = verifier[:128]
	}
	h := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(h[:])
	return verifier, challenge, nil
}

var xEndpoint = oauth2.Endpoint{
	AuthURL:  "https://twitter.com/i/oauth2/authorize",
	TokenURL: "https://api.twitter.com/2/oauth2/token",
}

func getXOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.C.XOAuth.ClientID,
		ClientSecret: config.C.XOAuth.ClientSecret,
		RedirectURL:  config.C.XOAuth.RedirectURL,
		Scopes:       []string{"users.read", "tweet.read"},
		Endpoint:     xEndpoint,
	}
}

func XLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.C.OAuthXEnabled() {
			c.Redirect(http.StatusFound, "/login?error=x_disabled")
			return
		}
		cfg := getXOAuthConfig()
		verifier, challenge, err := pkceCodeVerifierAndChallenge()
		if err != nil {
			c.Redirect(http.StatusFound, "/login?error=session")
			return
		}
		state := uuid.New().String()
		session := sessions.Default(c)
		session.Set("oauth_state", state)
		session.Set("oauth_code_verifier", verifier)
		if err := session.Save(); err != nil {
			c.Redirect(http.StatusFound, "/login?error=session")
			return
		}
		authURL := cfg.AuthCodeURL(state, oauth2.SetAuthURLParam("code_challenge", challenge), oauth2.SetAuthURLParam("code_challenge_method", "S256"))
		c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}

type xUserData struct {
	Data struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"data"`
}

func xExchangeCodeForToken(code, codeVerifier string) (accessToken string, err error) {
	cfg := getXOAuthConfig()
	data := url.Values{}
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", cfg.RedirectURL)
	data.Set("code_verifier", codeVerifier)
	req, err := http.NewRequestWithContext(context.Background(), "POST", cfg.Endpoint.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(cfg.ClientID, cfg.ClientSecret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token exchange: %s", string(body))
	}
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	if json.Unmarshal(body, &tok) != nil {
		return "", fmt.Errorf("invalid token response")
	}
	return tok.AccessToken, nil
}

func XCallback(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.C.OAuthXEnabled() {
			c.Redirect(http.StatusFound, "/login?error=x_disabled")
			return
		}
		code := c.Query("code")
		state := c.Query("state")
		if code == "" {
			c.Redirect(http.StatusFound, "/login?error=code")
			return
		}
		session := sessions.Default(c)
		savedState, _ := session.Get("oauth_state").(string)
		verifier, _ := session.Get("oauth_code_verifier").(string)
		if savedState == "" || savedState != state || verifier == "" {
			session.Delete("oauth_state")
			session.Delete("oauth_code_verifier")
			session.Save()
			c.Redirect(http.StatusFound, "/login?error=state")
			return
		}
		session.Delete("oauth_state")
		session.Delete("oauth_code_verifier")
		session.Save()

		accessToken, err := xExchangeCodeForToken(code, verifier)
		if err != nil {
			c.Redirect(http.StatusFound, "/login?error=exchange")
			return
		}
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "https://api.twitter.com/2/users/me?user.fields=profile_image_url", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.Redirect(http.StatusFound, "/login?error=userinfo")
			return
		}
		defer resp.Body.Close()
		var xu xUserData
		if json.NewDecoder(resp.Body).Decode(&xu) != nil || xu.Data.ID == "" {
			c.Redirect(http.StatusFound, "/login?error=userinfo")
			return
		}
		email := xu.Data.Username + "@x.user"

		var user model.User
		err = db.Where("x_id = ? OR username = ?", xu.Data.ID, xu.Data.Username).First(&user).Error
		if err != nil {
			user = model.User{
				Username:    xu.Data.Username,
				Email:       email,
				Password:    "",
				Name:        xu.Data.Name,
				XID:         xu.Data.ID,
				LoginMethod: "x",
			}
			if db.Create(&user).Error != nil {
				c.Redirect(http.StatusFound, "/login?error=create")
				return
			}
		} else {
			if user.XID == "" {
				user.XID = xu.Data.ID
				user.LoginMethod = "x"
				user.Name = xu.Data.Name
				db.Save(&user)
			}
		}
		session = sessions.Default(c)
		session.Set("user_id", user.Model.ID)
		session.Save()
		c.Redirect(http.StatusFound, "/")
	}
}
