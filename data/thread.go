package data

import "time"

type Thread struct {
	Id         int
	Uuid       string
	Topic      string
	UserId     int
	CreatedAt  time.Time
	Slug       string
	Discussion string
}

type Post struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	ThreadId  int
	CreatedAt time.Time
}

// format the CreatedAt date to display nicely on the screen
func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format("Jun 2, 2023 at 3:04pm")
}

func (post *Post) CreatedAtDate() string {
	return post.CreatedAt.Format("Jun 2, 2023 at 4:27pm")
}

func (thread Thread) NumReplies() (count int) {
	rows, err := Db.Query("SELECT count(*) - 1 FROM posts where thread_id = $1", thread.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}
	}

	rows.Close()
	return
}

// Get posts to a thread
func (thread *Thread) Posts() (posts []Post, err error) {
	rows, err := Db.Query("SELECT id, uuid, body, user_id, thread_id, created_at FROM posts where thread_id = $1", thread.Id)
	if err != nil {
		return
	}

	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}

	rows.Close()
	return
}

// Create a new thread
func (user *User) CreateThread(uuid string, topic string, slug string, discussion string) (conv Thread, err error) {
	statement := "insert into threads (uuid, topic, user_id, created_at,slug,discussion)" +
		" values ($1, $2, $3, $4,$5,$6) returning id, uuid, topic, user_id, created_at,slug,discussion"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(uuid, topic, user.Id, time.Now(), slug, discussion).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt, &conv.Slug, &conv.Discussion)
	return
}

// Create a new post to a thread
func (user *User) CreatePost(conv Thread, body string) (post Post, err error) {
	statement := "INSERT INTO posts (uuid, body, user_id, thread_id, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, body, user_id, thread_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(CreateUUID(), body, user.Id, conv.Id, time.Now()).Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
	return
}

// Get all threads in the database and returns it
func Threads() (threads []Thread, err error) {
	rows, err := Db.Query("SELECT id, uuid, topic, user_id, created_at,slug,discussion FROM threads ORDER BY created_at DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		conv := Thread{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt, &conv.Slug, &conv.Discussion); err != nil {
			return
		}
		threads = append(threads, conv)
	}
	rows.Close()
	return
}

// Get a thread by the UUID
func ThreadByUUID(uuid string) (conv Thread, err error) {
	conv = Thread{}
	err = Db.QueryRow("SELECT id, uuid, topic, user_id, created_at FROM threads WHERE uuid = $1", uuid).
		Scan(
			&conv.Id,
			&conv.Uuid,
			&conv.Topic,
			&conv.UserId,
			&conv.CreatedAt,
		)
	return
}

// Get a thread by the slug
func ThreadBySlug(slug string) (conv Thread, err error) {
	conv = Thread{}
	err = Db.QueryRow("SELECT id,uuid,topic,user_id,created_at,slug,discussion FROM threads WHERE slug=$1", slug).
		Scan(&conv.Id,
			&conv.Uuid,
			&conv.Topic,
			&conv.UserId,
			&conv.CreatedAt,
			&conv.Slug,
			&conv.Discussion,
		)
	return
}

// Get the user who started this thread
func (thread *Thread) User() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, user_name, created_at FROM users WHERE id = $1", thread.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.UserName, &user.CreatedAt)
	return
}

// Get the user who wrote the post
func (post *Post) User() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, user_name, created_at FROM users WHERE id = $1", post.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.UserName, &user.CreatedAt)
	return
}
