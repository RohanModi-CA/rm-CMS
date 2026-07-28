import { bufferToBase64, hashString } from "./auth/crypto-utils.js";

async function deriveAESKey(password, salt) {
    const saltBuffer = new TextEncoder().encode(salt);
    const hash = await hashwasm.argon2id({
        password: password, salt: saltBuffer, parallelism: 1, iterations: 2, memorySize: 16384, hashLength: 32, outputType: 'binary'
    });
    return await crypto.subtle.importKey("raw", hash, "AES-GCM", false, ["wrapKey"]);
}

document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("signup-button").addEventListener("click", async () => {
        const username = document.getElementById("signup-username").value;
        const password = document.getElementById("signup-password").value;
        if (!username || !password) return alert("Enter username and password");

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
            
            const userObject = {
                username: await hashString(username),
                salt: salt,
                auth_check: authCheckBase64,
                public_key: pubKeyBase64,
                vault: { iv: await bufferToBase64(iv), wrapped_key: await bufferToBase64(wrappedPrivateKey) },
                "encrypted-keys": {}
            };
            document.getElementById("signup-json-output").value = JSON.stringify(userObject, null, 4);
        } catch (e) { alert("Error generating keys."); }
    });
});
