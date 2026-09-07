import { useCallback, useEffect, useState } from 'react'
import Toast from '../components/Toast'
import { createAsset, listAssets, retireAsset } from '../features/assets/assetsApi'

const initialForm = {
	name: '',
	type: '',
	category: '',
	serialNumber: '',
}

const getAssetList = (data) => {
	if (Array.isArray(data)) {
		return data
	}

	return data?.assets ?? data?.items ?? []
}

const getAssetId = (asset) => asset.id ?? asset.assetId

const fetchAssets = async () => {
	const response = await listAssets()

	return getAssetList(response.data)
}

const statusStyles = {
	available: 'bg-emerald-100 text-emerald-700',
	assigned: 'bg-blue-100 text-blue-700',
	retired: 'bg-slate-200 text-slate-600',
	maintenance: 'bg-amber-100 text-amber-700',
}

const getStatusStyle = (status) =>
	statusStyles[status?.toLowerCase()] ?? 'bg-slate-100 text-slate-600'

function AssetsPage() {
	const [assets, setAssets] = useState([])
	const [isLoading, setIsLoading] = useState(true)
	const [isModalOpen, setIsModalOpen] = useState(false)
	const [isSubmitting, setIsSubmitting] = useState(false)
	const [retiringId, setRetiringId] = useState(null)
	const [form, setForm] = useState(initialForm)
	const [formError, setFormError] = useState('')
	const [pageError, setPageError] = useState('')
	const [feedback, setFeedback] = useState('')

	const loadAssets = useCallback(async (showLoading = true) => {
		if (showLoading) {
			setIsLoading(true)
		}

		setPageError('')

		try {
			setAssets(await fetchAssets())
		} catch {
			setPageError('Unable to load assets. Please try again.')
		} finally {
			setIsLoading(false)
		}
	}, [])

	useEffect(() => {
		let isMounted = true

		const loadInitialAssets = async () => {
			try {
				setAssets(await fetchAssets())
			} catch {
				if (isMounted) {
					setPageError('Unable to load assets. Please try again.')
				}
			} finally {
				if (isMounted) {
					setIsLoading(false)
				}
			}
		}

		loadInitialAssets()

		return () => {
			isMounted = false
		}
	}, [])

	const handleFormChange = (event) => {
		const { name, value } = event.target

		setForm((currentForm) => ({ ...currentForm, [name]: value }))
		setFormError('')
	}

	const closeModal = () => {
		if (!isSubmitting) {
			setIsModalOpen(false)
			setForm(initialForm)
			setFormError('')
		}
	}

	const handleCreate = async (event) => {
		event.preventDefault()

		if (Object.values(form).some((value) => !value.trim())) {
			setFormError('Complete all fields before creating the asset.')
			return
		}

		setIsSubmitting(true)
		setFormError('')
		setFeedback('')

		try {
			await createAsset({
				name: form.name.trim(),
				type: form.type.trim(),
				category: form.category.trim(),
				serialNumber: form.serialNumber.trim(),
			})
			closeModal()
			setFeedback('Asset created successfully.')
			await loadAssets(false)
		} catch {
			setFormError('Unable to create the asset. Please try again.')
		} finally {
			setIsSubmitting(false)
		}
	}

	const handleRetire = async (asset) => {
		const assetId = getAssetId(asset)

		setRetiringId(assetId)
		setFeedback('')

		try {
			await retireAsset(assetId)
			setFeedback('Asset retired successfully.')
			await loadAssets(false)
		} catch {
			setPageError('Unable to retire the asset. Please try again.')
		} finally {
			setRetiringId(null)
		}
	}

	return (
		<main className="min-h-screen bg-slate-50 px-4 py-8 text-slate-900 sm:px-6">
			<Toast message={feedback} onClose={() => setFeedback('')} />
			<div className="mx-auto max-w-7xl">
				<div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
					<div>
						<p className="text-sm font-medium text-slate-500">Admin portal</p>
						<h1 className="mt-1 text-3xl font-semibold tracking-tight">Inventory</h1>
						<p className="mt-2 text-sm text-slate-600">
							Manage the organization&apos;s assets and their current status.
						</p>
					</div>
					<button
						className="inline-flex items-center justify-center rounded-lg bg-slate-900 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-slate-700"
						onClick={() => {
							setFeedback('')
							setIsModalOpen(true)
						}}
						type="button"
					>
						Add asset
					</button>
				</div>

				{pageError && (
					<div className="mb-4 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700" role="alert">
						{pageError}
					</div>
				)}

				<section className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
					{isLoading ? (
						<div className="p-8 text-center text-sm text-slate-500">Loading assets...</div>
					) : assets.length === 0 ? (
						<div className="p-10 text-center">
							<h2 className="text-lg font-medium text-slate-900">No assets yet</h2>
							<p className="mt-2 text-sm text-slate-500">Add an asset to start building the inventory.</p>
						</div>
					) : (
						<div className="overflow-x-auto">
							<table className="min-w-full divide-y divide-slate-200 text-left text-sm">
								<thead className="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
									<tr>
										<th className="px-6 py-3 font-medium" scope="col">Name</th>
										<th className="px-6 py-3 font-medium" scope="col">Type</th>
										<th className="px-6 py-3 font-medium" scope="col">Category</th>
										<th className="px-6 py-3 font-medium" scope="col">Serial number</th>
										<th className="px-6 py-3 font-medium" scope="col">Status</th>
										<th className="px-6 py-3 font-medium" scope="col"><span className="sr-only">Actions</span></th>
									</tr>
								</thead>
								<tbody className="divide-y divide-slate-100">
									{assets.map((asset) => {
										const assetId = getAssetId(asset)
										const status = asset.status ?? '-'

										return (
											<tr className="text-slate-700" key={assetId}>
												<td className="whitespace-nowrap px-6 py-4 font-medium text-slate-900">{asset.name ?? '-'}</td>
												<td className="whitespace-nowrap px-6 py-4">{asset.type ?? asset.assetType ?? '-'}</td>
												<td className="whitespace-nowrap px-6 py-4">{asset.category ?? '-'}</td>
												<td className="whitespace-nowrap px-6 py-4">{asset.serialNumber ?? asset.serial_number ?? '-'}</td>
												<td className="whitespace-nowrap px-6 py-4">
													<span className={`inline-flex rounded-full px-2.5 py-1 text-xs font-medium ${getStatusStyle(status)}`}>
														{status}
													</span>
												</td>
												<td className="whitespace-nowrap px-6 py-4 text-right">
													{status.toLowerCase() === 'available' && (
														<button
															className="text-sm font-medium text-red-700 hover:text-red-900 disabled:cursor-not-allowed disabled:opacity-50"
															disabled={retiringId === assetId}
															onClick={() => handleRetire(asset)}
															type="button"
														>
															{retiringId === assetId ? 'Retiring...' : 'Retire'}
														</button>
													)}
												</td>
											</tr>
										)
									})}
								</tbody>
							</table>
						</div>
					)}
				</section>
			</div>

			{isModalOpen && (
				<div aria-labelledby="add-asset-title" aria-modal="true" className="fixed inset-0 z-10 overflow-y-auto bg-slate-900/40 px-4 py-8" role="dialog">
					<div className="mx-auto max-w-lg rounded-xl bg-white p-6 shadow-xl sm:p-8">
						<div className="mb-6">
							<h2 className="text-xl font-semibold" id="add-asset-title">Add asset</h2>
							<p className="mt-1 text-sm text-slate-500">Add the identifying details for a new inventory item.</p>
						</div>
						<form className="space-y-4" onSubmit={handleCreate}>
							{[
								['name', 'Name', 'MacBook Pro'],
								['type', 'Type', 'Laptop'],
								['category', 'Category', 'Hardware'],
								['serialNumber', 'Serial number', 'SN-12345'],
							].map(([name, label, placeholder]) => (
								<label className="block" htmlFor={`asset-${name}`} key={name}>
									<span className="mb-1.5 block text-sm font-medium text-slate-700">{label}</span>
									<input
										className="w-full rounded-lg border border-slate-300 px-3 py-2.5 text-sm outline-none transition placeholder:text-slate-400 focus:border-slate-500 focus:ring-2 focus:ring-slate-200"
										id={`asset-${name}`}
										name={name}
										onChange={handleFormChange}
										placeholder={placeholder}
										value={form[name]}
									/>
								</label>
							))}
							{formError && <p className="text-sm text-red-600" role="alert">{formError}</p>}
							<div className="flex justify-end gap-3 border-t border-slate-100 pt-5">
								<button className="rounded-lg px-4 py-2.5 text-sm font-medium text-slate-600 hover:bg-slate-100 disabled:opacity-60" disabled={isSubmitting} onClick={closeModal} type="button">Cancel</button>
								<button className="rounded-lg bg-slate-900 px-4 py-2.5 text-sm font-medium text-white hover:bg-slate-700 disabled:cursor-not-allowed disabled:opacity-60" disabled={isSubmitting} type="submit">{isSubmitting ? 'Creating...' : 'Create asset'}</button>
							</div>
						</form>
					</div>
				</div>
			)}
		</main>
	)
}

export default AssetsPage