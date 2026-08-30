import { EmptyState } from '../../components/ui/EmptyState'

export function OverdueLoansPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">Overdue loans</h1>
        <p className="text-sm text-slate-500">Loans past their due date and still active.</p>
      </div>
      <EmptyState
        title="Overdue tracking coming soon"
        description="Overdue detection and the suspension policy land with HU-08."
        hu="HU-08"
      />
    </div>
  )
}
