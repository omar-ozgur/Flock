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
- **[Search for user]** POST /users/search (Takes: user; Returns: status, message, users)
- **[Get user]** GET /users/{id} (Returns: status, message, user)
- **[Update user]** PUT /users/{id} (Takes: token, user; Returns: status, message, user)
- **[Delete user]** DELETE /users/{id} (Takes: token; Returns: status, message)
- **[Get events a user is attending]** GET /users/{id}/attendance (Returns: status, message, events)
- **[Get all events]** GET /events (Returns: status, message, events)
- **[Create event]** POST /events (Takes: token, event; Returns: status, message, event)
- **[Search events]** POST /events/search (Takes: event; Returns: status, message, events)
- **[Get event]** GET /events/{id} (Returns: status, message, event)
- **[Update event]** PUT /events/{id} (Takes: token, event; Returns: status, message, event}
- **[Delete event]** DELETE /events/{id} (Takes: token; Returns: status, message)
- **[Get event attendees]** GET /events/{id}/attendees (Returns: status, message, attendees)
- **[Attend event]** POST /events/{id}/attend (Takes: token; Returns: status, message, attendee)
- **[Remove event attendance]** DELETE /events/{id}/attendance (Takes: token; Returns: status, message)

# Models
```
type User struct {
  Id           int       `valid:"-"`
  First_name   string    `valid:"required"`
  Last_name    string    `valid:"required"`
  Email        string    `valid:"email,required"`
  Fb_id        string    `valid:"-"`
  Password     []byte    `valid:"required"`
  Time_created time.Time `valid:"-"`
}
```

```
type Event struct {
  Id           int       `valid:"-"`
  Title        string    `valid:"required"`
  Description  string    `valid:"required"`
  Location     string    `valid:"required"`
  User_id      int       `valid:"-"`
  Latitude     string    `valid:"latitude,required"`
  Longitude    string    `valid:"longitude,required"`
  Zip          int       `valid:"required"`
  Time_created time.Time `valid:"-"`
  Time_expires time.Time `valid:"-"`
}
```

```
type Attendee struct {
  Id       int `valid:"-"`
  Event_id int `valid:"required"`
  User_id  int `valid:"required"`
}
```
