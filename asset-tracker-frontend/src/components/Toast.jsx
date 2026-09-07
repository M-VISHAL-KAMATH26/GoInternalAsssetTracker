function Toast({ message, type = 'success', onClose }) {
	if (!message) {
		return null
	}

	const styles = type === 'error'
		? 'border-red-200 bg-red-50 text-red-700'
		: 'border-emerald-200 bg-emerald-50 text-emerald-700'

	return (
		<div className={`fixed right-4 top-4 z-50 flex max-w-sm items-start gap-4 rounded-lg border px-4 py-3 text-sm shadow-lg ${styles}`} role={type === 'error' ? 'alert' : 'status'}>
			<p>{message}</p>
			<button className="font-semibold opacity-70 hover:opacity-100" onClick={onClose} type="button">
				<span className="sr-only">Dismiss notification</span>
				<span aria-hidden="true">&times;</span>
			</button>
		</div>
	)
}

export default Toast