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
