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
- **[Get all posts]** GET /posts (Returns: status, message, posts)
- **[Create post]** POST /posts (Takes: token, post; Returns: status, message, post)
- **[Search posts]** POST /posts/search (Takes: post; Returns: status, message, posts)
- **[Get post]** GET /posts/{id} (Returns: status, message, post)
- **[Update post]** PUT /posts/{id} (Takes: token, post; Returns: status, message, post}
- **[Delete post]** DELETE /posts/{id} (Takes: token; Returns: status, message)
- **[Get post attendees]** GET /posts/{id}/attendees (Returns: status, message, attendees)
- **[Attend post]** POST /posts/{id}/attend (Takes: token; Returns: status, message, attendee)
- **[Remove post attendance]** DELETE /posts/{id}/attendance (Takes: token; Returns: status, message)
