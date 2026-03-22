document.addEventListener("DOMContentLoaded", function() {
	const searchInput = document.getElementById("video-search");
	if (!searchInput) return;

	document.querySelectorAll('button[data-url]').forEach(function(btn) {
		btn.addEventListener('click', function() {
			const url = this.dataset.url;
			const filename = this.dataset.filename;
			downloadVideo(url, filename);
		});
	});

	function performSearch() {
		const videoCards = Array.from(document.querySelectorAll(".video-card"));
		const query = searchInput.value.trim();

		if (!query) {
			videoCards.forEach((card) => (card.style.display = ""));
			return;
		}

		const videos = videoCards.map((card) => ({
			element: card,
			title: card.dataset.title || "",
			channelName: card.dataset.channelName || "",
		}));

		const fuse = new Fuse(videos, {
			keys: ["title", "channelName"],
			threshold: 0.4,
			includeScore: true,
		});

		const results = fuse.search(query);
		const resultElements = new Set(results.map((r) => r.item.element));

		videoCards.forEach((card) => {
			card.style.display = resultElements.has(card) ? "" : "none";
		});
	}

	searchInput.addEventListener("input", performSearch);

	if (searchInput.value.trim()) {
		performSearch();
	}
});

function addVideoHandler(event) {
	if (event.detail.successful) {
		bootstrap.Modal.getInstance(
			document.getElementById("add-video-modal"),
		).hide();
	} else {
		const field = event.detail.elt.querySelector(`[name="video_id"]`);
		field.setCustomValidity(event.detail.xhr.responseText);
		field.onfocus = () => field.reportValidity();
		field.onchange = () => field.setCustomValidity("");
		field.reportValidity();
	}
}

function subscribeHandler(event) {
	if (event.detail.successful) {
		bootstrap.Modal.getInstance(
			document.getElementById("subscription-modal"),
		).hide();
	} else {
		const field = event.detail.elt.querySelector(`[name="channel_id"]`);
		field.setCustomValidity(event.detail.xhr.responseText);
		field.onfocus = field.reportValidity;
		field.onchange = () => field.setCustomValidity("");
		field.reportValidity();
	}
}

if ("serviceWorker" in navigator) {
	navigator.serviceWorker.register("/assets/sw.js").then((registration) => {
		console.log("Service Worker registered with scope:", registration.scope);
	}).catch((error) => {
		console.log("Service Worker registration failed:", error);
	});
}

function downloadVideo(url, filename) {
	const xhr = new XMLHttpRequest();
	xhr.open("GET", url, true);
	xhr.responseType = "blob";
	xhr.onload = function () {
		if (this.status === 200) {
			const blob = this.response;
			const link = document.createElement("a");
			link.href = window.URL.createObjectURL(blob);
			link.download = filename;
			link.click();
			window.URL.revokeObjectURL(link.href);
		}
	};
	xhr.send();
}

function openResolutionModal(button) {
	const videoID = button.dataset.videoId;
	const title = button.dataset.videoTitle;

	const modalElement = document.getElementById("resolution-modal");
	const modalTitle = modalElement.querySelector(".modal-title");
	const formatOptions = document.getElementById("format-options");

	modalTitle.textContent = `Download: ${title}`;
	formatOptions.innerHTML = '<div class="spinner-border spinner-border-sm text-primary" role="status"><span class="visually-hidden">Loading...</span></div>';

	const modal = new bootstrap.Modal(modalElement);
	modal.show();

	fetch(`/api/videos/${videoID}/formats`)
		.then(response => response.json())
		.then(data => {
			if (data.formats && data.formats.length > 0) {
				const uniqueResolutions = new Map();
				data.formats.forEach(format => {
					if (!uniqueResolutions.has(format.height)) {
						uniqueResolutions.set(format.height, format);
					}
				});

				const sortedResolutions = Array.from(uniqueResolutions.values())
					.sort((a, b) => b.height - a.height);

				formatOptions.innerHTML = sortedResolutions.map(format => {
					const resolutionLabel = format.height > 0 ? `${format.height}p` : 'Unknown';
					const note = format.note ? ` - ${format.note}` : '';
					const formatId = format.format_id ? format.format_id : '';
					return `<button type="button" class="btn btn-outline-primary w-100 text-start" data-format="${formatId}">
						<div class="d-flex justify-content-between align-items-center">
							<span>${resolutionLabel}${note}</span>
							<i class="bi bi-download"></i>
						</div>
					</button>`;
				}).join('');

				formatOptions.querySelectorAll('button[data-format]').forEach(btn => {
					btn.addEventListener('click', function() {
						const format = this.dataset.format;
						downloadWithFormat(videoID, format);
						modal.hide();
					});
				});
			} else {
				formatOptions.innerHTML = '<div class="alert alert-warning">No formats available</div>';
			}
		})
		.catch(error => {
			formatOptions.innerHTML = '<div class="alert alert-danger">Failed to load formats</div>';
			console.error('Error fetching formats:', error);
		});
}

function downloadWithFormat(videoID, format) {
	const formData = new FormData();
	formData.append("format", format);

	const xhr = new XMLHttpRequest();
	xhr.open("POST", `/api/videos/${videoID}/download`, true);
	xhr.onload = function() {
		if (this.status === 200 || this.status === 202) {
			location.reload();
		} else {
			alert('Download failed: ' + (this.responseText || 'Unknown error'));
		}
	};
	xhr.onerror = function() {
		alert('Network error occurred');
	};
	xhr.send(formData);
}
