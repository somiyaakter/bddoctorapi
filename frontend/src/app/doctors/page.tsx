import {
	ApiRequestError,
	listDoctors,
	listLocations,
	listSpecialties,
} from '@/lib/api'

import { DoctorGrid } from '../../components/ui/doctor/doctor-grid'
import { FilterPanel } from '../../components/ui/search/filter-panel'
import { SearchBar } from '../../components/ui/search/search-bar'
import { EmptyState } from '../../components/ui/shared/empty-state'
import { ErrorState } from '../../components/ui/shared/error-state'
import { Pagination } from '../../components/ui/shared/pagination'

interface DoctorsPageProps {
	searchParams: Promise<{
		q?: string
		location_id?: string
		specialty_id?: string
		page?: string
	}>
}


export default async function DoctorsPage({ searchParams }: DoctorsPageProps) {
	const params = await searchParams

	const page = params.page ? parseInt(params.page, 10) : 1

	const locationId = params.location_id
		? parseInt(params.location_id, 10)
		: undefined

	const specialtyId = params.specialty_id
		? parseInt(params.specialty_id, 10)
		: undefined

	let doctorsRes
	let locationsRes
	let specialtiesRes

	try {
		;[doctorsRes, locationsRes, specialtiesRes] = await Promise.all([
			listDoctors({
				page,
				pageSize: 12,
				q: params.q,
				locationId,
				specialtyId,
			}),

			listLocations(),

			listSpecialties(locationId),
		])
	} catch (err) {
		const message =
			err instanceof ApiRequestError
				? err.message
				: 'Failed to load doctors. Please try again.'

		return (
			<div className="mx-auto max-w-6xl px-4 py-16 sm:px-6">
				<ErrorState message={message} />
			</div>
		)
	}

	const doctors = doctorsRes.data ?? []
	const pagination = doctorsRes.pagination

	return (
		<div className="mx-auto max-w-6xl px-4 py-8 sm:px-6">
			{/* Header */}
			<div className="mb-6">
				<h1 className="font-display text-2xl font-semibold text-ink sm:text-3xl">
					Find a Doctor
				</h1>

				<p className="mt-1 text-sm text-slate">
					{pagination?.total_items?.toLocaleString() ?? 0} doctors across
					Bangladesh
				</p>
			</div>

			{/* Search */}
			<div className="mb-6">
				<SearchBar defaultValue={params.q} />
			</div>

			{/* Content */}
			<div className="grid grid-cols-1 gap-6 lg:grid-cols-[240px_1fr]">
				{/* Filters */}
			<aside className="lg:sticky lg:top-20 lg:self-start">
					<FilterPanel
						locations={locationsRes.data ?? []}
						specialties={specialtiesRes.data ?? []}
						selectedLocationId={locationId}
						selectedSpecialtyId={specialtyId}
					/>
				</aside>

				{/* Doctors */}
				<div>
					{doctors.length === 0 ? (
						<EmptyState
							title="No doctors found"
							description="Try adjusting your search or filters."
						/>
					) : (
						<>
							<DoctorGrid doctors={doctors} />

							{pagination && pagination.total_pages > 1 && (
								<div className="mt-8">
									<Pagination
										currentPage={pagination.page}
										totalPages={pagination.total_pages}
									/>
								</div>
							)}
						</>
					)}
				</div>
			</div>
		</div>
	)
}
