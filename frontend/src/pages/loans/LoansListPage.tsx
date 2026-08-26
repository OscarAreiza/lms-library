import { EmptyState } from '../../components/ui/EmptyState'

export function LoansListPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">Loans</h1>
        <p className="text-sm text-slate-500">Active loans and full history.</p>
      </div>
      <EmptyState
        title="Loan registration and history coming soon"
        description="Registering a loan and its return lands with HU-06 and HU-07."
        hu="HU-06 / HU-07"
      />
    </div>
  )
}
