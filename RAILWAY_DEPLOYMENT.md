# Railway Deployment Guide

## Issue After Refactoring

**Problem**: Railway build fails with `stat /app/server: directory not found`

**Cause**: Railway is trying to build from `./server` (old structure) but the code has been refactored to `./cmd/api` (new structure).

---

## Solution: Update Build Configuration

### Option 1: Using nixpacks.toml (Recommended)

The `nixpacks.toml` file has been created in the root directory with the correct build path:

```toml
[phases.setup]
nixPkgs = ["go"]

[phases.install]
cmds = ["go mod download"]

[phases.build]
cmds = ["go build -o main ./cmd/api"]

[start]
cmd = "./main"
```

**Railway will automatically detect and use this file.**

---

### Option 2: Using Dockerfile

A `Dockerfile` has been created for more control:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 3000
CMD ["./main"]
```

**To use Dockerfile in Railway:**
1. Go to your Railway project settings
2. Under "Build", select "Dockerfile" as the builder
3. Railway will use the Dockerfile instead of Nixpacks

---

### Option 3: Railway Dashboard Settings

If you prefer not to use config files, update in Railway dashboard:

1. Go to your Railway project
2. Click on your service
3. Go to **Settings** → **Build**
4. Set **Build Command**: `go build -o main ./cmd/api`
5. Set **Start Command**: `./main`
6. Click **Save**

---

## Environment Variables

Make sure these are set in Railway:

### Required Variables

```bash
DATABASE_URL=postgresql://user:pass@host:5432/dbname
JWT_SECRET=your-secret-key-minimum-32-characters-required
PORT=3000
```

### Optional Variables

```bash
ENV=production
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
DB_MAX_CONNECTIONS=20
DB_IDLE_TIMEOUT=300s
JWT_EXPIRATION=24h
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
```

---

## Deployment Steps

### 1. Push Changes to Git

```bash
git add .
git commit -m "Fix Railway deployment after refactoring"
git push origin main
```

### 2. Railway Auto-Deploy

Railway will automatically:
1. Detect the `nixpacks.toml` file
2. Build using `./cmd/api` instead of `./server`
3. Deploy the new version

### 3. Verify Deployment

Check the deployment logs in Railway dashboard:

```
✅ Build successful
✅ Starting application
✅ Server listening on :3000
```

### 4. Test Endpoints

```bash
# Health check
curl https://your-app.railway.app/health

# Pokemon endpoint
curl https://your-app.railway.app/pokemon?name=pikachu

# Auth endpoint (should return 400 for empty body)
curl -X POST https://your-app.railway.app/api/auth/login
```

---

## Troubleshooting

### Build Still Fails

**Check Railway Logs:**
1. Go to Railway dashboard
2. Click on your service
3. View "Deployments" tab
4. Click on the failed deployment
5. Check build logs

**Common Issues:**

1. **Old cache**: Railway might be using cached build
   - Solution: Trigger a new deployment or clear cache

2. **Wrong builder**: Railway using wrong build method
   - Solution: Explicitly set builder in settings

3. **Missing files**: Some files not committed to git
   - Solution: Check `.gitignore` and commit necessary files

### Application Crashes on Start

**Check Runtime Logs:**

1. **Database connection fails**
   - Verify `DATABASE_URL` is set correctly
   - Check database is accessible from Railway

2. **JWT secret missing**
   - Verify `JWT_SECRET` is set (min 32 chars)

3. **Port binding issues**
   - Railway automatically sets `PORT` env var
   - Your app should use `os.Getenv("PORT")` (already implemented)

### Database Connection Issues

**If using Railway PostgreSQL:**

1. Go to your PostgreSQL service in Railway
2. Copy the `DATABASE_URL` from "Connect" tab
3. Add it to your API service environment variables

**Format:**
```
postgresql://postgres:password@host.railway.internal:5432/railway
```

---

## Files Created

✅ `nixpacks.toml` - Railway Nixpacks configuration (recommended)
✅ `Dockerfile` - Alternative Docker build (optional)
✅ `railway.toml` - Railway-specific settings (optional)

**Railway will automatically use `nixpacks.toml` if present.**

---

## Quick Fix Checklist

- [ ] `nixpacks.toml` file exists in root
- [ ] Build command points to `./cmd/api`
- [ ] Environment variables are set in Railway
- [ ] `DATABASE_URL` is correct
- [ ] `JWT_SECRET` is set (32+ chars)
- [ ] Code is pushed to git
- [ ] Railway deployment triggered

---

## Verification

After deployment, your API should be accessible at:

```
https://your-app.railway.app
```

**Test endpoints:**
- `GET /health` → `{"status":"healthy"}`
- `GET /pokemon?name=pikachu` → Pokemon data
- `POST /api/auth/register` → Registration endpoint
- `POST /api/battle/start` → Battle endpoint

---

## Next Steps

1. ✅ Push code with `nixpacks.toml`
2. ✅ Wait for Railway to auto-deploy
3. ✅ Verify deployment in Railway dashboard
4. ✅ Test API endpoints
5. ✅ Update frontend to use new Railway URL (if changed)

---

## Support

If issues persist:

1. Check Railway build logs
2. Check Railway runtime logs
3. Verify all environment variables
4. Test locally: `go build -o main ./cmd/api && ./main`
5. Check Railway community forum or Discord

---

**Your deployment should work now!** 🚀

The key change is: `./server` → `./cmd/api`
