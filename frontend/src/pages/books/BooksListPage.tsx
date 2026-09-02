import { Link } from 'react-router-dom'

import { Button } from '../../components/ui/Button'
import { EmptyState } from '../../components/ui/EmptyState'

export function BooksListPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-slate-900">Catalog</h1>
          <p className="text-sm text-slate-500">Search the book catalog and check availability.</p>
        </div>
        <Link to="/books/new">
          <Button>Register book</Button>
        </Link>
      </div>
      <EmptyState
        title="Catalog search and editing coming soon"
        description="Search (HU-05) and editing (HU-09) land separately. Registration (HU-04) is already available above."
        hu="HU-05 / HU-09"
      />
    </div>
  )
}
