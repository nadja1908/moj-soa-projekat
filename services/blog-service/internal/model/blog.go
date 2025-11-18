package model

import "time"

// BlogPost predstavlja blog post u sistemu
type BlogPost struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	LikesCount  int       `json:"likesCount"`
}

// BlogComment predstavlja komentar na blog post
type BlogComment struct {
	ID          int64     `json:"id"`
	BlogID      int64     `json:"blogId"`
	UserID      int64     `json:"userId"`
	CommentText string    `json:"commentText"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// BlogLike predstavlja lajk na blog post
type BlogLike struct {
	ID        int64     `json:"id"`
	BlogID    int64     `json:"blogId"`
	UserID    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

// BlogImage predstavlja sliku povezanu sa blog postom
type BlogImage struct {
	ID        int64     `json:"id"`
	BlogID    int64     `json:"blogId"`
	ImageURL  string    `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
}
