# Roadmap

## v0.1

### Distribution

- Create Docker Image
- Create Docker Compose file
- Create Binary

### Server

- Configuration
    - Address and Port

- Storage
    - In-Memory

- API
    - OpenAPI 3.0 spec
        - https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#impl-gin
    - `gin` framework

    - User registration and authentication
        - User entity
            - ID
            - Username
            - Password
            - IsAdmin
            - CreatedAt
            - LastUpdatedAt
        - Endpoints
            - `POST /api/init`
            - `GET /api/users`
            - `GET /api/users/<id>`
            - `POST /api/users/register`
                - Password hashing
            - `POST /api/users/login` 
            - `PATCH /api/users/<id>`
            - `DELETE /api/users/<id>`


## v0.2

### Server

- Configuration
    - Add support for DB backends:
        - SQLite

- WebSockets
    - One-on-one messaging
    - `gorilla/websocket`
        - https://github.com/gorilla/websocket

### Client

- Implement Configuration

- Create Commands
    - admin
        - init
        - users
            - list
            - create
            - update
            - delete


## v1.0

### CI/CD

- Automatic binary build
- Automatic docker hub deployment
- Automatic tests
- Automatic linter

### Distribution

- Provide Docker Image
- Provide Docker Compose file
- Provide binary

### Server

- User registration and authentication
    - JWT authentication / authorization
- User roles
- One-on-one messaging
- Group chats
- End-to-end encryption
- User status
- File sharing
- Configuration file
    - DB Backends
    - Address and Port
    - History Retention Policy
    - File Retention Policy
    - Free Registration or Invitation Only

### Clients

- CLI
    - Configuration
        - File `~/.config/hoodie.yaml`
        - Environment Variables
            - `HOODIE_ADDRESS`
            - `HOODIE_PORT`
        - Global CLI Options
            - `-a ADDR`
            - `-p PORT`
    - Commands
        - admin
            - init
            - users
                - list
                - create
                - update
                - delete
                - invite


