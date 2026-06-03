import { base64ToBuffer, bufferToBase64 } from "./crypto-utils.js";

async function loadUsers() {
    const userSelect = document.getElementById("user-select");
    try {
        const response = await fetch("/auth/keys.json");
        const data = await response.json();
        
        userSelect.innerHTML = "";
        data.users.forEach(user => {
            const opt = document.createElement("option");
            opt.value = user.public_key;
            opt.dataset.hash = user.username;
            opt.textContent = `User: ${user.username.substring(0, 8)}...`;
            userSelect.appendChild(opt);
        });
    } catch (e) {
        userSelect.innerHTML = "<option>Error loading keys.json</option>";
    }
}

async function encryptSecret() {
    const pubKeyBase64 = document.getElementById("user-select").value;
    const userHash = document.getElementById("user-select").selectedOptions[0].dataset.hash;
    const propName = document.getElementById("prop-name").value;
    const propValue = document.getElementById("prop-value").value;
    const output = document.getElementById("output");

    if (!propName || !propValue) return alert("Enter name and value.");

    try {
        // Import Public Key
        const pubKeyBuffer = await base64ToBuffer(pubKeyBase64);
        const publicKey = await crypto.subtle.importKey(
            "spki",
            pubKeyBuffer,
            { name: "RSA-OAEP", hash: "SHA-256" },
            true,
            ["encrypt"]
        );

        // Encrypt
        const enc = new TextEncoder();
        const encryptedBuffer = await crypto.subtle.encrypt(
            { name: "RSA-OAEP" },
            publicKey,
            enc.encode(propValue)
        );

        const encryptedBase64 = await bufferToBase64(encryptedBuffer);

        // Generate JSON Structure
        const result = {
            "users": {
                [userHash]: {
                    [propName]: encryptedBase64
                }
            }
        };

        output.value = JSON.stringify(result, null, 4);
        alert("Success! Copy the JSON into a file in /auth/data_jsons/");
    } catch (e) {
        console.error(e);
        alert("Encryption failed.");
    }
}

document.addEventListener("DOMContentLoaded", () => {
    loadUsers();
    document.getElementById("encrypt-btn").onclick = encryptSecret;
});
