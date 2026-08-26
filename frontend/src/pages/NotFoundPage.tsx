import { Link } from 'react-router-dom'

// library-docs/12-ux-ui/navigation-map.md: "Undefined routes show a 404 screen
// linking back to /dashboard."
export function NotFoundPage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-3 bg-slate-50 text-center">
      <p className="text-6xl font-bold text-primary-600">404</p>
      <p className="text-slate-600">This page doesn't exist.</p>
      <Link to="/dashboard" className="text-sm font-medium text-primary-600 hover:underline">
        Back to dashboard
      </Link>
    </div>
  )
}
