# Quick Deploy Guide

The fastest way to get PokeTacTix deployed. For detailed instructions, see [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md).

## 1. Database (5 minutes)

```bash
# 1. Go to https://neon.tech
# 2. Create project: "poketactix-db"
# 3. Copy connection string
```

Save this connection string - you'll need it next.

## 2. Backend (10 minutes)

```bash
# 1. Go to https://railway.app
# 2. New Project → Deploy from GitHub
# 3. Select your repo
# 4. Add these environment variables:
```

| Variable | Value |
|----------|-------|
| `DATABASE_URL` | Your Neon connection string |
| `JWT_SECRET` | Run: `openssl rand -base64 32` |
| `PORT` | `8080` |
| `ENV` | `production` |
| `CORS_ORIGINS` | `https://localhost` (update later) |
| `RATE_LIMIT_ENABLED` | `true` |

```bash
# 5. Wait for deploy
# 6. Copy your Railway URL: https://your-app.up.railway.app
# 7. Test: curl https://your-app.up.railway.app/health
```

## 3. Frontend (10 minutes)

```bash
# 1. Update frontend/.env.production with your Railway URL
# 2. Go to https://netlify.com
# 3. New site → Import from GitHub
# 4. Configure:
```

| Setting | Value |
|---------|-------|
| Base directory | `frontend` |
| Build command | `npm run build` |
| Publish directory | `frontend/dist` |

```bash
# 5. Add environment variable:
```

| Variable | Value |
|----------|-------|
| `VITE_API_URL` | Your Railway URL |

```bash
# 6. Deploy
# 7. Copy your Netlify URL: https://your-app.netlify.app
```

## 4. Update CORS (2 minutes)

```bash
# 1. Go back to Railway
# 2. Update CORS_ORIGINS variable with your Netlify URL
# 3. Railway will auto-redeploy
```

## 5. Test (5 minutes)

Visit your Netlify URL and test:
- ✅ Homepage loads with wallpaper
- ✅ Register account
- ✅ Login
- ✅ Start battle
- ✅ Visit shop

## Done! 🎉

Your app is live at: `https://your-app.netlify.app`

---

## Troubleshooting

**CORS errors?**
- Check `CORS_ORIGINS` in Railway matches your Netlify URL exactly

**Can't connect to database?**
- Verify `DATABASE_URL` in Railway
- Check Neon dashboard - database might be paused

**Frontend shows "Network Error"?**
- Check `VITE_API_URL` in Netlify
- Verify Railway backend is running

**Assets not loading?**
- Ensure `frontend/public/assets/` is in your repo
- Check Netlify deploy logs

---

For detailed troubleshooting, see [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md)
