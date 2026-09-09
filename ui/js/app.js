const fileInput = document.getElementById('file-input');
const uploadArea = document.getElementById('upload-area');
const previewSection = document.getElementById('preview-section');
const previewFilename = document.getElementById('preview-filename');
const previewImage = document.getElementById('preview-image');
const previewSize = document.getElementById('preview-size');
const previewType = document.getElementById('preview-type');
const removeFileBtn = document.getElementById('remove-file');
const processBtn = document.getElementById('process-btn');
const uploadError = document.getElementById('upload-error');

const jobCard = document.getElementById('job-card');
const jobIdLabel = document.getElementById('job-id');
const statusBadge = document.getElementById('status-badge');
const jobTimeline = document.getElementById('job-timeline');
const pollingInfo = document.getElementById('polling-info');
const tryAgainSection = document.getElementById('try-again-section');
const tryAgainBtn = document.getElementById('try-again-btn');

const resultsSection = document.getElementById('results-section');
const resultsEmpty = document.getElementById('results-empty');
const resultsActive = document.getElementById('results-active');
const resultsGrid = document.getElementById('results-grid');
const resultsError = document.getElementById('results-error');
const resultsErrorMessage = document.getElementById('results-error-message');

let selectedFile = null;
let isSubmitting = false;
let currentJobId = null;
let currentStatusUrl = null;
let pollingController = null;
let pollingInterval = null;

const API_BASE = '';

const ACCEPTED_TYPES = {
    'image/jpeg': true,
    'image/png': true
};

const MAX_SIZE = 10 * 1024 * 1024;

fileInput.addEventListener('change', handleFileSelect);

uploadArea.addEventListener('dragover', (e) => {
    e.preventDefault();
    uploadArea.style.borderColor = '#4361ee';
    uploadArea.style.background = '#f8f9ff';
});

uploadArea.addEventListener('dragleave', () => {
    uploadArea.style.borderColor = '#c0c0c0';
    uploadArea.style.background = '';
});

uploadArea.addEventListener('drop', (e) => {
    e.preventDefault();
    uploadArea.style.borderColor = '#c0c0c0';
    uploadArea.style.background = '';
    const files = e.dataTransfer.files;
    if (files.length > 0) {
        handleFile(files[0]);
    }
});

removeFileBtn.addEventListener('click', () => {
    resetUpload();
});

processBtn.addEventListener('click', () => {
    submitImage();
});

tryAgainBtn.addEventListener('click', () => {
    tryAgain();
});

function handleFileSelect(e) {
    const file = e.target.files[0];
    if (file) {
        handleFile(file);
    }
}

function handleFile(file) {
    hideError();

    if (!ACCEPTED_TYPES[file.type]) {
        showError('Unsupported file type. Please select a JPEG or PNG image.');
        return;
    }

    if (file.size > MAX_SIZE) {
        showError('File is too large. Maximum size is 10 MB.');
        return;
    }

    if (file.size === 0) {
        showError('File is empty.');
        return;
    }

    selectedFile = file;

    previewFilename.textContent = file.name;
    previewImage.src = URL.createObjectURL(file);
    previewSize.textContent = formatFileSize(file.size);
    previewType.textContent = file.type;

    uploadArea.hidden = true;
    previewSection.hidden = false;
    processBtn.disabled = false;
}

function resetUpload() {
    selectedFile = null;
    fileInput.value = '';
    uploadArea.hidden = false;
    previewSection.hidden = true;
    processBtn.disabled = true;
    hideError();
    URL.revokeObjectURL(previewImage.src);
}

function showError(msg) {
    uploadError.textContent = msg;
    uploadError.hidden = false;
}

function hideError() {
    uploadError.hidden = true;
}

async function submitImage() {
    if (isSubmitting || !selectedFile) return;

    isSubmitting = true;
    processBtn.disabled = true;
    processBtn.textContent = 'Uploading...';

    if (currentJobId !== null) {
        cancelPolling();
        jobCard.hidden = true;
        resultsSection.hidden = true;
        currentJobId = null;
        currentStatusUrl = null;
    }

    const formData = new FormData();
    formData.append('image', selectedFile);

    try {
        const response = await fetch(`${API_BASE}/v1/images`, {
            method: 'POST',
            body: formData,
        });

        const data = await response.json();

        if (!response.ok) {
            const msg = data.error || 'Upload failed';
            if (typeof msg === 'object') {
                const errors = Object.values(msg).join(', ');
                showError(errors);
            } else {
                showError(msg);
            }
            return;
        }

        currentJobId = data.job_id;
        currentStatusUrl = data.status_url;

        showJobCard(data);
        showResultsSection();
        startPolling();

    } catch (err) {
        showError('Network error: ' + err.message);
    } finally {
        isSubmitting = false;
        processBtn.textContent = 'Process image';
        processBtn.disabled = true;
    }
}

function showJobCard(data) {
    jobCard.hidden = false;
    jobIdLabel.textContent = `Job #${data.job_id}`;
    updateStatusBadge(data.status);
    updateTimeline(data.status, null, null, null);
    pollingInfo.hidden = false;
    tryAgainSection.hidden = true;
}

function updateStatusBadge(status) {
    statusBadge.textContent = status;
    statusBadge.className = 'badge badge-' + status;
}

function updateTimeline(status, queuedAt, startedAt, completedAt) {
    const now = new Date().toLocaleTimeString();
    let html = '';

    html += `<div class="timeline-item complete">
        <div>Upload accepted</div>
        <div class="timeline-time">${now}</div>
    </div>`;

    html += `<div class="timeline-item complete">
        <div>Original stored</div>
        <div class="timeline-time">${now}</div>
    </div>`;

    if (status === 'queued' || status === 'processing' || status === 'completed') {
        const cls = status === 'completed' ? 'complete' : 'active';
        const time = startedAt ? formatTime(startedAt) : '';
        html += `<div class="timeline-item ${cls}">
            <div>Generating variants</div>
            ${time ? `<div class="timeline-time">${time}</div>` : ''}
        </div>`;
    }

    if (status === 'completed') {
        const time = completedAt ? formatTime(completedAt) : now;
        html += `<div class="timeline-item complete">
            <div>Complete</div>
            <div class="timeline-time">${time}</div>
        </div>`;
    } else if (status === 'failed') {
        html += `<div class="timeline-item" style="color: #c0392b;">
            <div>Failed</div>
        </div>`;
    } else {
        html += `<div class="timeline-item pending">
            <div>Complete</div>
        </div>`;
    }

    jobTimeline.innerHTML = html;
}

function showResultsSection() {
    resultsSection.hidden = false;
    resultsEmpty.hidden = false;
    resultsActive.hidden = false;
    resultsGrid.hidden = true;
    resultsError.hidden = true;
}

function showResults(job) {
    resultsEmpty.hidden = true;
    resultsActive.hidden = true;
    resultsError.hidden = true;
    resultsGrid.hidden = false;
    resultsGrid.innerHTML = '';

    if (!job.variants || job.variants.length === 0) {
        return;
    }

    for (const variant of job.variants) {
        const card = document.createElement('div');
        card.className = 'result-card';
        card.innerHTML = `
            <img src="${API_BASE}${variant.url}" alt="${variant.name}">
            <div class="result-info">
                <div class="result-name">${variant.name}</div>
                <div class="result-dims">${variant.width} x ${variant.height}</div>
                <div class="result-actions">
                    <a href="${API_BASE}${variant.url}" target="_blank">View</a>
                    <a href="${API_BASE}${variant.url}" download>Download</a>
                </div>
            </div>
        `;
        resultsGrid.appendChild(card);
    }
}

function showFailedResults(message) {
    resultsEmpty.hidden = true;
    resultsActive.hidden = true;
    resultsGrid.hidden = true;
    resultsError.hidden = false;
    resultsErrorMessage.textContent = message || 'An unknown error occurred.';
}

function startPolling() {
    cancelPolling();

    pollingController = new AbortController();

    pollingInterval = setInterval(async () => {
        if (!currentStatusUrl) {
            cancelPolling();
            return;
        }

        try {
            const response = await fetch(`${API_BASE}${currentStatusUrl}`, {
                signal: pollingController.signal,
            });

            if (!response.ok) {
                stopPollingWithError();
                return;
            }

            const job = await response.json();
            updateStatusBadge(job.status);
            updateTimeline(job.status, job.queued_at, job.started_at, job.completed_at);

            if (job.status === 'completed') {
                cancelPolling();
                pollingInfo.hidden = true;
                showResults(job);
            } else if (job.status === 'failed') {
                cancelPolling();
                pollingInfo.hidden = true;
                showFailedResults(job.error);
            }

        } catch (err) {
            if (err.name === 'AbortError') return;
            stopPollingWithError();
        }
    }, 1000);
}

function cancelPolling() {
    if (pollingInterval) {
        clearInterval(pollingInterval);
        pollingInterval = null;
    }
    if (pollingController) {
        pollingController.abort();
        pollingController = null;
    }
}

function stopPollingWithError() {
    cancelPolling();
    pollingInfo.hidden = true;
    tryAgainSection.hidden = false;
}

function tryAgain() {
    tryAgainSection.hidden = true;

    if (currentStatusUrl) {
        pollingInfo.hidden = false;
        startPolling();
    }
}

function formatFileSize(bytes) {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
}

function formatTime(isoStr) {
    if (!isoStr) return '';
    return new Date(isoStr).toLocaleTimeString();
}
