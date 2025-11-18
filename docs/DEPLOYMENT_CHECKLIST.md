# PokeTacTix Deployment Checklist

Use this checklist to ensure a smooth deployment process.

## Pre-Deployment

- [ ] Code is pushed to GitHub
- [ ] All tests pass locally
- [ ] Frontend builds successfully (`cd frontend && npm run build`)
- [ ] Backend builds successfully (`go build ./cmd/api`)
- [ ] Database migrations are in `internal/database/migrations/`
- [ ] Assets are in `frontend/public/assets/` (wallpaper.jpg, pokeball.png)

## Database (Neon)

- [ ] Created Neon account
- [ ] Created new project
- [ ] Copied connection string
- [ ] Connection string includes `?sslmode=require`
- [ ] Noted free tier limits (0.5 GB storage, 191.9 hours/month)

## Backend (Railway)

- [ ] Created Railway account
- [ ] Created new project from GitHub repo
- [ ] Set environment variables:
  - [ ] `DATABASE_URL` (from Neon)
  - [ ] `JWT_SECRET` (generated with `openssl rand -base64 32`)
  - [ ] `PORT=8080`
  - [ ] `ENV=production`
  - [ ] `CORS_ORIGINS` (will update after Netlify)
  - [ ] `RATE_LIMIT_ENABLED=true`
- [ ] Deployment successful
- [ ] Health endpoint works: `https://your-app.up.railway.app/health`
- [ ] Checked logs for errors
- [ ] Database migrations ran successfully

## Frontend (Netlify)

- [ ] Created Netlify account
- [ ] Created `netlify.toml` in project root
- [ ] Created `frontend/_redirects` file
- [ ] Created `frontend/.env.production` with Railway URL
- [ ] Created new site from GitHub repo
- [ ] Set build settings:
  - [ ] Base directory: `frontend`
  - [ ] Build command: `npm run build`
  - [ ] Publish directory: `frontend/dist`
- [ ] Set environment variable:
  - [ ] `VITE_API_URL` (Railway URL)
- [ ] Deployment successful
- [ ] Site loads correctly
- [ ] Assets load (wallpaper, pokeball)

## CORS Configuration

- [ ] Copied Netlify URL
- [ ] Updated `CORS_ORIGINS` in Railway with Netlify URL
- [ ] Railway redeployed automatically
- [ ] Tested frontend can connect to backend

## Testing

- [ ] Homepage loads with wallpaper and pokeballs
- [ ] Can register new account
- [ ] Can login
- [ ] Dashboard loads
- [ ] Can start a battle
- [ ] Shop loads and shows Pokemon
- [ ] Can purchase Pokemon (if have coins)
- [ ] Deck manager works
- [ ] Profile page shows stats
- [ ] Can logout

## Security

- [ ] JWT_SECRET is strong and unique
- [ ] No sensitive data in GitHub repo
- [ ] `.env` file is in `.gitignore`
- [ ] CORS is properly configured
- [ ] Rate limiting is enabled
- [ ] HTTPS is enabled (automatic)

## Monitoring

- [ ] Bookmarked Neon dashboard
- [ ] Bookmarked Railway dashboard
- [ ] Bookmarked Netlify dashboard
- [ ] Checked Railway logs for errors
- [ ] Checked Netlify deploy logs

## Documentation

- [ ] Updated README with production URLs
- [ ] Documented any custom configuration
- [ ] Saved all credentials securely (password manager)

## Optional Enhancements

- [ ] Set up custom domain on Netlify
- [ ] Set up custom domain on Railway
- [ ] Add error tracking (Sentry, etc.)
- [ ] Add analytics (Google Analytics, etc.)
- [ ] Set up monitoring alerts
- [ ] Configure automatic backups

## Post-Deployment

- [ ] Shared app URL with team/users
- [ ] Monitored for first 24 hours
- [ ] Checked for any errors in logs
- [ ] Verified database usage in Neon
- [ ] Verified Railway credit usage

---

## Quick URLs Reference

After deployment, fill these in:

- **Frontend**: https://_____.netlify.app
- **Backend**: https://_____.up.railway.app
- **Database**: Neon dashboard URL
- **GitHub**: Your repository URL

---

## Emergency Contacts

- **Neon Support**: https://neon.tech/docs/introduction/support
- **Railway Support**: https://railway.app/help
- **Netlify Support**: https://www.netlify.com/support/

---

## Rollback Plan

If something goes wrong:

1. **Frontend**: Netlify allows instant rollback to previous deploy
2. **Backend**: Railway allows rollback to previous deployment
3. **Database**: Neon has automatic backups (7-day retention on free tier)

---

**Last Updated**: [Add date when you deploy]
