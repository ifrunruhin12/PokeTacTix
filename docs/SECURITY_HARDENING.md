# Security Hardening Implementation

This document describes the security improvements implemented in Task 9 of the PokeTacTix web application.

## Overview

All security hardening tasks have been completed to protect the application from common vulnerabilities and attacks.

## 1. Input Validation and Sanitization (Task 9.1)

### Implementation

Created comprehensive validation middleware in `internal/middleware/validation.go`:

#### Features
- **String Sanitization**: HTML escaping to prevent XSS attacks
- **Username Validation**: 3-50 characters, alphanumeric + underscores only
- **Email Validation**: Proper email format validation with regex
- **Password Validation**: Enforces strong password requirements:
  - Minimum 8 characters
  - At least 1 uppercase letter
  - At least 1 lowercase letter
  - At least 1 number
  - At least 1 special character
- **Pokemon Name Validation**: Letters, hyphens, spaces, and apostrophes only
- **Battle Mode Validation**: Only "1v1" or "5v5" allowed
- **Battle Move Validation**: Only valid moves (attack, defend, pass, sacrifice, surrender)
- **Card ID Validation**: Positive integers, no duplicates

#### Applied To
- Auth handlers: Registration and login inputs
- Battle handlers: Battle mode, moves, and battle IDs
- Shop handlers: Pokemon names
- All user-facing endpoints

### Benefits
- Prevents XSS attacks through HTML escaping
- Prevents injection attacks through strict input validation
- Provides clear error messages for invalid inputs
- Consistent validation across all endpoints

## 2. SQL Injection Prevention (Task 9.2)

### Implementation

Verified and fixed all database queries to use parameterized queries exclusively.

#### Key Changes

**Fixed in `internal/stats/repository.go`:**
- Removed `fmt.Sprintf` query building in `updatePlayerStats()`
- Replaced with separate parameterized queries for each mode/result combination
- Added strict validation of mode and result parameters before query execution

**Verified in all repositories:**
- `internal/auth/repository.go`: ✅ All queries use `$1, $2, ...` parameters
- `internal/cards/repository.go`: ✅ All queries use parameterized binding
- `internal/shop/repository.go`: ✅ All queries use pgx parameter binding
- `internal/stats/repository.go`: ✅ Fixed to use parameterized queries only

### Benefits
- Complete protection against SQL injection attacks
- No user input is ever concatenated into SQL strings
- All queries use pgx parameter binding ($1, $2, etc.)
- Strict validation before query execution

## 3. CORS Configuration (Task 9.3)

### Implementation

Updated CORS configuration in `cmd/api/main.go`:

#### Production Mode
- Only allows explicitly configured origins from environment variables
- If no origins configured, blocks all cross-origin requests
- Logs warning if no origins are set in production

#### Development Mode
- Allows localhost origins by default (ports 3000, 5173, 8080)
- Can be overridden with environment variables
- Logs configured origins for debugging

#### Configuration
```go
corsConfig := cors.Config{
    AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
    AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
    AllowCredentials: true,
    MaxAge:           3600,
}
```

### Benefits
- Prevents unauthorized cross-origin requests in production
- Allows credentials (cookies, authorization headers)
- Flexible configuration via environment variables
- Secure defaults for both development and production

## 4. Security Headers and HTTPS Enforcement (Task 9.4)

### Implementation

Created security middleware in `internal/middleware/security.go`:

#### Security Headers Added

1. **X-Frame-Options: DENY**
   - Prevents clickjacking attacks
   - Blocks embedding in iframes

2. **X-Content-Type-Options: nosniff**
   - Prevents MIME type sniffing
   - Forces browser to respect declared content types

3. **Referrer-Policy: strict-origin-when-cross-origin**
   - Controls referrer information leakage
   - Only sends origin for cross-origin requests

4. **X-XSS-Protection: 1; mode=block**
   - Enables XSS filter in legacy browsers
   - Blocks page if XSS detected

5. **Content-Security-Policy**
   - Restricts resource loading
   - Prevents inline script execution (with exceptions for compatibility)
   - Blocks framing

6. **Strict-Transport-Security** (HTTPS only)
   - Forces HTTPS for 1 year
   - Includes subdomains
   - Only set when using HTTPS

#### HTTPS Redirect

- Automatically redirects HTTP to HTTPS in production
- Checks both protocol and X-Forwarded-Proto header (for proxies)
- Only active in production environment

### Benefits
- Comprehensive protection against common web vulnerabilities
- Prevents clickjacking, XSS, and MIME sniffing attacks
- Forces HTTPS in production
- Compatible with reverse proxies and load balancers

## Environment Configuration

### Required Environment Variables

```bash
# Server
ENV=production  # or "development"
PORT=3000

# CORS
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# Database (already configured)
DATABASE_URL=postgresql://...

# JWT (already configured)
JWT_SECRET=...
```

### Development vs Production

| Feature | Development | Production |
|---------|------------|------------|
| CORS Origins | localhost:3000,5173,8080 | Explicit list required |
| HTTPS Redirect | Disabled | Enabled |
| Security Headers | All enabled | All enabled |
| Input Validation | Enabled | Enabled |
| SQL Injection Prevention | Enabled | Enabled |

## Testing Recommendations

### Manual Testing

1. **Input Validation**
   - Try registering with weak passwords
   - Try SQL injection in username/email fields
   - Try XSS payloads in text inputs

2. **CORS**
   - Test from unauthorized origins
   - Verify credentials are sent correctly

3. **Security Headers**
   - Check response headers with browser dev tools
   - Verify CSP blocks unauthorized scripts

4. **HTTPS Redirect**
   - Test HTTP requests in production
   - Verify redirect to HTTPS

### Automated Testing

Consider adding security tests:
- SQL injection test suite
- XSS payload test suite
- CORS policy tests
- Header verification tests

## Security Checklist

- [x] Input validation on all endpoints
- [x] XSS prevention through HTML escaping
- [x] SQL injection prevention with parameterized queries
- [x] CORS properly configured for production
- [x] Security headers on all responses
- [x] HTTPS enforcement in production
- [x] Password strength requirements
- [x] Email format validation
- [x] Battle move validation
- [x] Pokemon name sanitization

## Future Improvements

1. **Rate Limiting**: Already implemented in auth, consider adding to other endpoints
2. **Request Size Limits**: Add max body size limits
3. **API Key Authentication**: For service-to-service communication
4. **Audit Logging**: Log all security-relevant events
5. **Penetration Testing**: Regular security audits
6. **Dependency Scanning**: Automated vulnerability scanning
7. **WAF Integration**: Web Application Firewall for additional protection

## References

- OWASP Top 10: https://owasp.org/www-project-top-ten/
- OWASP Cheat Sheet Series: https://cheatsheetseries.owasp.org/
- Go Security Best Practices: https://golang.org/doc/security/
- Fiber Security: https://docs.gofiber.io/guide/security

## Compliance

This implementation addresses the following security requirements:
- Requirement 14.1: Password hashing and validation
- Requirement 14.2: SQL injection prevention
- Requirement 14.3: JWT security (already implemented)
- Requirement 14.4: Input validation and XSS prevention
- Requirement 14.5: Rate limiting (already implemented)
- Requirement 14.6: CORS configuration
- Requirement 14.7: Security headers
- Requirement 14.8: HTTPS enforcement
