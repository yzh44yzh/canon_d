# Fyne Android build and packaging constraints

**Research date:** 2026-11-27

**Scope:** Canon D on Android, Fyne mobile packaging, a developer-built debug APK, and Go-embedded Markdown decks.

**Evidence status:** facts below come from first-party Fyne, Go, Android, and Canon D sources. No Android SDK/NDK build or physical-device test was run for this note; items that cannot be established from primary sources are marked as unresolved.

## Recommendation

Use a single repository Go module and a dedicated Fyne main package, for example `cmd/mobile`, with `FyneApp.toml`, `Icon.png`, and the app sources there.

| Component | Concrete choice | Basis |
| --- | --- | --- |
| Canon D Go | **Go 1.25.0 or newer**, matching this repository's `go 1.25.0`. | Canon D `go.mod`, `AGENTS.md`; see project contract. |
| Fyne library | **`fyne.io/fyne/v2` v2.8.2**, pinned in `go.mod`. | Latest Fyne release inspected; its `go.mod` requires Go 1.22.0. [Fyne v2.8.2](https://github.com/fyne-io/fyne/releases/tag/v2.8.2); [v2.8.2 `go.mod`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.2/go.mod) |
| Fyne CLI | **`fyne.io/tools/cmd/fyne` v1.7.4**, pinned with `go install`. | Fyne moved the command to `fyne.io/tools/cmd/fyne`; the inspected CLI module requires Go 1.22.0 and itself pins Fyne v2.8.1. The application can still use Fyne v2.8.2. [Fyne tools v1.7.4](https://github.com/fyne-io/tools/releases/tag/v1.7.4); [CLI installation](https://docs.fyne.io/started/quick/); [tools `go.mod`](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/go.mod) |
| CGO | **`CGO_ENABLED=1`.** | Fyne uses CGo for native graphics and its mobile shared library. The CLI explicitly enables CGO and `-buildmode=c-shared`, supplies NDK `CC`/`CXX`, `GOOS=android`, and per-architecture `GOARCH`. [Fyne cross-compiling](https://docs.fyne.io/started/cross-compiling/); [mobile environment source](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/env.go); [shared-library build source](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go) |
| Android target | **Android API level 21 or newer; `android` (all supported ABIs) for a portable APK, or `android/arm64` for a smaller arm64-only test APK.** | The Fyne CLI rejects `-androidapi < 21`, emits `minSdkVersion=21`, and emits `targetSdkVersion=29` for non-release packaging. `fyne package` has no documented `-androidapi` flag; the minimum is fixed in this CLI version. [build constants](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/build.go); [API-level contract](https://developer.android.com/guide/topics/manifest/uses-sdk-element) |
| Android SDK | Installed Android SDK, `ANDROID_HOME`, at least one `$ANDROID_HOME/platforms/android-*`, and `adb` in `PATH`. | Fyne mobile packaging requires Android SDK; the CLI needs `ANDROID_HOME`, reads `$ANDROID_HOME/platforms`, finds `adb`, and emits `android:minSdkVersion`/`targetSdkVersion`. `ANDROID_SDK_ROOT` is not an alternative read by this source. [Fyne Mobile Packaging](https://docs.fyne.io/started/mobile/); [platform resource lookup](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/binres/sdk.go); [Android ADB](https://developer.android.com/tools/adb) |
| Android NDK | Installed matching NDK and **`ANDROID_NDK_HOME`**. NDK must have LLVM `clang`, `clang++`, and `nm` under `toolchains/llvm/prebuilt/<host>/bin`; source checks NDK **r19c or newer**. | `$ANDROID_HOME/ndk-bundle` is a fallback, but explicit `ANDROID_NDK_HOME` is clearer and handles side-by-side NDK installations. The CLI itself supplies compiler variables; no separate `gomobile init` is part of Fyne packaging. [NDK lookup and compiler checks](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/env.go); [Android NDK](https://developer.android.com/ndk/downloads) |
| Compile platform | Any installed SDK `platforms/android-*`; this CLI chooses the lexicographically greatest directory and reads its `android.jar`. Keep it pinned in CI for reproducibility. | `LatestAPIResourcesPath`. [platform resource lookup](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/binres/sdk.go) |
| ABIs | `arm`, `arm64`, `386`, `amd64` are built by `android`; map to `armeabi-v7a`, `arm64-v8a`, `x86`, `x86_64`. | `allArchs` and `ndk` in the CLI. `arm` uses `GOARM=7`. [mobile environment source](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/env.go) |

Fyne's current official documentation says Go 1.16+ for `gomobile`, but that is a separate Go Mobile workflow. Canon D should use Fyne packaging here and retain its stricter **Go 1.25.0+** requirement. [Go Mobile](https://go.dev/wiki/Mobile) is not an additional Fyne `fyne package` dependency established by the primary Fyne packaging source.

### Version decisions left to Canon D

- **Numeric NDK version:** Fyne does not pin a numeric NDK release in the inspected documentation; the code checks layout and compilers and rejects r19c-era toolchains. No current numeric NDK release is invented here. Pin a known-good NDK in Canon D CI and verify it on a real device. [NDK environment source](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/env.go)
- **Application ID:** `com.github.yzh44yzh.canon_d` is a reasonable stable reverse-DNS choice based on Canon D's GitHub remote, but namespace ownership is not verified by Fyne. The ID must contain a dot, not contain `-`, have no empty package segment, and segments cannot start with `_` or a digit. Once released, changing it creates a different Android application. [application-ID validation](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/commands/package.go)
- **minSdkVersion:** API **21** is what Fyne CLI v1.7.4 emits. Android's general documentation recommends explicitly declaring an accurate value, but claiming support below 21 requires inspection and testing of Fyne, not merely changing a number. A real API-21 device/emulator is the useful compatibility check. [Android `uses-sdk`](https://developer.android.com/guide/topics/manifest/uses-sdk-element); [Fyne build source](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/binres/sdk.go)
- **targetSdkVersion:** this Fyne CLI emits **29 for non-release/debug** packaging. Do not infer target SDK 36 from `androidTargetRelease=36`; that value applies to release mode, while Canon D is requesting debug APKs. Production-release hardening is out of scope. [build constants](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/build.go)

## Application metadata and assets

### Required app icon

A PNG icon must exist. With `-icon` unspecified, Fyne uses `<source-dir>/Icon.png`; a missing icon fails packaging. Fyne can convert another supported image to PNG, but Android app packaging expects a PNG. No `AndroidManifest.xml` is required. Fyne generates launcher resources and `android:icon` when an icon is available. [Fyne `validate()`](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/commands/package.go); [Fyne icon bundling](https://docs.fyne.io/extend/bundle/); [Android icon contract](https://developer.android.com/guide/topics/resources/drawable-resource)

No separate adaptive-icon assets are required. This CLI additionally recognizes `Icon-foreground.png`, `Icon-background.png`, and `Icon-monochrome.png`, or adaptive-icon metadata, but those are optional. No icon pixel-size requirement beyond a valid PNG is stated by Fyne; choose a square, high-resolution source and inspect the generated launcher resources. [adaptive-icon handling](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go)

### `FyneApp.toml`

The file belongs in the directory where `fyne package` runs, normally `cmd/mobile`. Fyne reads it if present; it can replace repeated CLI flags. Recommended minimal content for the initial app:

```toml
Website = "https://github.com/yzh44yzh/canon_d"

[Details]
Name = "CanonD"
ID = "com.github.yzh44yzh.canon_d"
Version = "1.0.0"
Build = 1
Icon = "Icon.png"

[Migrations]
fyneDo = true
```

- `Name`, `ID`, `Version`, `Build`, and `Icon` are the Fyne `Details` keys. `Build` is a positive integer and corresponds to Android `versionCode`; `Version` corresponds to `versionName`. CLI flags can override file values. [Fyne app metadata](https://docs.fyne.io/started/metadata/); [generated Android manifest](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/templates/data/AndroidManifest.xml)
- **Fyne v2.8 + CLI v1.7.4:** `fyneDo=true` opts into the single-goroutine Fyne event model. Fyne v2.9 makes this model the default. Fyne events/callbacks share a goroutine; calls from application-created goroutines must use `fyne.Do`, and `fyne.DoAndWait` is available when completion must be synchronous. This is relevant to mobile lifecycle saving and UI updates, but lifecycle integration belongs to the separate lifecycle decision. [Fyne goroutines](https://docs.fyne.io/started/goroutines/); [`fyne.Lifecycle`](https://docs.fyne.io/api/v2/fyne/lifecycle/)
- Fyne's `app.New()` is sufficient when `FyneApp.toml` is loaded by the packager. `app.NewWithID(...)` explicitly sets an app ID in Go; keep it identical to `Details.ID`, especially if mobile app-private storage depends on it. This packaging research does not require Fyne storage. [Fyne app package](https://docs.fyne.io/api/v2/app/pkg/)
- On a successful `fyne package`, this CLI version loads `FyneApp.toml`, increments `Details.Build`, and saves it. A project that wants reproducible manual versioning should verify and manage that mutation explicitly. [package completion source](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/commands/package.go)

### Embedded Markdown

No Fyne asset directory or `fyne bundle` command is required for embedded decks. Keep each original `.md` file unchanged, embed it in the Go mobile package, and pass `fyne.NewStaticResource(filename, bytes)` or call `parser.MakeDeck(filename, string(bytes))` after reading the `embed.FS`. The parser already accepts in-memory strings. Preserve exact original Markdown bytes and the exact base filename, because both are used by deck parsing/rendering and filename-based card identity.

```go
//go:embed decks/*.md
var deckFS embed.FS
```

`go:embed` cannot embed files whose names begin with `.` or `_`, and cannot embed files outside the package directory; keep the mobile package's embedded deck source layout accordingly. Embedded source files do not need to be copied into a separate Fyne `assets/` directory for this Go-only deck path. [Go `embed` package](https://pkg.go.dev/embed); [Canon D parser](../../internal/parser/parser.go); [Fyne resources](https://docs.fyne.io/extend/bundle/); [Fyne `StaticResource`](https://docs.fyne.io/api/v2/fyne/staticresource/)

Fyne's own `fyne bundle` generates Go `fyne.StaticResource` values using `go:embed`; it is useful for images/fonts, but is not mandatory for Canon D's Markdown decks. Fyne can also package `assets/`, but relying on extracted filesystem assets would reintroduce unnecessary runtime filesystem requirements. [Fyne resource bundling](https://docs.fyne.io/extend/bundle/); [Fyne mobile asset handling](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go)

### Android manifest and signing facts

Fyne generates `AndroidManifest.xml`; a custom manifest is optional for advanced needs. The generated app uses `org.golang.app.GoNativeActivity`, `android.app.lib_name`, launcher intents, `android:debuggable=true` for non-release builds, `android:versionCode`, and `android:versionName`. Manual `<uses-sdk>` is rejected because Fyne synthesizes it. [Fyne manifest template](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/templates/data/AndroidManifest.xml); [manifest processing](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/binres/binres.go)

The generated template currently includes `READ_EXTERNAL_STORAGE`, `WRITE_EXTERNAL_STORAGE`, and `INTERNET`. Canon D's embedded-deck/private-statistics design does not need external-storage access; permissions should be removed with a custom manifest if production needs minimal permissions. That is not a requirement for the initial debug APK research. [Fyne manifest template](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/templates/data/AndroidManifest.xml)

Android requires digitally signed APKs. Fyne's non-release writer emits its v1 JAR signature files (`META-INF/MANIFEST.MF`, `CERT.SF`, `CERT.RSA`) using an embedded development certificate that is intentionally stable across builds. This is suitable only for installing debug builds; it is not a production keystore. `fyne package` without `-release` does not ask for production signing. `fyne release` is a separate path. [Android signing](https://developer.android.com/studio/publish/app-signing); [Fyne APK writer](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/writer.go); [Fyne debug certificate](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go); [Fyne release docs](https://docs.fyne.io/started/distribution/)

## Exact debug-APK workflow

Assume a Go main package at `cmd/mobile`, `CanonD` app name, stable `com.github.yzh44yzh.canon_d` ID, `Icon.png`, and original decks embedded in Go. `ANDROID_HOME`, `ANDROID_NDK_HOME`, and `$ANDROID_HOME/platform-tools` are environment/toolchain prerequisites, not app secrets.

```sh
# Go binary installed to GOBIN or GOPATH/bin, which must also be in PATH.
go install fyne.io/tools/cmd/fyne@v1.7.4

# Run from the repository root.
go mod tidy

# Run from the directory containing go code, FyneApp.toml, Icon.png,
# and the go:embed directive.
cd cmd/mobile

CGO_ENABLED=1 \
PATH="$ANDROID_HOME/platform-tools:$PATH" \
fyne package \
  -os android \
  -app-id com.github.yzh44yzh.canon_d \
  -name CanonD \
  -icon Icon.png \
  -app-version 1.0.0 \
  -app-build 1 \
  -verbose
```

The application module must require `fyne.io/fyne/v2` v2.8.2, and `fyne` must resolve to CLI v1.7.4. Use `go mod edit -require=fyne.io/fyne/v2@v2.8.2` followed by `go mod tidy` if not already pinned. The CLI command is **`fyne package`**, with **`-os`**, **`-app-id`**, **`-name`**, **`-icon`**, **`-app-version`**, **`-app-build`**, and **`-verbose`**. There is no need to use `fyne build` first. Do not use `--src` for mobile; this CLI rejects mobile source-dir mode and must run inside the main package. Do not pass `-release`. [package flags](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/commands/package.go); [CLI migration](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.2/cmd/fyne/main.go); [official command example](https://docs.fyne.io/started/mobile/)

For `FyneApp.toml`-based builds, omit `-app-version`, `-app-build`, `-app-id`, `-name`, and `-icon`; retain `-os`, `-app-id` if explicitness is preferred, and `-verbose`. Note that successful packaging increments `Details.Build`, as described above.

Output naming comes from Fyne's `androidPkgName(appName)`, not the Go module name. With `CanonD`, output is **`CanonD.apk`**. [output-name source](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go)

A smaller arm64-only smoke-test command changes only targeting:

```sh
CGO_ENABLED=1 fyne package -os android/arm64 -app-id com.github.yzh44yzh.canon_d -name CanonD -icon Icon.png
```

Then connect a device with USB debugging or start an emulator and install:

```sh
adb install -r CanonD.apk
```

The official page uses `adb install myapp.apk`; `-r` additionally permits reinstalling an update signed with the same development key. Verify `adb devices`, generated manifest metadata, ABIs, signature, app launch, embedded decks, and persistence on a real device before calling the package ready. `apksigner verify --verbose CanonD.apk` is a useful Android SDK Build Tools check. No such device build is claimed here. [official `adb install`](https://developer.android.com/tools/adb#install); [Android `apksigner`](https://developer.android.com/studio/command-line/apksigner); [Fyne Mobile Packaging](https://docs.fyne.io/started/mobile/)

## Constraints and gotchas that affect implementation

1. **CGo and NDK are mandatory.** Fyne mobile packaging cannot use the host-only desktop C compiler. The CLI finds NDK LLVM tools per ABI and forces `CGO_ENABLED=1`, `GOOS=android`, matching `GOARCH`, `CC`, `CXX`, and `GOARM=7` for 32-bit ARM. Do not override those with host `CC`/`CXX` in the command environment; a durable project `CC` setting could interfere outside the CLI-controlled `go` invocation. [mobile environment source](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/env.go); [CGo requirement](https://docs.fyne.io/started/cross-compiling/)
2. **Correct env var names matter.** `ANDROID_HOME` and `ANDROID_NDK_HOME` are read. `ANDROID_SDK_ROOT` and `NDK_PATH` are not alternatives for actual Fyne packaging. `$ANDROID_HOME/ndk-bundle` is checked if `ANDROID_NDK_HOME` is absent/invalid; side-by-side NDKs otherwise need explicit configuration. `adb` must be in `PATH`. [Fyne NDK lookup](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/env.go); [Fyne Android SDK lookup](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/binres/sdk.go); [official requirements](https://docs.fyne.io/started/mobile/)
3. **No Fyne command mobile source-dir mode.** Run `fyne package` in `cmd/mobile`. `-src` is rejected for Android/iOS. `-executable` is ignored for mobile builds. [package validation](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/commands/package.go)
4. **App identity is stable.** `-app-id` is required for Android in this CLI, must be globally unique reverse-DNS-style, cannot contain `-`, cannot have empty segments, and package parts cannot begin `_` or digits. `ID` in `FyneApp.toml` is equivalent metadata; `app.NewWithID` should not contradict it. Android compares signing certificates for upgrades, so app ID and debug signing key must remain stable. [application-ID validation](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/commands/package.go); [Android signing identity](https://developer.android.com/studio/publish/app-signing)
5. **Version handling differs by input.** CLI defaults are version `0.0.1`, build `1`; integer and `x.y.z` forms are accepted. Android `versionCode` must increase for upgrades. With `FyneApp.toml`, this CLI increments `Details.Build` after successful packaging. Choose one source of truth in Canon D. [package defaults](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/commands/package.go); [Android versioning](https://developer.android.com/tools/publishing/versioning); [package completion source](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/commands/package.go)
6. **Debug does not mean unsigned.** Fyne signs the APK, uses a non-production debug key, and emits `android:debuggable=true`. Android production packages need a private, backed-up signing key, consistent certificate across updates, and any required release scheme; this initial effort intentionally stops before production signing. [Fyne writer](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/writer.go); [Android signing](https://developer.android.com/studio/publish/app-signing)
7. **`fyneDo` is required on v2.8.** Fyne says v2.8 requires opt-in with `[Migrations] fyneDo=true`; v2.9 makes it default. Do not mix direct Fyne calls from arbitrary goroutines with lifecycle callbacks. Use `fyne.Do`/`fyne.DoAndWait` where calls originate off the Fyne event goroutine. [Fyne goroutines](https://docs.fyne.io/started/goroutines/)
8. **Packaging is not proof of app runtime compatibility.** A successful APK only proves the toolchain and resources completed. No numeric NDK pin, API-21 real-device run, soft-keyboard check, lifecycle persistence check, or production store-signing check is established by these sources. These should be acceptance checks in Canon D. No claim is made that this note itself built or installed an APK.

## Source index

### First-party Fyne

- [Fyne Mobile Packaging](https://docs.fyne.io/started/mobile/)
- [Fyne App Metadata](https://docs.fyne.io/started/metadata/)
- [Fyne Bundling Resources](https://docs.fyne.io/extend/bundle/)
- [Fyne Using Goroutines](https://docs.fyne.io/started/goroutines/)
- [Fyne Cross Compiling](https://docs.fyne.io/started/cross-compiling/)
- [Fyne Distributing to App Stores](https://docs.fyne.io/started/distribution/)
- [Fyne v2.8.2 release](https://github.com/fyne-io/fyne/releases/tag/v2.8.2)
- [Fyne tools v1.7.4 release](https://github.com/fyne-io/tools/releases/tag/v1.7.4)
- Fyne tools v1.7.4 source: `cmd/fyne/internal/commands/package.go`, `package-mobile.go`; `cmd/fyne/internal/mobile/env.go`, `build.go`, `build_androidapp.go`, `writer.go`, `manifest.go`, `binres/binres.go`, `binres/sdk.go`; `cmd/fyne/internal/templates/data/AndroidManifest.xml`
- [Fyne `StaticResource`](https://docs.fyne.io/api/v2/fyne/staticresource/)

### First-party Go and Android

- [Go `embed`](https://pkg.go.dev/embed)
- [Go Mobile](https://go.dev/wiki/Mobile) — separate from Fyne packaging
- [Android `uses-sdk` / API levels](https://developer.android.com/guide/topics/manifest/uses-sdk-element)
- [Android app signing](https://developer.android.com/studio/publish/app-signing)
- [Android versioning](https://developer.android.com/tools/publishing/versioning)
- [Android ADB](https://developer.android.com/tools/adb)
- [Android `apksigner`](https://developer.android.com/studio/command-line/apksigner)
- [Android NDK](https://developer.android.com/ndk/downloads)
