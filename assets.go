package entydad

import "embed"

// AssetsFS embeds this package's static CSS/JS so the app can copy them at boot via
// pyeza.CopyNamespacedAssets — replaces the old CopyStyles/CopyStaticAssets + runtime.Caller hack.
//
//go:embed assets
var AssetsFS embed.FS
