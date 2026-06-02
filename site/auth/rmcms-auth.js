import { decryptRSA } from "../crypto-utils.js";

// Private module variables
let sessionKey = null;
let sessionUserHash = null;

const DB_NAME = "rmcms_auth_db";
const STORE_NAME = "session";

async function getDB() {
    return new Promise((resolve, reject) => {
        const request = indexedDB.open(DB_NAME, 1);
        request.onupgradeneeded = () => {
            if (!request.result.objectStoreNames.contains(STORE_NAME)) {
                request.result.createObjectStore(STORE_NAME);
            }
        };
        request.onsuccess = () => resolve(request.result);
        request.onerror = () => reject(request.error);
    });
}

export const rmcmsAuth = {
    async unlock(privateKey, usernameHash) {
        sessionKey = privateKey;
        sessionUserHash = usernameHash;

        try {
            const db = await getDB();
            const tx = db.transaction(STORE_NAME, "readwrite");
            const store = tx.objectStore(STORE_NAME);
            store.put(privateKey, "privateKey");
            store.put(usernameHash, "usernameHash");
        } catch (e) {
            console.error("Persistence failed:", e);
        }
    },

    async tryRestoreSession() {
        if (sessionKey) return true;
        try {
            const db = await getDB();
            const tx = db.transaction(STORE_NAME, "readonly");
            const store = tx.objectStore(STORE_NAME);
            const keyReq = store.get("privateKey");
            const hashReq = store.get("usernameHash");

            return new Promise((resolve) => {
                tx.oncomplete = () => {
                    if (keyReq.result && hashReq.result) {
                        sessionKey = keyReq.result;
                        sessionUserHash = hashReq.result;
                        resolve(true);
                    } else {
                        resolve(false);
                    }
                };
            });
        } catch (e) { return false; }
    },

    isAuthenticated() {
        return sessionKey !== null;
    },

    async getProperty(filename, propertyName) {
        if (!sessionKey) await this.tryRestoreSession();
        if (!sessionKey) {
            alert("Authentication required.");
            throw new Error("Not authenticated.");
        }

        try {
            const response = await fetch(`auth_data/${filename}`);
            const data = await response.json();
            const userEntry = data.users[sessionUserHash];
            if (!userEntry) throw new Error(`User entry missing.`);
            const encryptedValue = userEntry[propertyName];
            return await decryptRSA(encryptedValue, sessionKey);
        } catch (e) {
            alert(`Auth Error: ${e.message}`);
            throw e; 
        }
    },

    async logout() {
        const db = await getDB();
        const tx = db.transaction(STORE_NAME, "readwrite");
        tx.objectStore(STORE_NAME).clear();
        sessionKey = null;
        sessionUserHash = null;
    }
};
