# RM-CMS Auth

`rmcms-auth.js` is the small client-side auth runtime for this project.

It handles three things:

- restoring a session from IndexedDB
- unlocking a session after login
- decrypting per-user properties from JSON files

## Setup

The auth flow depends on browser crypto support and IndexedDB.

Login and signup also use `hash-wasm` for Argon2id.

## Login Flow

The login page uses `login.js`.

The flow is:

1. the user enters a username and password
2. the username is hashed with SHA-256
3. `/auth/keys.json` is loaded
4. the user record is found
5. Argon2id derives an AES key from the password and the stored salt
6. the wrapped RSA private key is unwrapped
7. the private key is checked by decrypting `auth_check`
8. on success, the session is stored with `rmcmsAuth.unlock()`

Example:

```js
import { rmcmsAuth } from "./auth/rmcms-auth.js";

await rmcmsAuth.tryRestoreSession();
if (rmcmsAuth.isAuthenticated()) {
    console.log("session restored");
}
```

## Signup Flow

The signup page uses `signup.js`.

The flow is:

1. the user picks a username and password
2. a random salt is generated
3. Argon2id derives an AES key
4. an RSA keypair is generated
5. the private key is wrapped with AES-GCM
6. the public key and the wrapped private key are written to JSON

The generated object currently contains:

- `username`
- `salt`
- `auth_check`
- `public_key`
- `vault.iv`
- `vault.wrapped_key`
- `encrypted-keys`

## Session Storage

`rmcmsAuth.unlock(privateKey, usernameHash)` stores two values:

- `privateKey`
- `usernameHash`

They are saved in IndexedDB under the `rmcms_auth_db` database and `session`
object store.

`rmcmsAuth.tryRestoreSession()` reloads them into module state.

`rmcmsAuth.logout()` clears both memory and IndexedDB.

## Secret Retrieval

`rmcmsAuth.getProperty(filename, propertyName)` fetches a secret JSON file and
decrypts one property for the current user.

The current lookup path is:

```text
/auth/data_jsons/<filename>
```

The file is expected to contain a `users` object keyed by the hashed username.

Example:

```js
import { rmcmsAuth } from "/auth/rmcms-auth.js";

const secret = await rmcmsAuth.getProperty("test.json", "my_secret");
console.log(secret);
```

## Secret Manager

`auth/manager.js` is a helper for creating encrypted values.

It loads `/auth/keys.json`, lets the operator pick a user, encrypts a plaintext
value with that user’s public key, and outputs JSON for `/auth/data_jsons/`.

## Current Shape Notes

The repo currently has two different expectations for `keys.json`:

- `login.js` and `auth/manager.js` expect `data.users` to be an array
- `rmcmsAuth.getProperty()` expects `data.users` to be an object keyed by hash

The docs here reflect that current code, but the data model should be made
consistent before the auth flow can be considered finished.

