# URL Shortener Frontend

Frontend application for URL shortener service built with Next.js, React, and shadcn/ui.

## Tech Stack

- Next.js 15.5.2
- React 19.1.0
- TypeScript
- React Query (TanStack Query)
- React Hook Form
- Zod (validation)
- shadcn/ui components
- Tailwind CSS
- next-intl (internationalization)

## Getting Started

Install dependencies:

```bash
pnpm install
```

Run the development server:

```bash
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000) to see the application.

## Environment Variables

Create a `.env.local` file with:

```
API_URL=http://localhost:8080
```

## Project Structure

```
src/
├── app/              # Next.js app directory
├── components/       # React components
│   ├── client/      # Client components
│   └── ui/          # shadcn/ui components
├── api/             # API client and endpoints
├── hooks/           # Custom React hooks
├── i18n/            # Internationalization
└── utils/           # Utility functions
```
