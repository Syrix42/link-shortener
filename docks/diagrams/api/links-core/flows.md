

## Dashboard / List My Links (Authenticated) GET /api/links
```mermaid
sequenceDiagram
    autonumber
    actor User as Authenticated User
    participant Router
    participant AuthMiddleware as Auth Middleware
    participant Handler
    participant BusinessLogic as Link Service / Business Logic
    participant Repository
    participant DB

    User->>Router: GET /api/links
    Router->>AuthMiddleware: Authenticate token

    alt Invalid / missing token
        AuthMiddleware-->>User: 401 Unauthorized
    else Authenticated
        AuthMiddleware-->>Router: Authenticated user context
        Router->>Handler: Route request
        Handler->>Handler: Validate query params (page, limit) [optional]
        Handler->>BusinessLogic: GetUserLinks(userId)

        BusinessLogic->>Repository: Fetch links by owner (not deleted)
        Repository->>DB: SELECT links WHERE user_id = ? AND deleted_at IS NULL
        DB-->>Repository: Link rows
        Repository-->>BusinessLogic: Links data

        BusinessLogic-->>Handler: Dashboard link list (metadata)
        Handler-->>User: 200 OK (links + clickCount + lastClickedAt)
    end
```

## Create Short Link (Authenticated) — POST /api/links

```mermaid 
sequenceDiagram
    autonumber
    actor User as Authenticated User
    participant Router
    participant AuthMiddleware as Auth Middleware
    participant Handler
    participant BusinessLogic as Link Service / Business Logic
    participant Repository
    participant DB

    User->>Router: POST /api/links { longUrl, shortCode }
    Router->>AuthMiddleware: Authenticate token

    alt Invalid / missing token
        AuthMiddleware-->>User: 401 Unauthorized
    else Authenticated
        AuthMiddleware-->>Router: Authenticated user context
        Router->>Handler: Route request
        Handler->>Handler: Validate body (URL, shortCode format)
        Handler->>BusinessLogic: CreateLink(userId, longUrl, shortCode)

        BusinessLogic->>BusinessLogic: Check reserved shortCode
        BusinessLogic->>Repository: Check shortCode uniqueness
        Repository->>DB: SELECT by short_code
        DB-->>Repository: Exists / Not found
        Repository-->>BusinessLogic: Result

        alt shortCode already exists
            BusinessLogic-->>Handler: Conflict error
            Handler-->>User: 409 Conflict
        else Valid and unique
            BusinessLogic->>Repository: Insert link
            Repository->>DB: INSERT links row
            DB-->>Repository: Created row
            Repository-->>BusinessLogic: Link created
            BusinessLogic-->>Handler: Created link DTO
            Handler-->>User: 201 Created (shortUrl, metadata)
        end
    end
```

## Delete Short Link (Authenticated, Owner Only) — DELETE /api/links/:shortCode
```mermaid 
sequenceDiagram
    autonumber
    actor User as Authenticated User
    participant Router
    participant AuthMiddleware as Auth Middleware
    participant Handler
    participant BusinessLogic as Link Service / Business Logic
    participant Repository
    participant DB

    User->>Router: DELETE /api/links/{shortCode}
    Router->>AuthMiddleware: Authenticate token

    alt Invalid / missing token
        AuthMiddleware-->>User: 401 Unauthorized
    else Authenticated
        AuthMiddleware-->>Router: Authenticated user context
        Router->>Handler: Route request
        Handler->>Handler: Validate path param shortCode
        Handler->>BusinessLogic: DeleteLink(userId, shortCode)

        BusinessLogic->>Repository: Find link by shortCode
        Repository->>DB: SELECT link by short_code
        DB-->>Repository: Link row / Not found
        Repository-->>BusinessLogic: Result

        alt Link not found or already deleted
            BusinessLogic-->>Handler: Not found
            Handler-->>User: 404 Not Found
        else Link exists
            BusinessLogic->>BusinessLogic: Check ownership (ownerUserId == userId) or admin
            alt Not owner and not admin
                BusinessLogic-->>Handler: Forbidden
                Handler-->>User: 403 Forbidden
            else Authorized
                BusinessLogic->>Repository: Soft delete link
                Repository->>DB: UPDATE links SET deleted_at = NOW()
                DB-->>Repository: Updated
                Repository-->>BusinessLogic: Success
                BusinessLogic-->>Handler: Deleted
                Handler-->>User: 204 No Content
            end
        end
    end
```

## Public Redirect (No Auth) — GET /:shortCode

```mermaid
sequenceDiagram
    autonumber
    actor Visitor as Public Visitor
    participant Router
    participant Handler
    participant BusinessLogic as Redirect Service / Business Logic
    participant Repository
    participant DB

    Visitor->>Router: GET /{shortCode}
    Router->>Handler: Route request
    Handler->>Handler: Validate shortCode format
    Handler->>BusinessLogic: ResolveAndRedirect(shortCode)

    BusinessLogic->>Repository: Find active link by shortCode
    Repository->>DB: SELECT long_url, click_count, deleted_at WHERE short_code = ?
    DB-->>Repository: Link row / Not found
    Repository-->>BusinessLogic: Result

    alt Not found / deleted
        BusinessLogic-->>Handler: Link not redirectable
        Handler-->>Visitor: 404 Not Found
    else Found and active
        BusinessLogic->>Repository: Increment click_count and set last_clicked_at
        Repository->>DB: UPDATE links SET click_count = click_count + 1, last_clicked_at = NOW()
        DB-->>Repository: Updated
        Repository-->>BusinessLogic: Success

        BusinessLogic-->>Handler: Resolved longUrl
        Handler-->>Visitor: 302 Redirect (Location: longUrl)
    end
```
