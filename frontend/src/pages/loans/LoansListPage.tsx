import { Link } from 'react-router-dom'

import { Button } from '../../components/ui/Button'
import { EmptyState } from '../../components/ui/EmptyState'

export function LoansListPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-slate-900">Loans</h1>
          <p className="text-sm text-slate-500">Active loans and full history.</p>
        </div>
        <Link to="/loans/new">
          <Button>Register loan</Button>
        </Link>
      </div>
      <EmptyState
        title="Loan history coming soon"
        description="Viewing active loans and full history lands with HU-07. Registration (HU-06) is already available above."
        hu="HU-07"
      />
    </div>
  )
}
