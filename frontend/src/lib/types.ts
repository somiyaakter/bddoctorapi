export interface Chamber {
	id: number
	doctor_id: number
	name: string
	address: string
	visiting_hour: string
	appointment_phone: string
}

export interface Doctor {
	id: number
	name: string
	bmdc_reg_no: string
	degrees: string
	experience_years: number
	specialties: string
	designation: string
	workplace: string
	image_url: string
	profile_url: string
	chambers: Chamber[]
	created_at: string
	updated_at: string
}

export interface Location {
	id: number
	name: string
	url: string
}

export interface Specialty {
	id: number
	location_id: number
	name: string
	url: string
}

export interface Pagination {
	page: number
	page_size: number
	total_items: number
	total_pages: number
}

export interface ApiListResponse<T> {
	data: T[]
	pagination?: Pagination
}

export interface ApiSingleResponse<T> {
	data: T
}
