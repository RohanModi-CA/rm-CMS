import { rmcmsAuth } from "/auth/rmcms-auth.js";

let stylesAdded = false;

function addStyles() {
    if (stylesAdded) return;
    stylesAdded = true;

    const link = document.createElement("link");
    link.rel = "stylesheet";
    link.href = "/elements/global.css";
    document.head.appendChild(link);

}

function setAuthIndicator(link, authenticated) {
    link.textContent = authenticated ? "Logged in" : "Log in";
    link.dataset.state = authenticated ? "authenticated" : "unauthenticated";
}

export function mountHeader(container, options = {}) {
    if (!container) throw new Error("mountHeader needs a container element.");

    addStyles();

    const homeHref = options.homeHref || "/convert.html";
    const loginHref = options.loginHref || "/login.html";
    const homeLabel = options.homeLabel || "Home";

    const header = document.createElement("header");
    header.className = "rm-elements-header";

    const left = document.createElement("div");
    left.className = "rm-elements-header-left";

    const homeLink = document.createElement("a");
    homeLink.href = homeHref;
    homeLink.textContent = homeLabel;
    left.appendChild(homeLink);

    const right = document.createElement("div");
    right.className = "rm-elements-header-right";

    const authLink = document.createElement("a");
    authLink.href = loginHref;
    authLink.target = "_blank";
    authLink.rel = "noopener noreferrer";
    authLink.className = "rm-elements-auth-indicator";
    authLink.dataset.state = "loading";
    authLink.textContent = "Checking...";
    right.appendChild(authLink);

    header.appendChild(left);
    header.appendChild(right);
    container.replaceChildren(header);

    (async () => {
        try {
            await rmcmsAuth.tryRestoreSession();
            setAuthIndicator(authLink, rmcmsAuth.isAuthenticated());
        } catch {
            setAuthIndicator(authLink, false);
        }
    })();
}
