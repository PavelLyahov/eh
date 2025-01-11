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

    const response = await fetch(API_BASE, {
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

document.getElementById("addPhraseForm").addEventListener("submit", addPhrase);

loadPhrases();
