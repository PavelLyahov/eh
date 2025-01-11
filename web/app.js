const API_BASE = "/api/phrases";

async function loadPhrases() {
    const response = await fetch(API_BASE);
    const phrases = await response.json();

    const phrasesList = document.getElementById("phrasesList");
    phrasesList.innerHTML = "";

    phrases.forEach((phrase) => {
        const li = document.createElement("li");
        li.innerHTML = `
            <span>
                <strong>${phrase.text}</strong> — ${phrase.translation} (Level: ${phrase.level})
            </span>
            <button onclick="deletePhrase(${phrase.id})">Delete</button>
        `;
        phrasesList.appendChild(li);
    });
}

async function addPhrase(event) {
    event.preventDefault();

    const text = document.getElementById("text").value;
    const translation = document.getElementById("translation").value;
    const level = document.getElementById("level").value;
    const file = fileSelector.value;

    const response = await fetch(`/api/files/phrases?file=${file}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ text, translation, level: parseInt(level) }),
    });

    if (response.ok) {
        document.getElementById("addPhraseForm").reset();
        loadPhrases();
    } else {
        alert("Failed to add phrase.");
    }
}

async function deletePhrase(id) {
    const response = await fetch(`${API_BASE}/?id=${id}`, {
        method: "DELETE",
    });

    if (response.ok) {
        loadPhrases();
    } else {
        alert("Failed to delete phrase.");
    }
}

const fileSelector = document.getElementById("fileSelector");

async function loadFiles() {
    const response = await fetch("/api/files");
    const files = await response.json();

    fileSelector.innerHTML = files.map(
        (file) => `<option value="${file}">${file}</option>`
    ).join("");
}

async function createFile(event) {
    event.preventDefault();
    const newFileName = document.getElementById("newFileName").value;

    const response = await fetch("/api/files", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ fileName: newFileName }),
    });

    if (response.ok) {
        document.getElementById("createFileForm").reset();
        loadFiles();
    } else {
        alert("Failed to create file.");
    }
}



// Инициализация
document.getElementById("createFileForm").addEventListener("submit", createFile);
loadFiles();


document.getElementById("addPhraseForm").addEventListener("submit", addPhrase);

loadPhrases();
