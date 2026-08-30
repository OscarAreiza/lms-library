import type { PropsWithChildren } from 'react'
import { Navigate } from 'react-router-dom'

import { isAuthenticated } from '../../lib/auth'
import { AppLayout } from './AppLayout'

// Every route under /dashboard, /students, /books, /loans redirects to /login
// without a session — library-docs/12-ux-ui/navigation-map.md, "Navigation rules".
export function ProtectedRoute({ children }: PropsWithChildren) {
  if (!isAuthenticated()) {
    return <Navigate to="/login" replace />
  }
  return <AppLayout>{children}</AppLayout>
}
