import { EmptyState } from '../../components/ui/EmptyState'

export function StudentsListPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">Students</h1>
        <p className="text-sm text-slate-500">Search, register, and manage student records.</p>
      </div>
      <EmptyState
        title="Student registry coming soon"
        description="Registration, search, editing, and deactivation land with HU-02 and HU-03."
        hu="HU-02 / HU-03"
      />
    </div>
  )
}
