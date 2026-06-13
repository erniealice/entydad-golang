package auth

// noAccessHTML is the static page served at GET /auth/no-access. It is
// reached after a successful credential check that found zero active
// principal bindings — the user IS valid but has nothing to do here. The
// cookie has been cleared before this handler runs (see domain_auth login
// handler, no-principal branch).
//
// Kept as an inline string to mirror the signout-loading page: no template
// dependencies, no data, no flash; same bytes every render.
const noAccessHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>No access</title>
    <style>
        :root { color-scheme: light dark; }
        html, body { height: 100%; margin: 0; }
        body {
            display: flex;
            align-items: center;
            justify-content: center;
            font-family: system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background: #F5F5F7;
            color: #1F1F23;
        }
        .no-access-card {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 1rem;
            padding: 2rem 2.5rem;
            text-align: center;
            max-width: 28rem;
        }
        .no-access-title { font-size: 1.25rem; font-weight: 600; margin: 0; }
        .no-access-body  { font-size: 0.95rem; color: #555; margin: 0; line-height: 1.5; }
        .no-access-link  { font-size: 0.875rem; color: #2563eb; text-decoration: underline; margin-top: 0.5rem; }
        @media (prefers-color-scheme: dark) {
            body { background: #1F1F23; color: #F5F5F7; }
            .no-access-body { color: #C9C9CC; }
            .no-access-link { color: #6FA8DC; }
        }
    </style>
</head>
<body>
    <main class="no-access-card" role="status" data-testid="no-access-page">
        <p class="no-access-title">No access yet</p>
        <p class="no-access-body">Your account is recognised, but no workspace, portal grant, or delegation is currently active for you. Ask an administrator to grant you access, then sign in again.</p>
        <a class="no-access-link" href="/auth/login" data-testid="no-access-back-to-login">Back to sign in</a>
    </main>
</body>
</html>`

// logoutLoadingHTML is the tiny standalone page rendered at GET /auth/logout.
// It auto-submits a POST to /action/auth/logout on load, so a user hitting
// /auth/logout in a browser tab completes the same flow as the sidebar
// profile menu form.
//
// Kept as an inline string (rather than a template file) because:
//   - No layout dependencies (the app-shell session may already be invalid).
//   - Zero query parameters or dynamic data — same bytes every render.
//   - Standalone styling so the transition feels intentional, not a flash.
const logoutLoadingHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Signing out…</title>
    <style>
        :root { color-scheme: light dark; }
        html, body { height: 100%; margin: 0; }
        body {
            display: flex;
            align-items: center;
            justify-content: center;
            font-family: system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background: #F5F5F7;
            color: #1F1F23;
        }
        .signout-card {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 1rem;
            padding: 2rem 2.5rem;
            text-align: center;
        }
        .signout-spinner {
            width: 2.25rem;
            height: 2.25rem;
            border: 3px solid rgba(0, 0, 0, 0.1);
            border-top-color: #2E5C8A;
            border-radius: 50%;
            animation: signout-spin 0.8s linear infinite;
        }
        .signout-title {
            font-size: 1.0625rem;
            font-weight: 500;
            margin: 0;
        }
        .signout-subtitle {
            font-size: 0.875rem;
            color: #6b6b6b;
            margin: 0;
            line-height: 1.5;
        }
        .signout-fallback-link {
            margin-top: 1rem;
            font-size: 0.8125rem;
            color: #6b6b6b;
            text-decoration: underline;
        }
        @keyframes signout-spin { to { transform: rotate(360deg); } }
        @media (prefers-color-scheme: dark) {
            body { background: #1F1F23; color: #F5F5F7; }
            .signout-spinner { border-color: rgba(255,255,255,0.12); border-top-color: #6FA8DC; }
            .signout-subtitle, .signout-fallback-link { color: #B5B5B8; }
        }
    </style>
</head>
<body>
    <main class="signout-card" role="status" aria-live="polite">
        <div class="signout-spinner" aria-hidden="true"></div>
        <p class="signout-title">Signing you out…</p>
        <p class="signout-subtitle">Hang tight — we're ending your session and returning to the sign-in page.</p>
        <form id="signoutForm" method="POST" action="/action/auth/logout" hidden>
            <noscript>
                <button type="submit">Continue</button>
            </noscript>
        </form>
        <noscript>
            <p class="signout-fallback-link">JavaScript is disabled — press the Continue button above to finish signing out.</p>
        </noscript>
    </main>
    <script>
        (function () {
            var form = document.getElementById('signoutForm');
            if (form) { form.submit(); }
        }());
    </script>
</body>
</html>`
