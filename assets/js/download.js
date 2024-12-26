function showLoadingIndicator() {
    const loadingIndicator = document.createElement('div');
    loadingIndicator.id = 'loading-indicator';
    loadingIndicator.className = 'fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 cursor-not-allowed select-none';
    loadingIndicator.innerHTML = `
      <div class="bg-white p-10 rounded shadow-md flex items-center text-xl space-x-4">
        <span class="text-gray-700">Download genereren...</span>
      </div>
    `;
    document.body.appendChild(loadingIndicator);

    document.body.style.pointerEvents = 'none';
    document.body.style.userSelect = 'none';
}

function hideLoadingIndicator() {
    const loadingIndicator = document.getElementById('loading-indicator');
    if (loadingIndicator) {
        loadingIndicator.remove();
    }
    document.body.style.pointerEvents = '';
    document.body.style.userSelect = '';
}

function handleDownloadClick(event) {
    const downloadUrl = event.target.getAttribute('data-download-url');
    if (!downloadUrl) {
        return;
    }

    showLoadingIndicator()

    fetch(downloadUrl).then((response) => {
        if (!response.ok) {
            return;
        }

        response.blob().then((blob) => {
            const blobUrl = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = blobUrl;
            a.download = event.target.getAttribute('data-download-name');
            a.style.display = 'none';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(blobUrl);
            hideLoadingIndicator()
        }).catch((err) => console.error(err));

    }).catch((err) => console.error(err));
}

document.addEventListener('click', function(event) {
    if (event.target.hasAttribute('data-download-url')) {
        handleDownloadClick(event);
    }
});
