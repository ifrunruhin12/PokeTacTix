# PokeTacTix Deployment Guide

Complete guide to deploy PokeTacTix to production using Neon (Database), Railway (Backend), and Netlify (Frontend).

## Overview

- **Database**: Neon PostgreSQL (Free tier)
- **Backend**: Railway (Go API)
- **Frontend**: Netlify (React + Vite)

## Prerequisites

- GitHub account
- Neon account (https://neon.tech)
- Railway account (https://railway.app)
- Netlify account (https://netlify.com)
- Your code pushed to a GitHub repository

---

## Part 1: Database Setup (Neon)

### Step 1: Create Neon Project

1. Go to https://neon.tech and sign in
2. Click "Create a project"
3. Configure your project:
   - **Name**: `poketactix-db` (or your preferred name)
   - **Region**: Choose closest to your users
   - **PostgreSQL version**: 16 (latest)
4. Click "Create project"

### Step 2: Get Connection String

1. After project creation, you'll see the connection details
2. Copy the **Connection string** - it looks like:
   ```
   postgresql://username:password@ep-xxx-xxx.region.aws.neon.tech/dbname?sslmode=require
   ```
3. Save this - you'll need it for Railway

### Step 3: Run Database Migrations

You have two options:

**Option A: Using local psql**
```bash
# Install psql if you don't have it
# On macOS: brew install postgresql
# On Ubuntu: sudo apt-get install postgresql-client

# Connect to Neon database
psql "postgresql://username:password@ep-xxx-xxx.region.aws.neon.tech/dbname?sslmode=require"

# Run migrations manually by copying content from:
# internal/database/migrations/*.up.sql
```

**Option B: Let Railway run migrations on first deploy**
- Railway will automatically run migrations when the backend starts
- The migrations are in `internal/database/migrations/`

### Step 4: Neon Free Tier Limits

Be aware of free tier limits:
- **Storage**: 0.5 GB
- **Compute**: 191.9 hours/month (always-on)
- **Branches**: 10
- **Projects**: 1

**Tips**:
- Monitor usage in Neon dashboard
- Database auto-pauses after 5 minutes of inactivity (free tier)
- First query after pause may be slower (cold start)

---

## Part 2: Backend Setup (Railway)

### Step 1: Create Railway Project

1. Go to https://railway.app and sign in
2. Click "New Project"
3. Select "Deploy from GitHub repo"
4. Authorize Railway to access your GitHub
5. Select your `PokeTacTix` repository

### Step 2: Configure Build Settings

Railway should auto-detect your Go app, but verify:

1. In your Railway project, go to **Settings**
2. Check **Build Configuration**:
   - **Build Command**: (leave empty, Railway auto-detects)
   - **Start Command**: `./cmd/api/main`
   - **Root Directory**: `/` (root of repo)

### Step 3: Set Environment Variables

In Railway project, go to **Variables** tab and add:

```bash
# Database
DATABASE_URL=postgresql://username:password@ep-xxx-xxx.region.aws.neon.tech/dbname?sslmode=require

# JWT Secret (generate a random string)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server Configuration
PORT=8080
ENVIRONMENT=production

# CORS (you'll update this after deploying frontend)
CORS_ORIGINS=https://your-app.netlify.app

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
```

**Important**: 
- Generate a strong JWT_SECRET: `openssl rand -base64 32`
- You'll update `CORS_ORIGINS` after deploying to Netlify

### Step 4: Deploy

1. Railway will automatically deploy after you set variables
2. Wait for deployment to complete (2-5 minutes)
3. Once deployed, you'll get a URL like: `https://your-app.up.railway.app`
4. Test the health endpoint: `https://your-app.up.railway.app/health`

### Step 5: Verify Database Connection

Check Railway logs to ensure:
- Database connection successful
- Migrations ran successfully
- No connection errors

### Step 6: Set Up Custom Domain (Optional)

1. In Railway project, go to **Settings** → **Domains**
2. Click "Generate Domain" or add your custom domain
3. Follow DNS configuration instructions if using custom domain

---

## Part 3: Frontend Setup (Netlify)

### Step 1: Create Netlify Configuration

First, create a `netlify.toml` file in your project root:

```toml
[build]
  base = "frontend"
  command = "npm run build"
  publish = "dist"

[[redirects]]
  from = "/*"
  to = "/index.html"
  status = 200

[build.environment]
  NODE_VERSION = "20"
```

Also create `frontend/_redirects` file:

```
/*    /index.html   200
```

### Step 2: Update API URL

Create `frontend/.env.production`:

```bash
VITE_API_URL=https://your-app.up.railway.app
```

Replace `your-app.up.railway.app` with your actual Railway URL.

### Step 3: Deploy to Netlify

1. Go to https://netlify.com and sign in
2. Click "Add new site" → "Import an existing project"
3. Choose "Deploy with GitHub"
4. Authorize Netlify and select your repository
5. Configure build settings:
   - **Base directory**: `frontend`
   - **Build command**: `npm run build`
   - **Publish directory**: `frontend/dist`
6. Click "Deploy site"

### Step 4: Set Environment Variables in Netlify

1. Go to **Site settings** → **Environment variables**
2. Add:
   ```
   VITE_API_URL=https://your-app.up.railway.app
   ```

### Step 5: Update CORS in Railway

Now that you have your Netlify URL:

1. Go back to Railway
2. Update the `CORS_ORIGINS` environment variable:
   ```
   CORS_ORIGINS=https://your-app.netlify.app
   ```
3. Replace with your actual Netlify URL
4. Railway will automatically redeploy

### Step 6: Set Up Custom Domain (Optional)

1. In Netlify, go to **Domain settings**
2. Click "Add custom domain"
3. Follow DNS configuration instructions
4. Netlify provides free SSL certificates

---

## Part 4: Verification & Testing

### Test the Full Stack

1. **Visit your Netlify URL**: `https://your-app.netlify.app`
2. **Register a new account**
3. **Test key features**:
   - Login/Logout
   - View Dashboard
   - Start a battle
   - Visit the shop
   - Check profile stats

### Common Issues & Solutions

#### Issue: CORS Errors

**Symptom**: Frontend can't connect to backend, browser console shows CORS errors

**Solution**:
1. Verify `CORS_ORIGINS` in Railway matches your Netlify URL exactly
2. Make sure Railway redeployed after changing the variable
3. Check Railway logs for CORS-related messages

#### Issue: Database Connection Fails

**Symptom**: Railway logs show "failed to connect to database"

**Solution**:
1. Verify `DATABASE_URL` is correct in Railway
2. Check Neon dashboard - database might be paused (free tier)
3. Ensure connection string includes `?sslmode=require`

#### Issue: Frontend Shows "Network Error"

**Symptom**: Frontend loads but API calls fail

**Solution**:
1. Check `VITE_API_URL` in Netlify environment variables
2. Verify Railway backend is running (check Railway dashboard)
3. Test backend directly: `curl https://your-app.up.railway.app/health`

#### Issue: Migrations Not Running

**Symptom**: Backend starts but database tables don't exist

**Solution**:
1. Check Railway logs for migration errors
2. Manually run migrations using psql (see Part 1, Step 3)
3. Verify migration files are in the repository

#### Issue: Assets Not Loading (wallpaper, pokeball)

**Symptom**: Homepage background or pokeball images missing

**Solution**:
1. Verify `frontend/public/assets/` folder is in your repository
2. Check Netlify deploy logs - ensure assets were copied
3. Images should be at `/assets/wallpaper.jpg` and `/assets/pokeball.png`

---

## Part 5: Monitoring & Maintenance

### Monitor Your Services

**Neon Dashboard**:
- Check database storage usage
- Monitor query performance
- View connection stats

**Railway Dashboard**:
- Monitor CPU and memory usage
- Check deployment logs
- View request metrics

**Netlify Dashboard**:
- Monitor build times
- Check bandwidth usage
- View deploy logs

### Free Tier Limits Summary

| Service | Limit | What Happens When Exceeded |
|---------|-------|---------------------------|
| Neon | 0.5 GB storage | Need to upgrade |
| Neon | 191.9 hours/month | Database pauses, need to upgrade |
| Railway | $5 credit/month | Service stops, need to add payment |
| Netlify | 100 GB bandwidth/month | Site stops serving, need to upgrade |

### Backup Strategy

**Database Backups**:
- Neon automatically backs up your database
- Free tier: 7-day retention
- Can restore from Neon dashboard

**Code Backups**:
- Your code is on GitHub (already backed up)
- Tag releases: `git tag v1.0.0 && git push --tags`

---

## Part 6: Environment Variables Reference

### Railway (Backend)

```bash
# Required
DATABASE_URL=postgresql://...
JWT_SECRET=your-secret-key
PORT=8080
ENVIRONMENT=production
CORS_ORIGINS=https://your-app.netlify.app

# Optional
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
LOG_LEVEL=info
```

### Netlify (Frontend)

```bash
# Required
VITE_API_URL=https://your-app.up.railway.app
```

---

## Part 7: Updating Your Deployment

### Update Backend

1. Push changes to GitHub
2. Railway automatically detects and redeploys
3. Monitor Railway logs during deployment
4. Verify health endpoint after deployment

### Update Frontend

1. Push changes to GitHub
2. Netlify automatically detects and rebuilds
3. Monitor Netlify deploy logs
4. Clear browser cache and test

### Database Migrations

When adding new migrations:

1. Create migration files in `internal/database/migrations/`
2. Push to GitHub
3. Railway will run new migrations on next deploy
4. Verify in Railway logs that migrations succeeded

---

## Part 8: Cost Optimization Tips

1. **Neon**: Database auto-pauses after 5 min inactivity (free tier)
2. **Railway**: Monitor usage, $5/month credit should be enough for low traffic
3. **Netlify**: Optimize images and assets to reduce bandwidth
4. **Caching**: Implement caching to reduce database queries

---

## Part 9: Security Checklist

- [ ] Strong JWT_SECRET set in Railway
- [ ] CORS_ORIGINS properly configured
- [ ] Database connection uses SSL (`?sslmode=require`)
- [ ] Rate limiting enabled
- [ ] No sensitive data in GitHub repository
- [ ] Environment variables not hardcoded
- [ ] HTTPS enabled (automatic on Railway and Netlify)

---

## Part 10: Next Steps

After successful deployment:

1. **Set up monitoring**: Consider adding error tracking (Sentry, etc.)
2. **Analytics**: Add Google Analytics or similar
3. **Custom domains**: Set up your own domain names
4. **CI/CD**: Already set up via GitHub integration
5. **Testing**: Set up automated tests in CI pipeline

---

## Support & Resources

- **Neon Docs**: https://neon.tech/docs
- **Railway Docs**: https://docs.railway.app
- **Netlify Docs**: https://docs.netlify.com
- **Your Repository**: Check README.md for local development

---

## Quick Reference Commands

```bash
# Test backend locally
go run cmd/api/main.go

# Test frontend locally
cd frontend && npm run dev

# Build frontend
cd frontend && npm run build

# Generate JWT secret
openssl rand -base64 32

# Connect to Neon database
psql "postgresql://..."

# Check Railway logs
# Use Railway dashboard or CLI: railway logs

# Check Netlify logs
# Use Netlify dashboard
```

---

**Congratulations!** 🎉 Your PokeTacTix game is now live in production!
