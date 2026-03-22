document.addEventListener("DOMContentLoaded", function() {
	const searchInput = document.getElementById("video-search");
	if (!searchInput) return;

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

document.addEventListener('click', function(e) {
	const btn = e.target.closest('button[data-url]');
	if (!btn) return;
	downloadVideo(btn.dataset.url, btn.dataset.filename);
});

var currentVideoID = '';

function openResolutionModal(button) {
	currentVideoID = button.dataset.videoId;
	document.getElementById('resolution-modal-title').textContent = button.dataset.videoTitle;
	new bootstrap.Modal(document.getElementById('resolution-modal')).show();
}

document.addEventListener('click', function(e) {
	const btn = e.target.closest('.resolution-btn');
	if (!btn || !currentVideoID) return;
	const formData = new FormData();
	formData.append('format', btn.dataset.height);
	fetch('/videos/' + currentVideoID + '/download', {
		method: 'POST',
		body: formData,
	}).then(function(resp) {
		if (resp.ok) {
			bootstrap.Modal.getInstance(document.getElementById('resolution-modal')).hide();
			location.reload();
		}
	});
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
