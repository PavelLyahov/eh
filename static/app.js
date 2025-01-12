const API_BASE = "/api";
let selectedFile = null;

window.renameFile = renameFile;
window.deleteFile = deleteFile;
window.selectFile = selectFile;
window.addPhrase = addPhrase;
window.deletePhrase = deletePhrase;

async function loadFiles() {
    const response = await fetch(`${API_BASE}/files`);
    const files = await response.json();

    const fileList = document.getElementById("fileList");
    fileList.innerHTML = "";

    files.forEach((file) => {
        const li = document.createElement("li");
        li.className = selectedFile === file ? "selected" : "";
        li.innerHTML = `
            <span>${file}</span>
            <div>
                <button onclick="renameFile('${file}')">Rename</button>
                <button onclick="deleteFile('${file}')">Delete</button>
                <button onclick="selectFile('${file}')">Open</button>
            </div>
        `;
        fileList.appendChild(li);
    });
}

async function createFile(event) {
    event.preventDefault();

    const newFileName = document.getElementById("newFileName").value;

    const response = await fetch(`${API_BASE}/files`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ fileName: newFileName }),
    });

    if (response.ok) {
        document.getElementById("createFileForm").reset();
        await loadFiles();
    } else {
        handleError(response);
    }
}

async function renameFile(file) {
    const newName = prompt("Enter new name:", file);
    if (!newName) return;

    const response = await fetch(`${API_BASE}/files?oldName=${file}&newName=${newName}`, {
        method: "PATCH",
    });

    if (response.ok) {
        await loadFiles();
    } else {
        handleError(response);
    }
}

async function deleteFile(file) {
    if (!confirm(`Are you sure you want to delete ${file}?`)) return;

    const response = await fetch(`${API_BASE}/files?file=${file}`, {
        method: "DELETE",
    });

    if (response.ok) {
        if (selectedFile === file) {
            selectedFile = null;
            document.getElementById("phraseManager").style.display = "none";
        }
        await loadFiles();
    } else {
        handleError(response);
    }
}

async function selectFile(file) {
    selectedFile = file;
    document.getElementById("phraseManager").style.display = "block";
    document.getElementById("currentFileName").textContent = file;

    await loadPhrases(file);
    await loadFiles();
}

function handleError(response) {
    (errorHandlers[response.status] || errorHandlers.default)();
}

const errorHandlers = {
    400: () => alert("Bad request. Please check the input and try again."),
    404: () => alert("Requested resource not found."),
    409: () => alert("Conflict. The file already exists."),
    500: () => alert("Internal server error. Please try again later."),
    default: () => alert("An unknown error occurred."),
};


async function loadPhrases(file) {
    const response = await fetch(`${API_BASE}/files/phrases?file=${file}`);
    const phrases = await response.json();

    const phrasesList = document.getElementById("phrasesList");
    phrasesList.innerHTML = "";

    phrases.forEach((phrase) => {
        const li = document.createElement("li");
        li.innerHTML = `
            <span>
                <strong>${phrase.text}</strong> — ${phrase.translation}
            </span>
            <button onclick="deletePhrase('${file}', ${phrase.id})">Delete</button>
        `;
        phrasesList.appendChild(li);
    });
}

async function addPhrase(event) {
    event.preventDefault();

    if (!selectedFile) {
        alert("No file selected!");
        return;
    }

    const text = document.getElementById("text").value;
    const translation = document.getElementById("translation").value;
    const level = document.getElementById("level").value;

    const response = await fetch(`${API_BASE}/files/phrases?file=${selectedFile}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ text, translation, level: parseInt(level) }),
    });

    if (response.ok) {
        document.getElementById("addPhraseForm").reset();
        await loadPhrases(selectedFile);
    } else {
        alert("Failed to add phrase.");
    }
}

async function deletePhrase(file, id) {
    const response = await fetch(`${API_BASE}/files/phrases?id=${id}&file=${file}`, {
        method: "DELETE",
    });

    if (response.ok) {
        await loadPhrases(file);
    } else {
        alert("Failed to delete phrase.");
    }
}

document.querySelectorAll(".collapsible-header").forEach((header) => {
    header.addEventListener("click", () => {
        const section = header.parentElement;
        section.classList.toggle("open");
    });
});

document.getElementById("createFileForm").addEventListener("submit", createFile);
document.getElementById("addPhraseForm").addEventListener("submit", addPhrase);

loadFiles();
