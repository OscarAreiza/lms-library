import { EmptyState } from '../../components/ui/EmptyState'

export function BooksListPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">Catalog</h1>
        <p className="text-sm text-slate-500">Search the book catalog and check availability.</p>
      </div>
      <EmptyState
        title="Book catalog coming soon"
        description="Registration, search, and editing land with HU-04, HU-05, and HU-09."
        hu="HU-04 / HU-05 / HU-09"
      />
    </div>
  )
}
