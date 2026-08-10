async function loadFiles() {
	const container = document.getElementById("files");

	try {
		const response = await fetch("/3usSroyU/files");

		if (!response.ok) {
			throw new Error(`HTTP ${response.status}`);
		}

		const files = await response.json();

		container.replaceChildren();

		if (files.length === 0) {
			container.textContent = "No files.";
			return;
		}

		for (const file of files) {
			const item = document.createElement("div");

			const name = document.createElement("a");
			name.href = `/download/${file.id}`;
			name.textContent = file.name;

			const info = document.createElement("span");
			info.textContent =
				` (${file.size} bytes, ${file.mime_type}, uploaded ${file.upload_time})`;

			item.append(name, info);
			container.appendChild(item);
		}
	} catch (error) {
		console.error(error);
		container.textContent = "Failed to load files.";
	}

}

async function uploadFile(file) {
	const status = document.getElementById("upload-status");

	status.textContent = "Creating upload...";

	const createResponse = await fetch("/upload", {
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({
			name: file.name,
			mime_type: file.type
		})
	});

	if (!createResponse.ok) {
		throw new Error(`Create upload failed: HTTP ${createResponse.status}`);
	}

	const upload = await createResponse.json();

	const id = upload.id;
	let offset = upload.offset;

	while (offset < file.size) {
		const data = file.slice(offset);

		status.textContent =
			`Uploading ${offset} / ${file.size} bytes...`;

		const response = await fetch(`/upload/${id}`, {
			method: "PATCH",
			headers: {
				"Upload-Offset": offset.toString()
			},
			body: data
		});

		if (!response.ok) {
			throw new Error(`Upload failed: HTTP ${response.status}`);
		}

		const result = await response.json();

		offset = result.offset;
	}

	status.textContent = "Committing...";

	const commitResponse = await fetch(`/upload/${id}/commit`, {
		method: "POST"
	});

	if (!commitResponse.ok) {
		throw new Error(`Commit failed: HTTP ${commitResponse.status}`);
	}

	status.textContent = "Upload complete.";

	await loadFiles();

}

document.getElementById("upload-button").addEventListener("click", async () => {
	const input = document.getElementById("upload-file");
	const file = input.files[0];

	if (!file) {
		document.getElementById("upload-status").textContent =
			"Please select a file.";
		return;
	}

	const button = document.getElementById("upload-button");

	button.disabled = true;

	try {
		await uploadFile(file);
	} catch (error) {
		console.error(error);

		document.getElementById("upload-status").textContent =
			`Upload failed: ${error.message}`;
	} finally {
		button.disabled = false;
	}

});

loadFiles();

