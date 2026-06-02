import { bufferToBase64, base64ToBuffer, hashString } from "./crypto-utils.js";
import { rmcmsAuth } from "./auth/rmcms-auth.js";

async function deriveAESKey(password, salt) {
    const saltBuffer = new TextEncoder().encode(salt);
    const hash = await hashwasm.argon2id({
        password: password,
        salt: saltBuffer,
        parallelism: 1,
        iterations: 2,
        memorySize: 16384,
        hashLength: 32,
        outputType: 'binary'
    });
    return await crypto.subtle.importKey(
        "raw", hash, "AES-GCM", false, ["encrypt", "decrypt", "wrapKey", "unwrapKey"]
    );
}

function updateUI() {
    const authStatus = document.getElementById("auth-status");
    const logoutButton = document.getElementById("logout-button");
    const loginButton = document.getElementById("login-button");
    const usernameInput = document.getElementById("login-username");
    const passwordInput = document.getElementById("login-password");

    if (rmcmsAuth.isAuthenticated()) {
        authStatus.textContent = "Status: Authenticated";
        authStatus.style.color = "green";
        logoutButton.classList.remove("hidden");
        loginButton.classList.add("hidden");
        usernameInput.classList.add("hidden");
        passwordInput.classList.add("hidden");
    } else {
        authStatus.textContent = "Status: Not Authenticated";
        authStatus.style.color = "red";
        logoutButton.classList.add("hidden");
        loginButton.classList.remove("hidden");
        usernameInput.classList.remove("hidden");
        passwordInput.classList.remove("hidden");
    }
}

document.addEventListener("DOMContentLoaded", async () => {
    const showLogin = document.getElementById("show-login");
    const showSignup = document.getElementById("show-signup");
    const loginForm = document.getElementById("login-form");
    const signupForm = document.getElementById("signup-form");

    await rmcmsAuth.tryRestoreSession();
    updateUI();

    showLogin.addEventListener("click", () => {
        loginForm.classList.remove("hidden");
        signupForm.classList.add("hidden");
    });

    showSignup.addEventListener("click", () => {
        signupForm.classList.remove("hidden");
        loginForm.classList.add("hidden");
    });

    document.getElementById("logout-button").addEventListener("click", async () => {
        await rmcmsAuth.logout();
        updateUI();
    });

    document.getElementById("signup-button").addEventListener("click", async () => {
        const username = document.getElementById("signup-username").value;
        const password = document.getElementById("signup-password").value;
        const output = document.getElementById("signup-json-output");
        if (!username || !password) { alert("Enter username/password"); return; }
        try {
            const saltBuffer = new Uint8Array(8);
            crypto.getRandomValues(saltBuffer);
            const salt = Array.from(saltBuffer).map(b => b.toString(16).padStart(2, '0')).join('');
            const aesKey = await deriveAESKey(password, salt);
            const rsaPair = await crypto.subtle.generateKey(
                { name: "RSA-OAEP", modulusLength: 2048, publicExponent: new Uint8Array([1, 0, 1]), hash: "SHA-256" },
                true, ["encrypt", "decrypt", "wrapKey", "unwrapKey"]
            );
            const pubKeyBuffer = await crypto.subtle.exportKey("spki", rsaPair.publicKey);
            const pubKeyBase64 = await bufferToBase64(pubKeyBuffer);
            const enc = new TextEncoder();
            const authCheckBuffer = await crypto.subtle.encrypt({ name: "RSA-OAEP" }, rsaPair.publicKey, enc.encode("SUCCESS"));
            const authCheckBase64 = await bufferToBase64(authCheckBuffer);
            const iv = crypto.getRandomValues(new Uint8Array(12));
            const wrappedPrivateKey = await crypto.subtle.wrapKey("pkcs8", rsaPair.privateKey, aesKey, { name: "AES-GCM", iv: iv });
            const wrappedBase64 = await bufferToBase64(wrappedPrivateKey);
            const ivBase64 = await bufferToBase64(iv);
            const usernameHash = await hashString(username);
            const userObject = {
                username: usernameHash, salt: salt, auth_check: authCheckBase64, public_key: pubKeyBase64,
                vault: { iv: ivBase64, wrapped_key: wrappedBase64 }, "encrypted-keys": {}
            };
            output.value = JSON.stringify(userObject, null, 4);
        } catch (e) { alert("Error generating keys."); }
    });

    document.getElementById("login-button").addEventListener("click", async () => {
        const username = document.getElementById("login-username").value;
        const password = document.getElementById("login-password").value;
        if (!username || !password) return;
        try {
            const response = await fetch("keys.json");
            const data = await response.json();
            const usernameHash = await hashString(username);
            const user = data.users.find(u => u.username === usernameHash);
            if (!user) { alert("User not found"); return; }
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
                alert("Login Successful!");
                await rmcmsAuth.unlock(privateKey, usernameHash);
                updateUI();
            } else { alert("Authentication Failed"); }
        } catch (e) { alert("Login Failed."); }
    });
});
