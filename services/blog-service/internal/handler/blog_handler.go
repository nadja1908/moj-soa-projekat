package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"blog-service/internal/model"
	"blog-service/internal/store"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	store                  *store.Store
	stakeholdersServiceURL string
	followerServiceURL     string
	httpClient             *http.Client
}

func NewBlogHandler(store *store.Store, stakeholdersServiceURL, followerServiceURL string) *BlogHandler {
	return &BlogHandler{
		store:                  store,
		stakeholdersServiceURL: stakeholdersServiceURL,
		followerServiceURL:     followerServiceURL,
		httpClient:             &http.Client{Timeout: 5 * time.Second},
	}
}

// CreateBlogPost kreira novi blog post
// CreateBlogPost kreira novi blog post
func (h *BlogHandler) CreateBlogPost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req model.CreateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &model.BlogPost{
		UserID:      userID.(int64),
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
	}

	// 1️⃣ prvo kreiramo sam blog post
	if err := h.store.CreateBlogPost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog post"})
		return
	}

	// 2️⃣ ako postoje slike — čuvamo ih
	if len(req.Images) > 0 {
		if err := h.store.CreateBlogImages(post.ID, req.Images); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save images"})
			return
		}
	}

	// 3️⃣ vraćamo korisniku potvrdu
	c.JSON(http.StatusCreated, gin.H{
		"message": "Blog post created successfully",
		"post":    post,
	})
}

// fetchAuthor dohvata informacije o autoru iz stakeholders servisa
func (h *BlogHandler) fetchAuthor(userID int64) *model.Author {
	url := fmt.Sprintf("%s/internal/users/%d", h.stakeholdersServiceURL, userID)

	resp, err := h.httpClient.Get(url)
	if err != nil {
		fmt.Printf("Error fetching author: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Non-OK status from stakeholders service: %d\n", resp.StatusCode)
		return nil
	}

	var response struct {
		User struct {
			ID       int64  `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("Error decoding user response: %v\n", err)
		return nil
	}

	return &model.Author{
		ID:       response.User.ID,
		Username: response.User.Username,
		Email:    response.User.Email,
	}
}

// enrichPostsWithAuthors dodaje author informacije u blog postove
func (h *BlogHandler) enrichPostsWithAuthors(posts []model.BlogPost) {
	for i := range posts {
		if author := h.fetchAuthor(posts[i].UserID); author != nil {
			posts[i].Author = author
		}
	}
}

// GetAllBlogPosts vraća sve blog postove
func (h *BlogHandler) GetAllBlogPosts(c *gin.Context) {
	posts, err := h.store.GetAllBlogPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get blog posts"})
		return
	}

	// Dodaj author informacije
	h.enrichPostsWithAuthors(posts)

	c.JSON(http.StatusOK, posts)
}

// GetBlogPost vraća specifičan blog post
func (h *BlogHandler) GetBlogPost(c *gin.Context) {
	blogIDStr := c.Param("id")
	blogID, err := strconv.ParseInt(blogIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	post, err := h.store.GetBlogPostByID(blogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get blog post"})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		return
	}

	// Dobij i komentare
	comments, err := h.store.GetCommentsByBlogID(blogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post":     post,
		"comments": comments,
	})
}

// CreateComment kreira novi komentar
// SAGA PATTERN: Proverava da li korisnik prati autora pre kreiranja komentara
func (h *BlogHandler) CreateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	blogIDStr := c.Param("id")
	blogID, err := strconv.ParseInt(blogIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	// 1️⃣ SAGA STEP 1: Dohvati blog post da bi znali ko je author
	post, err := h.store.GetBlogPostByID(blogID)
	if err != nil || post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		return
	}

	// 2️⃣ SAGA STEP 2: Proveri da li korisnik prati autora (poziv ka Follower Service)
	currentUserID := userID.(int64)
	authorID := post.UserID

	// Ako korisnik komentariše svoj blog, dozvoli
	if currentUserID != authorID {
		isFollowing, err := h.checkIfUserFollowsAuthor(currentUserID, authorID)
		if err != nil {
			fmt.Printf("Error checking follower status: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify follower status"})
			return
		}

		if !isFollowing {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "You must follow the author to comment on their blog",
				"message": "SAGA validation failed: User does not follow the author",
			})
			return
		}
	}

	// 3️⃣ SAGA STEP 3: Kreiraj komentar
	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &model.BlogComment{
		BlogID:      blogID,
		UserID:      currentUserID,
		CommentText: req.CommentText,
	}

	if err := h.store.CreateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// checkIfUserFollowsAuthor proverava da li follower prati author (poziv ka Follower Service)
func (h *BlogHandler) checkIfUserFollowsAuthor(followerID, authorID int64) (bool, error) {
	url := fmt.Sprintf("%s/api/follower/check-following/%d/%d", h.followerServiceURL, followerID, authorID)

	resp, err := h.httpClient.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to call follower service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("follower service returned status: %d", resp.StatusCode)
	}

	var response struct {
		IsFollowing bool `json:"isFollowing"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, fmt.Errorf("failed to decode follower response: %w", err)
	}

	return response.IsFollowing, nil
}

// UpdateComment ažurira postojeći komentar
func (h *BlogHandler) UpdateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.store.UpdateComment(commentID, userID.(int64), req.CommentText); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
}

// DeleteComment briše komentar
func (h *BlogHandler) DeleteComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	if err := h.store.DeleteComment(commentID, userID.(int64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// LikeBlogPost lajkuje blog post
func (h *BlogHandler) LikeBlogPost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	blogIDStr := c.Param("id")
	blogID, err := strconv.ParseInt(blogIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	if err := h.store.LikeBlogPost(userID.(int64), blogID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like blog post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog post liked successfully"})
}

// UnlikeBlogPost uklanja lajk sa blog posta
func (h *BlogHandler) UnlikeBlogPost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	blogIDStr := c.Param("id")
	blogID, err := strconv.ParseInt(blogIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	if err := h.store.UnlikeBlogPost(userID.(int64), blogID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike blog post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog post unliked successfully"})
}

// Health check endpoint
func (h *BlogHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "blog-service"})
}

// UploadImage omogućava slanje slika na server
func (h *BlogHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// napravi folder ako ne postoji
	if err := os.MkdirAll("uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload dir"})
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	path := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		//"url": "/uploads/" + filename,
		//"url": "/api/blog/uploads/" + filename,
		"url": "http://localhost:8000/api/blog/uploads/" + filename,
	})
}
