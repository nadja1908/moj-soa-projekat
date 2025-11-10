package model

// CreateBlogRequest je model za kreiranje blog posta
type CreateBlogRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Content     string `json:"content" binding:"required"`
}

// CreateCommentRequest je model za kreiranje komentara
type CreateCommentRequest struct {
	CommentText string `json:"commentText" binding:"required"`
}