# Get Started

Complete beginner's guide to running PokeTacTix.

## 🎯 What You Need

**Docker Desktop** - That's it!

Download from: https://www.docker.com/products/docker-desktop

## 📦 Get the Code

**Option 1: Download ZIP**
1. Go to https://github.com/IfrunRuhin12/PokeTacTix
2. Click "Code" → "Download ZIP"
3. Extract the ZIP
4. Open terminal in that folder

**Option 2: Use Git**
```bash
git clone https://github.com/IfrunRuhin12/PokeTacTix.git
cd PokeTacTix
```

## 🚀 Start Everything

```bash
make dev
```

Or:
```bash
./scripts/docker-dev.sh
```

**First time?** It takes 2-5 minutes to download images.

## ⏳ Wait for It

You'll see:
```
✅ Database is ready!
📝 Running database migrations...
✅ Backend is ready!
✅ Frontend is ready!

🎉 All Services Running!

Frontend:     http://localhost:5173
Backend API:  http://localhost:3000
API Docs:     http://localhost:3000/api/docs
```

## 🎮 Use the App

1. Open http://localhost:5173
2. Register an account
3. Start playing!

## 🛑 Stop Services

```bash
make stop
```

## � Start Again

```bash
make dev
```

Your data is saved!

## ❓ Troubleshooting

**"Docker is not running"**
- Open Docker Desktop
- Wait for it to start
- Try again

**"Port already in use"**
```bash
make clean
make dev
```

**Still stuck?**
```bash
make logs    # See what's happening
make status  # Check services
```

## 💡 Tips

- Docker Desktop must be running
- First start is slow (downloads images)
- Use `make help` to see all commands
- Check `make logs` if something breaks

## 🎓 What's Next?

- Register an account
- Get starter Pokemon
- Try a 1v1 battle
- Explore the shop
- Build your deck

---

**Need more details?** See [development.md](development.md)
