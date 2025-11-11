# PokeTacTix Frontend

React + Vite frontend for the PokeTacTix web application.

## Tech Stack

- **React 18** - UI library
- **Vite** - Build tool and dev server
- **React Router** - Client-side routing
- **Tailwind CSS** - Utility-first CSS framework
- **Framer Motion** - Animation library
- **Axios** - HTTP client

## Getting Started

### Prerequisites

- Node.js 18+ and npm

### Installation

```bash
# Install dependencies
npm install

# Copy environment variables
cp .env.example .env

# Update .env with your API URL
```

### Development

```bash
# Start dev server (http://localhost:5173)
npm run dev
```

### Build

```bash
# Build for production
npm run build

# Preview production build
npm run preview
```

## Project Structure

```
src/
├── components/       # Reusable UI components
│   ├── auth/        # Authentication components
│   ├── battle/      # Battle-related components
│   ├── shop/        # Shop components
│   ├── deck/        # Deck management components
│   ├── profile/     # Profile components
│   └── common/      # Shared components
├── contexts/        # React contexts
├── hooks/           # Custom React hooks
├── pages/           # Page components
├── services/        # API services
├── utils/           # Utility functions
├── App.jsx          # Main app component
├── main.jsx         # Entry point
└── index.css        # Global styles
```

## Features

- JWT-based authentication
- Protected routes
- Responsive design
- Smooth animations
- Type-safe API client
- Global state management

## Environment Variables

- `VITE_API_URL` - Backend API URL (default: http://localhost:3000)

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint
