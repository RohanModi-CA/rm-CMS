import { base64ToBuffer, hashString } from "./auth/crypto-utils.js";
import { rmcmsAuth } from "./auth/rmcms-auth.js";

async function deriveAESKey(password, salt) {
    const saltBuffer = new TextEncoder().encode(salt);
    const hash = await hashwasm.argon2id({
        password: password, salt: saltBuffer, parallelism: 1, iterations: 2, memorySize: 16384, hashLength: 32, outputType: 'binary'
    });
    return await crypto.subtle.importKey("raw", hash, "AES-GCM", false, ["unwrapKey"]);
}

function updateUI() {
    const authStatus = document.getElementById("auth-status");
    const logoutButton = document.getElementById("logout-button");
    const loginFields = document.getElementById("login-form-fields");

    if (rmcmsAuth.isAuthenticated()) {
        authStatus.textContent = "Status: Authenticated";
        authStatus.style.color = "green";
        logoutButton.classList.remove("hidden");
        loginFields.classList.add("hidden");
    } else {
        authStatus.textContent = "Status: Not Authenticated";
        authStatus.style.color = "red";
        logoutButton.classList.add("hidden");
        loginFields.classList.remove("hidden");
    }
}

document.addEventListener("DOMContentLoaded", async () => {
    await rmcmsAuth.tryRestoreSession();
    updateUI();

    document.getElementById("logout-button").addEventListener("click", async () => {
        await rmcmsAuth.logout();
        updateUI();
    });

    document.getElementById("login-button").addEventListener("click", async () => {
        const username = document.getElementById("login-username").value;
        const password = document.getElementById("login-password").value;
        if (!username || !password) return;

        try {
            const response = await fetch("/auth/keys.json");
            const data = await response.json();
            const usernameHash = await hashString(username);
            const user = data.users.find(u => u.username === usernameHash);
            if (!user) return alert("User not found");

            const aesKey = await deriveAESKey(password, user.salt);
            const iv = await base64ToBuffer(user.vault.iv);
            const wrappedKey = await base64ToBuffer(user.vault.wrapped_key);

            const privateKey = await crypto.subtle.unwrapKey(
                "pkcs8", wrappedKey, aesKey, { name: "AES-GCM", iv: iv },
                { name: "RSA-OAEP", hash: "SHA-256" }, true, ["decrypt"]
            );

            const authCheckBuffer = await base64ToBuffer(user.auth_check);
            const decryptedBuffer = await crypto.subtle.decrypt({ name: "RSA-OAEP" }, privateKey, authCheckBuffer);
            if (new TextDecoder().decode(decryptedBuffer) === "SUCCESS") {
                await rmcmsAuth.unlock(privateKey, usernameHash);
                updateUI();
            } else { alert("Authentication Failed"); }
        } catch (e) { alert("Login Failed."); }
    });
});
