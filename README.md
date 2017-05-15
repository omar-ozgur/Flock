# Flock-api
Join events on-the-fly

# ENV Variables
- PORT
- DB_INFO
- FLOCK_TOKEN_SECRET

# Routes
- **[Sign up]** POST /signup (Takes: user; Returns: status, message, user)
- **[Log in]** POST /login (Takes: email, password; Returns: status, message, token)
- **[Get current user]** GET /profile (Takes: token; Returns: status, message, user)
- **[Get all users]** GET /users (Returns: status, message, users)
- **[Get user]** GET /users/{id} (Returns: status, message, user)
- **[Update user]** PUT /users/{id} (Takes: token, user; Returns: status, message, user)
- **[Delete user]** DELETE /users/{id} (Takes: token; Returns: status, message)
- **[Get posts a user is attending]** GET /users/{id}/attendance (Returns: status, message, posts)
- **[Get all posts]** GET /posts (Returns: status, message, posts)
- **[Create post]** POST /posts (Takes: token, post; Returns: status, message, post)
- **[Search posts]** POST /posts/search (Takes: post; Returns: status, message, posts)
- **[Get post]** GET /posts/{id} (Returns: status, message, post)
- **[Update post]** PUT /posts/{id} (Takes: token, post; Returns: status, message, post}
- **[Delete post]** DELETE /posts/{id} (Takes: token; Returns: status, message)
- **[Get post attendees]** GET /posts/{id}/attendees (Returns: status, message, attendees)
- **[Attend post]** POST /posts/{id}/attend (Takes: token; Returns: status, message, attendee)
- **[Remove post attendance]** DELETE /posts/{id}/attendance (Takes: token; Returns: status, message)

# Models
```
type User struct {
  Id           int       `valid:"-"`
  First_name   string    `valid:"alphanum,required"`
  Last_name    string    `valid:"alphanum,required"`
  Email        string    `valid:"email,required"`
  Fb_id        int       `valid:"-"`
  Password     []byte    `valid:"required"`
  Time_created time.Time `valid:"-"`
}
```

```
type Post struct {
  Id           int       `valid:"-"`
  Title        string    `valid:"alphanum,required"`
  Location     string    `valid:"alphanum,required"`
  User_id      int       `valid:"required"`
  Latitude     string    `valid:"latitude,required"`
  Longitude    string    `valid:"longitude,required"`
  Zip          int       `valid:"required"`
  Time_created time.Time `valid:"-"`
  Time_expires time.Time `valid:"-"`
}
```

```
type Attendee struct {
  Id      int `valid:"-"`
  Post_id int `valid:"-"`
  User_id int `valid:"-"`
}
```
