import Image from 'next/image'
import type { Doctor } from '@/lib/types'

export function DoctorDetails({ doctor }: { doctor: Doctor }) {
	const jsonLd = {
		'@context': 'https://schema.org',
		'@type': 'Physician',
		name: doctor.name,
		medicalSpecialty: doctor.specialties || undefined,
		url: `https://yourdomain.com/doctors/${doctor.id}`,
		image: doctor.image_url || undefined,
		identifier: doctor.bmdc_reg_no || undefined,
	}

	return (
		<main className="min-h-screen bg-paper">
			{/* SEO structured data */}
			<script
				type="application/ld+json"
				dangerouslySetInnerHTML={{
					__html: JSON.stringify(jsonLd),
				}}
			/>

			<div className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8 lg:py-12">

				{/* Breadcrumb */}
				<nav className="mb-6 flex items-center gap-2 text-sm">
					<span className="text-slate">Doctors</span>

					<span className="text-hairline">/</span>

					<span className="truncate font-medium text-ink">
						{doctor.name}
					</span>
				</nav>

				{/* ================================================== */}
				{/* DOCTOR HEADER */}
				{/* ================================================== */}

				<section className="overflow-hidden rounded-2xl border border-hairline bg-white shadow-sm">

					{/* Brand accent */}
					<div className="h-1.5 bg-brand" />

					<div className="p-6 sm:p-8">
						<div className="flex flex-col gap-6 sm:flex-row sm:items-center">

							{/* Doctor Image */}
							<div className="relative h-28 w-28 shrink-0 overflow-hidden rounded-2xl border border-hairline bg-secondary sm:h-32 sm:w-32">

								{doctor.image_url ? (
									<Image
										src={doctor.image_url}
										alt={doctor.name}
										fill
										sizes="128px"
										className="object-cover"
									/>
								) : (
									<div className="flex h-full w-full items-center justify-center bg-secondary text-4xl font-semibold text-brand">
										{doctor.name.charAt(0).toUpperCase()}
									</div>
								)}
							</div>

							{/* Doctor Identity */}
							<div className="min-w-0 flex-1">

								<div className="flex flex-wrap items-center gap-2">

									<h1 className="font-display text-2xl font-semibold tracking-tight text-ink sm:text-3xl">
										{doctor.name}
									</h1>

									{doctor.bmdc_reg_no && (
										<span className="inline-flex items-center gap-1.5 rounded-full border border-brand/20 bg-brand/5 px-2.5 py-1 text-xs font-medium text-brand">
											<span className="flex h-4 w-4 items-center justify-center rounded-full bg-brand text-[9px] font-bold text-white">
												✓
											</span>

											Verified
										</span>
									)}
								</div>

								{/* Degrees */}
								{doctor.degrees && (
									<p className="mt-2 text-sm leading-6 text-slate sm:text-base">
										{doctor.degrees}
									</p>
								)}

								{/* Designation */}
								{doctor.designation && (
									<p className="mt-1 text-sm font-medium text-brand">
										{doctor.designation}
									</p>
								)}

								{/* BMDC */}
								{doctor.bmdc_reg_no && (
									<div className="mt-4 inline-flex flex-wrap items-center gap-x-3 gap-y-1 rounded-lg bg-paper px-3 py-2">
										<span className="text-[11px] font-semibold uppercase tracking-[0.12em] text-slate">
											BMDC Registration
										</span>

										<span className="font-mono text-xs font-medium text-ink">
											{doctor.bmdc_reg_no}
										</span>
									</div>
								)}
							</div>
						</div>
					</div>
				</section>

				{/* ================================================== */}
				{/* MAIN CONTENT */}
				{/* ================================================== */}

				<div className="mt-8 grid gap-8 lg:grid-cols-[minmax(0,1fr)_320px]">

					{/* ================================================== */}
					{/* LEFT COLUMN */}
					{/* ================================================== */}

					<div className="space-y-8">

						{/* ------------------------------------------------ */}
						{/* PROFESSIONAL INFORMATION */}
						{/* ------------------------------------------------ */}

						<section className="overflow-hidden rounded-2xl border border-hairline bg-white">

							<div className="border-b border-hairline px-6 py-5">
								<h2 className="font-display text-xl font-semibold text-ink">
									Professional Information
								</h2>

								<p className="mt-1 text-sm text-slate">
									Qualifications, specialty and professional experience
								</p>
							</div>

							<dl className="grid sm:grid-cols-2">

								{doctor.specialties && (
									<InfoItem
										label="Specialty"
										value={doctor.specialties}
									/>
								)}

								{doctor.designation && (
									<InfoItem
										label="Designation"
										value={doctor.designation}
									/>
								)}

								{doctor.workplace && (
									<InfoItem
										label="Workplace"
										value={doctor.workplace}
									/>
								)}

								{doctor.experience_years > 0 && (
									<InfoItem
										label="Experience"
										value={`${doctor.experience_years}+ years`}
									/>
								)}

							</dl>
						</section>

						{/* ------------------------------------------------ */}
						{/* CHAMBERS */}
						{/* ------------------------------------------------ */}

						{doctor.chambers.length > 0 && (
							<section>

								<div className="mb-5">
									<h2 className="font-display text-xl font-semibold text-ink">
										Chambers &amp; Appointments
									</h2>

									<p className="mt-1 text-sm text-slate">
										Where to visit and how to contact this doctor
									</p>
								</div>

								<div className="space-y-4">

									{doctor.chambers.map((chamber, index) => (
										<article
											key={chamber.id}
											className="overflow-hidden rounded-2xl border border-hairline bg-white transition-shadow hover:shadow-sm"
										>
											<div className="p-5 sm:p-6">

												<div className="flex items-start gap-4">

													{/* Chamber number */}
													<div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-secondary text-xs font-bold text-brand">
														{String(index + 1).padStart(2, '0')}
													</div>

													<div className="min-w-0 flex-1">

														{/* Chamber Name */}
														{chamber.name && (
															<h3 className="font-display text-base font-semibold text-ink">
																{chamber.name}
															</h3>
														)}

														{/* Address */}
														{chamber.address && (
															<div className="mt-3 flex items-start gap-2.5">

																<span
																	className="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-secondary text-xs text-brand"
																	aria-hidden="true"
																>
																	⌖
																</span>

																<p className="text-sm leading-6 text-slate">
																	{chamber.address}
																</p>

															</div>
														)}

														{/* Visiting Hours */}
														{chamber.visiting_hour && (
															<div className="mt-3 flex items-start gap-2.5">

																<span
																	className="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-secondary text-xs text-brand"
																	aria-hidden="true"
																>
																	◷
																</span>

																<div>
																	<p className="text-[11px] font-semibold uppercase tracking-[0.1em] text-slate">
																		Visiting Hours
																	</p>

																	<p className="mt-1 text-sm leading-5 text-ink">
																		{chamber.visiting_hour}
																	</p>
																</div>

															</div>
														)}

														{/* Appointment */}
														{chamber.appointment_phone && (
															<div className="mt-5 flex flex-col gap-3 border-t border-hairline pt-4 sm:flex-row sm:items-center">

																<a
																	href={`tel:${chamber.appointment_phone}`}
																	className="inline-flex items-center justify-center gap-2 rounded-lg bg-brand px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-brand-hover"
																>
																	<span aria-hidden="true">
																		☎
																	</span>

																	Call for Appointment
																</a>

																<span className="text-sm font-medium text-ink">
																	{chamber.appointment_phone}
																</span>

															</div>
														)}

													</div>
												</div>
											</div>
										</article>
									))}

								</div>
							</section>
						)}

					</div>

					{/* ================================================== */}
					{/* RIGHT SIDEBAR */}
					{/* ================================================== */}

					<aside className="space-y-4">

						{/* Quick Summary */}
						<section className="rounded-2xl border border-ink bg-ink p-6 text-white">

							<p className="text-[11px] font-semibold uppercase tracking-[0.14em] text-seal">
								Doctor Profile
							</p>

							<h2 className="mt-2 font-display text-xl font-semibold leading-snug">
								{doctor.name}
							</h2>

							{doctor.specialties && (
								<p className="mt-2 text-sm leading-6 text-white/65">
									{doctor.specialties}
								</p>
							)}

							<div className="mt-6 space-y-4 border-t border-white/10 pt-5">

								{doctor.experience_years > 0 && (
									<SummaryItem
										label="Experience"
										value={`${doctor.experience_years}+ years`}
									/>
								)}

								<SummaryItem
									label="Chambers"
									value={String(doctor.chambers.length)}
								/>

								{doctor.bmdc_reg_no && (
									<SummaryItem
										label="BMDC"
										value={doctor.bmdc_reg_no}
										mono
									/>
								)}

							</div>
						</section>

						{/* Directory Note */}
						<section className="rounded-2xl border border-hairline bg-secondary p-5">

							<div className="flex items-start gap-3">

								<div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-brand text-xs font-bold text-white">
									✓
								</div>

								<div>
									<h3 className="text-sm font-semibold text-ink">
										Directory Information
									</h3>

									<p className="mt-1.5 text-xs leading-5 text-slate">
										This profile is provided for directory and
										reference purposes. Please confirm appointment
										details directly with the chamber.
									</p>
								</div>

							</div>
						</section>

					</aside>
				</div>
			</div>
		</main>
	)
}

/* ================================================== */
/* INFO ITEM */
/* ================================================== */

function InfoItem({
	label,
	value,
}: {
	label: string
	value: string
}) {
	return (
		<div className="border-b border-hairline px-6 py-5 last:border-b-0 sm:nth-[odd]:border-r sm:nth-[3]:border-b-0">
			<dt className="text-[11px] font-semibold uppercase tracking-[0.12em] text-slate">
				{label}
			</dt>

			<dd className="mt-2 text-sm font-medium leading-6 text-ink">
				{value}
			</dd>
		</div>
	)
}

/* ================================================== */
/* SUMMARY ITEM */
/* ================================================== */

function SummaryItem({
	label,
	value,
	mono = false,
}: {
	label: string
	value: string
	mono?: boolean
}) {
	return (
		<div className="flex items-center justify-between gap-4">
			<span className="text-sm text-white/55">
				{label}
			</span>

			<span
				className={
					mono
						? 'font-mono text-xs font-medium text-white'
						: 'text-sm font-medium text-white'
				}
			>
				{value}
			</span>
		</div>
	)
}