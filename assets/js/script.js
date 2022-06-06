async function deleteDevice (id) {
	try {
		await fetch(`/api/devices/${id}`, { method: 'DELETE' })
		document.getElementById(id).remove()
	} catch(err) {
		alert("action unsuccessful. please try again.")
	}
}
