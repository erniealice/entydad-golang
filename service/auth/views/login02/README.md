# Login 02 View

This directory owns the current Entydad login-page view and its embedded HTML.
It supports classic password mode and Firebase mode from one presentation model.

## Key files

- `page.go` defines the page dependencies and data, including carousel slides,
  password-form visibility, signup availability, social providers, and public
  Firebase browser configuration.
- `action.go` provides the lightweight view action contract; real credential and
  session work remains in the parent auth module's handlers.
- `embed.go` exposes `templates/*.html` to the host renderer.

In classic mode the password form is always shown. In Firebase mode the parent
module supplies allowed social providers and decides whether password sign-in is
also available. Microsoft and Google buttons are therefore composition-driven,
not unconditional markup.

The view uses Entydad label models and Pyeza page contracts so branding and
business-specific text can be supplied by each host application.
