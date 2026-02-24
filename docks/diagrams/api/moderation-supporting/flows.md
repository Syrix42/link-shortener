# Admin Moderation Flows

## Link Deletion (POST)

```mermaid
sequenceDiagram
    autonumber
    actor Admin as Admin Client
    participant Router
    participant AuthMiddleware as Auth Middleware
    participant AuthorizationMiddleware as Authorization Middleware
    participant Handler
    participant BusinessLogic as Business Logic
    participant Repository
    participant DB

    Admin->>Router: Request delete link
    Router->>AuthMiddleware: Authenticate request
    AuthMiddleware-->>Router: Authenticated context
    Router->>AuthorizationMiddleware: Check admin permission

    alt Not authorized
        AuthorizationMiddleware-->>Admin: 403 Forbidden
    else Authorized
        AuthorizationMiddleware-->>Router: Permission granted
        Router->>Handler: Route request
        Handler->>Handler: Validate path/input
        Handler->>BusinessLogic: Pass delete request

        BusinessLogic->>BusinessLogic: Apply deletion rules
        BusinessLogic->>Repository: Delete or soft-delete target link
        Repository->>DB: Persist deletion state
        DB-->>Repository: Done
        Repository-->>BusinessLogic: Success

        BusinessLogic-->>Handler: Deletion success
        Handler-->>Admin: 200 OK / 204 No Content
    end
```

## GET USER LIST (GET)

```mermaid
sequenceDiagram
    autonumber
    actor Admin as Admin Client
    participant Router
    participant AuthMiddleware as Auth Middleware
    participant AuthorizationMiddleware as Authorization Middleware
    participant Handler
    participant BusinessLogic as Business Logic
    participant Repository
    participant DB

    Admin->>Router: Request delete/ban user
    Router->>AuthMiddleware: Authenticate request
    AuthMiddleware-->>Router: Authenticated context
    Router->>AuthorizationMiddleware: Check admin permission

    alt Not authorized
        AuthorizationMiddleware-->>Admin: 403 Forbidden
    else Authorized
        AuthorizationMiddleware-->>Router: Permission granted
        Router->>Handler: Route request
        Handler->>Handler: Validate path/input
        Handler->>BusinessLogic: Pass delete request

        BusinessLogic->>BusinessLogic: Apply deletion/ban rules
        BusinessLogic->>Repository: Delete or soft-delete target user
        Repository->>DB: Persist deletion state
        DB-->>Repository: Done
        Repository-->>BusinessLogic: Success

        BusinessLogic-->>Handler: Deletion success
        Handler-->>Admin: 200 OK / 204 No Content
    end
```

## GET Link List + Metadata (GET)


```mermaid
sequenceDiagram
    autonumber
    actor Admin as Admin Client
    participant Router
    participant AuthMiddleware as Auth Middleware
    participant AuthorizationMiddleware as Authorization Middleware
    participant Handler
    participant BusinessLogic as Business Logic
    participant Repository
    participant DB

    Admin->>Router: Request links list (with pagination)
    Router->>AuthMiddleware: Authenticate request
    AuthMiddleware-->>Router: Authenticated context
    Router->>AuthorizationMiddleware: Check admin permission

    alt Not authorized
        AuthorizationMiddleware-->>Admin: 403 Forbidden
    else Authorized
        AuthorizationMiddleware-->>Router: Permission granted
        Router->>Handler: Route request
        Handler->>Handler: Validate query params
        Handler->>BusinessLogic: Pass pagination/filter input

        BusinessLogic->>Repository: Fetch paginated links
        Repository->>DB: Query paginated link rows
        DB-->>Repository: Link rows
        Repository-->>BusinessLogic: Items

        BusinessLogic->>Repository: Fetch total count
        Repository->>DB: Query total count
        DB-->>Repository: Count
        Repository-->>BusinessLogic: Total

        BusinessLogic-->>Handler: List result + pagination metadata
        Handler-->>Admin: 200 OK
    end
 ```
 
 
