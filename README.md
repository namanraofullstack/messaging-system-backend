# Messaging System Backend

A messaging system backend project for intern onboarding tasks.

## 🚀 How to Run Project

### Prerequisites
- VS Code
- Docker Desktop
- Postman (for testing)

### Setup Steps

1. **Open VS Code**
   - Go to File → Open folder from workspace
   - Create new folder, name it and select it

2. **Open Terminal** (Ctrl + J or open terminal)

3. **Clone Repository**
   ```bash
   git clone https://github.com/namanraofullstack/messaging-system-backend.git
   ```

4. **Navigate to Project Directory**
   ```bash
   cd messaging-system-backend/messaging-system-backend/
   ```

5. **Configure Environment**
   - Change `.env.example` to `.env`
   - Fill in the credentials:

   ```env
   # Database Configuration
   DB_HOST=
   DB_PORT=
   DB_USER=
   DB_PASSWORD=
   DB_NAME=

   # Application Port
   PORT=

   # OpenAI API
   OPENAI_API_KEY=

   # JWT Secret Key (change in production)
   JWT_SECRET_KEY=

   # Redis Configuration
   REDIS_URL=
   REDIS_PASSWORD=
   REDIS_PORT=
   ```

6. **Run with Docker**
   - Make sure Docker Desktop is running
   ```bash
   docker compose down -v
   docker compose up --build
   ```

7. **Access Application**
   - Your app will be running on `localhost:8080`
   - Open Postman for testing

### Additional Notes
- The API for LLM summarization is from Hugging Face: https://huggingface.co/settings/tokens
- It uses a read token - generate the token when you edit the .env

## 📚 API Documentation

### 1. User Management

#### User Registration

A user should be able to register.

```bash
curl --location 'http://localhost:8080/register' \
--header 'Content-Type: application/json' \
--data '{"username":"naman", "password":"1234567"}'
```

**Success:**
```
201 Created
{"message":"registered"}
```

**Failure:**
```
400 Bad request
Invalid request

500 INTERNAL SERVER ERROR
Error saving user: pq: duplicate key value violates unique constraint "users_username_key"
```

#### User Login

A user should be able to log in.

```bash
curl --location 'http://localhost:8080/login' \
--header 'Content-Type: application/json' \
--data '{"username":"naman", "password":"1234567"}'
```

**Success:**
```
200 OK
{"token":"YOUR_TOKEN"}
```

**Failure:**
```
401 Unauthorized
Invalid username or password
```

#### User Logout

A user should be able to log out.

```bash
curl --location --request POST 'http://localhost:8080/logout' \
--header 'Authorization: Bearer <YOUR_TOKEN>'
```

**Success:**
```
200 OK
{
    "message": "Logged out successfully. Token has been blacklisted."
}
```

**Failure:**
```
500 Internal Server Error
Error blacklisting token: token signature is invalid: signature is invalid
```

#### Protected Route

A protected route to check if the user token is valid or not.

```bash
curl --location 'http://localhost:8080/protected' \
--header 'Authorization: Bearer <YOUR_TOKEN>'
```

**Success:**
```
200 OK
{"message":"Access granted to protected route."}
```

**Failure:**
```
401 Unauthorized
Invalid token
```

### 2. Peer to Peer Messaging (Direct Messages)

#### Send Message

A logged-in user can send a message to other users.

```bash
curl --location 'http://localhost:8080/send' \
--header 'Authorization: Bearer <YOUR_TOKEN>' \
--header 'Content-Type: application/json' \
--data '{
    "receiver_id": 3,
    "content": "Hey there!How are you?? What about your CET Exams??"
}'
```

**Success:**
```
200 OK
{"message":"Message sent"}
```

**Failure:**
```
401 Unauthorized 
Unauthorized
```

### 3. Group Messaging

#### Create Group

```bash
curl --location 'http://localhost:8080/group/create' \
--header 'Authorization: Bearer <YOUR_TOKEN>' \
--header 'Content-Type: application/json' \
--data '{"name": "hackers", "members": [1,2,3]}'
```

**Success:**
```
200 OK
{"group_id":5,"message":"Group created"}
```

**Failure:**
```
500 INTERNAL SERVER ERROR
group cannot have more than 25 members

401 Unauthorized
Unauthorized
```

**Note:** A User Group must have less than 25 members.

#### Send Group Message

A logged in user can post a message to a group.

```bash
curl --location 'http://localhost:8080/group/message' \
--header 'Authorization: Bearer <YOUR_TOKEN>' \
--header 'Content-Type: application/json' \
--data '{"group_id": 2, "content": "Hey developers lets finish our tasks and take it to production"}'
```

**Success:**
```
200 OK
{"message":"Message sent"}
```

**Failure:**
```
403 FORBIDDEN
You are not a member of this group

401 Unauthorized
Unauthorized
```

#### Add Member to Group

```bash
curl --location 'http://localhost:8080/group/add-member' \
--header 'Authorization: Bearer <YOUR_TOKEN>' \
--header 'Content-Type: application/json' \
--data '{"group_id": 1, "user_id": 6}'
```

**Success:**
```
200 OK
{"message":"Member added to group"}
```

**Failure:**
```
500 internal server error 
only admins can add members

500 INTERNAL SERVER ERROR
group already has 25 members

500 INTERNAL SERVER ERROR
error adding member
```

#### Remove Member from Group

```bash
curl --location 'http://localhost:8080/group/remove-member' \
--header 'Authorization: Bearer <YOUR_TOKEN>' \
--header 'Content-Type: application/json' \
--data '{
    "group_id": 4,
    "user_id": 3
}'
```

**Success:**
```
200 OK
{"message":"Member removed from group"}
```

**Failure:**
```
403 Forbidden
target user not found in group

403 forbidden
only admins can remove members

401 Unauthorized

403 forbidden
you cannot remove yourself
```

#### Promote to Admin

```bash
curl --location 'http://localhost:8080/group/promote' \
--header 'Authorization: Bearer <YOUR_TOKEN>' \
--header 'Content-Type: application/json' \
--data '{
    "group_id": 4,
    "user_id": 2
}'
```

**Success:**
```
200 OK
{"message":"Member promoted to admin"}
```

**Failure:**
```
403 Forbidden
group already has 2 admins

403 Forbidden
member not found in group

403 forbidden
only admins can promote members
```

#### Demote to Member

```bash
curl --location 'http://localhost:8080/group/demote' \
--header 'Authorization: Bearer <YOUR_TOKEN>' \
--header 'Content-Type: application/json' \
--data '{
    "group_id": 4,
    "user_id": 1
}'
```

**Success:**
```
200 OK
{"message":"Admin demoted to member"}
```

**Failure:**
```
403
user is not an admin

403
only admins can demote

403 
you cannot demote yourself

401 Unauthorized
```

**Note:** A User Group can have at max 2 admins.

### 4. Chat Management

#### Latest DM Previews

A logged in user must be able to see his latest 10 messages from 10 different users with a preview of the latest message in that specific chat.

```bash
curl --location 'http://localhost:8080/chats/latest-dm-previews' \
--header 'Authorization: Bearer <YOUR_TOKEN>' \
--header 'Content-Type: application/json'
```

**Success:**
```
200 OK
[
    {
        "id": 3,
        "sender_id": 1,
        "receiver_id": 1,
        "content": "Hey there!",
        "created_at": "2025-07-28T09:18:11.715147Z"
    },
    {
        "id": 5,
        "sender_id": 1,
        "receiver_id": 2,
        "content": "Hey there!tu tu main main",
        "created_at": "2025-07-28T09:18:27.757486Z"
    },
    {
        "id": 6,
        "sender_id": 1,
        "receiver_id": 3,
        "content": "Hey there!paisa de",
        "created_at": "2025-07-28T09:19:39.258786Z"
    }
]
```

**Failure:**
```
401 Unauthorized
```

#### Latest Group Previews

A logged in user can view the latest 10 groups (at max) – with a preview of latest message available to view in respective groups.

```bash
curl --location 'http://localhost:8080/chats/latest-group-previews' \
--header 'Authorization: Bearer <YOUR_TOKEN>' \
--header 'Content-Type: application/json'
```

**Success:**
```
200 OK
[
    {
        "id": 6,
        "name": "testing234",
        "last_message": "",
        "last_message_time": "2025-07-28T09:22:34.923136Z"
    },
    {
        "id": 7,
        "name": "testing234",
        "last_message": "",
        "last_message_time": "2025-07-28T09:22:34.923136Z"
    },
    {
        "id": 3,
        "name": "hackers",
        "last_message": "",
        "last_message_time": "2025-07-28T09:22:34.923136Z"
    },
    {
        "id": 1,
        "name": "testing",
        "last_message": "",
        "last_message_time": "2025-07-28T09:22:34.923136Z"
    },
    {
        "id": 5,
        "name": "testing234",
        "last_message": "",
        "last_message_time": "2025-07-28T09:22:34.923136Z"
    },
    {
        "id": 4,
        "name": "helloworld",
        "last_message": "Hey group!I need some money",
        "last_message_time": "2025-07-28T08:29:19.526639Z"
    },
    {
        "id": 2,
        "name": "developers",
        "last_message": "Hey developers lets finish our tasks and take it to production",
        "last_message_time": "2025-07-28T07:42:34.459401Z"
    }
]
```

**Failure:**
```
401 Unauthorized
```

#### View Chat Messages

A logged in user must be able to see the latest 10 messages when they select a specific chat (could be DM or Group).

**For DM:**
```bash
curl --location 'http://localhost:8080/chats/messages?type=dm&id=3' \
--header 'Authorization: Bearer <YOUR_TOKEN>' \
--header 'Content-Type: application/json'
```

**Success:**
```
200 OK
[
    {
        "id": 6,
        "sender_id": 1,
        "receiver_id": 3,
        "content": "Hey there!paisa de",
        "created_at": "2025-07-28T09:19:39.258786Z"
    },
    {
        "id": 2,
        "sender_id": 1,
        "receiver_id": 3,
        "content": "Hey there!give me money",
        "created_at": "2025-07-28T08:22:56.32981Z"
    },
    {
        "id": 1,
        "sender_id": 1,
        "receiver_id": 3,
        "content": "Hey there!How are you?? What about your CET Exams??",
        "created_at": "2025-07-28T07:34:01.443008Z"
    }
]
```

**Failure:**
```
400 BAD REQUEST
Invalid chat ID

401 Unauthorized
```

**For Group:**
```bash
curl --location 'http://localhost:8080/chats/messages?type=group&id=4' \
--header 'Authorization: Bearer <YOUR_TOKEN>' \
--header 'Content-Type: application/json'
```

**Success:**
```
[
    {
        "id": 3,
        "sender_id": 1,
        "group_id": 4,
        "content": "Hey group!I need some money",
        "created_at": "2025-07-28T08:29:19.526639Z"
    },
    {
        "id": 2,
        "sender_id": 1,
        "group_id": 4,
        "content": "Hey group!I need some money",
        "created_at": "2025-07-28T08:28:11.240992Z"
    }
]
```

**Failure:**
```
401 unauthorized

403 
Not a member of the group
```

### 5. AI Features

#### Group Summary

Provide an API that takes a Group_ID and returns a summary of the messages along with users who sent it using an LLM (e.g., OpenAI or any other public LLM API).

```bash
curl --location 'http://localhost:8080/groups/summary?group_id=4' \
--header 'Authorization: Bearer <YOUR_TOKEN>'
```

**Success:**
```
200 OK
{
    "summary": "Saawan needs money for his studies loan. He is angry at the group.       .   \"Foolde you all are foolish.\"  \"I need some money for my studies loan.\"   'Foalde you are foolish.'",
    "users": [
        "saawan"
    ]
}
```

**Failure:**
```
500 INTERNAL SERVER ERROR
{"error": "Failed to generate summary"}
```

### 6. Edit Messages

#### Edit Direct Messages

```bash
curl --location --request PUT 'http://localhost:8080/edit/direct' \
--header 'Authorization: Bearer <YOUR_TOKEN>' \
--header 'Content-Type: application/json' \
--data '{
    "message_id": 1,
    "new_content": "Hey there! Saawan (edited2)",
    "last_updated_at": "2025-07-28T13:08:38.699794Z"
}'
```

#### Edit Group Messages

```bash
curl --location --request PUT 'http://localhost:8080/edit/group' \
--header 'Authorization: Bearer <YOUR_TOKEN>' \
--header 'Content-Type: application/json' \
--data '{
    "message_id": 1,
    "new_content": "Edited group message",
    "last_updated_at": "2025-07-28T13:12:52.757404Z"
}'
```
---

## 🔧 System Assumptions

### 1. Authentication & Security
- Users authenticate using JWT tokens, which are signed with a secret key from the environment variable JWT_SECRET_KEY
- Passwords are securely hashed using bcrypt before storing in the database
- JWT tokens are blacklisted on logout using Redis, assuming Redis is available and properly configured

### 2. Environment and Development
- The app listens on a port defined by the PORT environment variable (default 8080)

### 3. API Design
- RESTful endpoints are used for user registration, login, logout, messaging, group management, and chat previews
- Protected routes require a valid JWT in the Authorization header
- Logout is implemented by blacklisting the JWT token, not by deleting server-side sessions

### 4. User & Group Logic
- Groups have a maximum of 25 members and up to 2 admins
- Only group admins can add/remove/promote/demote members

### 5. External Services
- Optional integration with OpenAI and Hugging Face APIs for message summarization, with API keys provided via environment variables

## 🔗 Repository

[GitHub Repository](https://github.com/namanraofullstack/messaging-system-backend)