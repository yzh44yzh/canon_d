# Research: Fyne Android build and packaging constraints

Ticket: [yzh44yzh/canon_d#6](https://github.com/yzh44yzh/canon_d/issues/6).
Parent map: [#1, Wayfinder: Initial Fyne Android app](https://github.com/yzh44yzh/canon_d/issues/1).

Research date: 2026-11-27. Source-access date for every external source below: 2026-11-27.
All external URLs cited in this document returned a successful HTTP status when checked on that date.

Scope: a **developer-built, directly installable debug APK** for Canon D, built from this
repository, which currently declares `go 1.25.0` and has no Fyne dependency
([local go.mod:1-12](../../go.mod#L1)). The eight decks are the Markdown files under
`decks/` listed in the parent issue
([issue #1 body](https://github.com/yzh44yzh/canon_d/issues/1)).

Verified facts and recommendations are separated deliberately. "Verified" means the claim is
quoted or directly derivable from a first-party Fyne, Go, or Android source. "Recommendation"
means a concrete choice inferred for Canon D, not a claim that upstream requires it.

## Answer at a glance

| Decision area | Concrete answer | Status |
|---|---|---|
| Fyne library | `fyne.io/fyne/v2 v2.8.3`; Go proxy `@latest` returned `v2.8.3` at `2026-11-20T13:27:29Z`, and the changelog heading is `## 2.8.3 - 24 Nov 2026` ([proxy](https://proxy.golang.org/fyne.io/fyne/v2/@latest), [changelog](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/CHANGELOG.md#L6), [tag](https://github.com/fyne-io/fyne/releases/tag/v2.8.3)). | Verified |
| Fyne CLI tool | `fyne.io/tools v1.7.4`; Go proxy `@latest` returned `v1.7.4` at `2026-10-19T09:49:46Z`; changelog heading `## 1.7.4 - 19 Oct 2026` ([proxy](https://proxy.golang.org/fyne.io/tools/@latest), [changelog](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/CHANGELOG.md#L6), [release](https://github.com/fyne-io/tools/releases/tag/v1.7.4)). | Verified |
| Go version | Keep `go 1.25.0`. Fyne `v2.8.3` and tools `v1.7.4` both declare `go 1.22.0` ([Fyne go.mod:3](https://github.com/fyne-io/fyne/blob/v2.8.3/go.mod#L3), [tools go.mod:3](https://github.com/fyne-io/tools/blob/v1.7.4/go.mod#L3)); Go requires a module's `go` line to be at least that of its dependencies ([Go modules reference](https://go.dev/doc/modules/gomod-ref#go)). | Verified compatibility from declared minimums; no first-party compatibility matrix exists |
| CGo and ABIs | The CLI forces `CGO_ENABLED=1`, `GOOS=android`, per-ABI `GOARCH`, NDK `CC`/`CXX`, and `GOARM=7` for 32-bit ARM; `android` builds `arm`, `arm64`, `386`, `amd64` (mapped to `armeabi-v7a`, `arm64-v8a`, `x86`, `x86_64`) ([Fyne cross-compiling](https://docs.fyne.io/started/cross-compiling/), [tools env.go](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/env.go), [shared-library build](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go)). | Verified |
| Android minimum | Manifest gets `minSdkVersion="21"` = Android 5.0 ([tools binres/sdk.go:14-15](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/binres/sdk.go#L14-L15), [Android API table](https://developer.android.com/guide/topics/manifest/uses-sdk-element)). | Verified |
| Debug target | Tools set `androidTargetDebug = 29` (Android 10) and use it when `!buildRelease`; release target is `36` ([tools build.go:58-60](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build.go#L58-L60), [build.go:122-127](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build.go#L122-L127), [API table](https://developer.android.com/guide/topics/manifest/uses-sdk-element)). | Verified |
| Target-floor install rules | Android 14 blocks `targetSdkVersion` below 23 and Android 15 blocks below 24; ADB has an explicit bypass only for testing older targets ([Android 14](https://developer.android.com/about/versions/14/behavior-changes-all#minimum-target-api-level), [Android 15](https://developer.android.com/about/versions/15/behavior-changes-all#minimum-target-api-level)). Fyne debug target 29 clears both floors. | Verified |
| Debug signing | `fyne package -os android` uses a PEM RSA key embedded in the tools source, with the comment `Must be consistent across builds so newer app versions can be installed` ([build_androidapp.go:397-410](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go#L397-L410), [build_androidapp.go:580-610](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go#L580-L610)). | Verified |
| Install | `adb install path_to_apk` is the documented install command; `pm install -r` means `Reinstall an existing app, keeping its data` ([Android adb docs](https://developer.android.com/tools/adb#install)). | Verified |
| Deck data | Use Go `//go:embed` with `embed.FS` from the package that owns the Markdown files; `..` patterns are forbidden ([Go embed source](https://github.com/golang/go/blob/master/src/embed/embed.go#L5-L99)). | Verified rule; Canon D placement is a recommendation |
| Icon | `Icon.png` in the directory where `fyne` runs is the default; a missing file produces `Missing application icon at "..."` ([tools package.go:328-332](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package.go#L328-L332)); Fyne docs show `-icon mobileIcon.png` ([Fyne mobile docs](https://docs.fyne.io/started/mobile/)). | Verified; exact pixel size is unresolved |

## Verified facts

### Versions and Go compatibility

- The Go module proxy is the first-party module index: it reports `v2.8.3` as the latest
  `fyne.io/fyne/v2` release at `2026-11-20T13:27:29Z`
  ([proxy response](https://proxy.golang.org/fyne.io/fyne/v2/@latest)). The Fyne changelog
  writes `## 2.8.3 - 24 Nov 2026`
  ([changelog:6](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/CHANGELOG.md#L6)).
  Note a source discrepancy: the GitHub Releases API still surfaced `v2.8.2` as the first
  result while the tag, changelog, and module proxy all expose `v2.8.3`
  ([Fyne releases](https://github.com/fyne-io/fyne/releases),
  [v2.8.3 tag](https://github.com/fyne-io/fyne/releases/tag/v2.8.3)).
- The Fyne CLI tool's latest module release is `fyne.io/tools v1.7.4`
  ([proxy response](https://proxy.golang.org/fyne.io/tools/@latest)); its changelog says
  `## 1.7.4 - 19 Oct 2026`
  ([tools changelog:6](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/CHANGELOG.md#L6)).
  The tools module itself requires `fyne.io/fyne/v2 v2.8.1`
  ([tools go.mod:6](https://github.com/fyne-io/tools/blob/v1.7.4/go.mod#L6)), which is an
  older minimum than the current library `v2.8.3`; no first-party note prohibits using the
  newer library.
- Both modules declare `go 1.22.0`
  ([Fyne go.mod:3](https://github.com/fyne-io/fyne/blob/v2.8.3/go.mod#L3),
  [tools go.mod:3](https://github.com/fyne-io/tools/blob/v1.7.4/go.mod#L3)). Go's module
  reference states: `The go directive sets the minimum version of Go required to use this
  module` and `The go line must be greater than or equal to the go line of all dependencies`
  ([Go modules reference](https://go.dev/doc/modules/gomod-ref#go)). The repository's
  `go 1.25.0` therefore satisfies both dependencies. Go's compatibility promise says
  programs written to Go 1 continue to compile and run over the specification's lifetime
  ([Go 1 compatibility](https://go.dev/doc/go1compat)).
- Fyne's changelog records the earlier fix `Resolve compile issue with Go 1.25.0 caused by
  golang.org/x/tools conflict`
  ([Fyne changelog](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/CHANGELOG.md),
  [tools changelog](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/CHANGELOG.md)).
  Go 1.25 release notes introduce no Android port change
  ([Go 1.25 release notes, "Ports"](https://go.dev/doc/go1.25#ports)).
- Fyne's Quick Start says to install with `go get fyne.io/fyne/v2@latest` and
  `go install fyne.io/tools/cmd/fyne@latest`
  ([Fyne Quick Start](https://docs.fyne.io/started/quick/)).

### Android SDK, NDK, build tools, and environment

- Fyne's mobile doc states: `For Android builds you must have the Android SDK and NDK
  installed with appropriate environment set up so that the tools (such as adb) can be
  found on the command line`
  ([Fyne mobile packaging](https://docs.fyne.io/started/mobile/)).
- Fyne's Quick Start recommends Android Studio's SDK Manager with the `NDK (Side by side)`
  package, or a standalone NDK with `ANDROID_NDK_HOME` set
  ([Fyne Quick Start, Android tab](https://docs.fyne.io/started/quick/)).
- Fyne tools look for the NDK in `$ANDROID_NDK_HOME` first, then `$ANDROID_HOME/ndk-bundle`,
  and error with `no Android NDK found in $ANDROID_HOME/ndk-bundle nor in $ANDROID_NDK_HOME`
  ([tools env.go:253-270](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/env.go#L253-L270)).
- The exact first-party NDK floor is `Make sure your NDK version is >= r19c`; Fyne also notes
  `we require r19c, not just r19b`
  ([tools env.go:166-173](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/env.go#L166-L173)).
  No exact newer NDK version is pinned by Fyne.
- CGo is part of the Fyne mobile contract: the CLI sets `CGO_ENABLED=1`, `GOOS=android`, the
  matching `GOARCH`, NDK-provided `CC`/`CXX`, and `GOARM=7` for 32-bit ARM
  ([Fyne cross-compiling](https://docs.fyne.io/started/cross-compiling/),
  [tools env.go](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/env.go),
  [shared-library build](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go)).
  Do not set a host `CC`/`CXX` in the project environment; the CLI supplies per-ABI compilers.
  The `android` target builds `arm`, `arm64`, `386`, and `amd64` (`armeabi-v7a`, `arm64-v8a`,
  `x86`, `x86_64`); `-os android/arm64` is a supported smaller smoke-test target
  ([tools env.go](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/env.go)).
- `ANDROID_HOME` is required: source errors include `could not find android tools, missing
  ANDROID_HOME` and `ANDROID_HOME env var not set`
  ([tools util/mobile.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/util/mobile.go),
  [tools binres/sdk.go:17-19](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/binres/sdk.go#L17-L19)).
  The build-tools directory must contain `aapt2`
  ([tools util/mobile.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/util/mobile.go)).
  No exact build-tools version is pinned by Fyne.
- Fyne selects the highest-numbered installed `$ANDROID_HOME/platforms/*/android.jar`
  ([tools binres/sdk.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/binres/sdk.go));
  it expects a generated toolchain under `$GOPATH/pkg/gomobile` and reports
  `toolchain not installed, run \`gomobile init\`` when absent
  ([tools env.go:57](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/env.go#L57)).
- No first-party Fyne source reviewed states a required JDK version. The release path names
  `jarsigner`, `apksigner`, and `zipalign`, but a plain debug `fyne package` does not use the
  release signing path
  ([tools release.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/release.go)).

### Minimum Android version, target SDK, and manifest

- Fyne tools hard-code `minAndroidAPI = 21`
  ([tools build.go:35-37](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build.go#L35-L37))
  and `const MinSDK = 21`
  ([tools binres/sdk.go:14-15](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/binres/sdk.go#L14-L15)).
  The generated manifest synthesizes `<uses-sdk android:minSdkVersion="...">` from `MinSDK`
  and `android:targetSdkVersion` from the selected target
  ([tools binres/binres.go:247-248](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/binres/binres.go#L247-L248),
  [binres.go:269-290](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/binres/binres.go#L269-L290)).
  A hand-written `<uses-sdk>` in a custom manifest is rejected with
  `manual declaration of uses-sdk in AndroidManifest.xml not supported`
  ([binres.go:247-248](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/binres/binres.go#L247-L248)).
- Android's API table maps API 21 to Android 5.0 and API 29 to Android 10
  ([Android `<uses-sdk>` documentation](https://developer.android.com/guide/topics/manifest/uses-sdk-element)).
  Android states that the system blocks installation when `minSdkVersion` exceeds the
  device's API level
  ([Android `<uses-sdk>` documentation](https://developer.android.com/guide/topics/manifest/uses-sdk-element)).
- Fyne tools choose `androidTargetRelease = 36` and `androidTargetDebug = 29`, with the
  comment `TODO once we have gomobile debug signing working for v2 android signs`
  ([tools build.go:58-60](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build.go#L58-L60),
  [build.go:122-127](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build.go#L122-L127)).
  Plain `fyne package -os android` does not set `release`, so it uses the debug target
  ([tools package-mobile.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package-mobile.go)).
- Android 11 requires apps targeting API 30+ to use signature scheme v2 or higher:
  `Users can't install or update apps that are only signed with APK Signature Scheme v1 on
  devices that run Android 11`
  ([Android 11 behavior changes](https://developer.android.com/about/versions/11/behavior-changes-11#minimum-signature-scheme)).
  Fyne's debug target 29 sits below that rule, consistent with the source TODO above.
- Android 14 says `apps with a targetSdkVersion lower than 23 can't be installed`, and
  Android 15 raises that to `lower than 24`
  ([Android 14](https://developer.android.com/about/versions/14/behavior-changes-all#minimum-target-api-level),
  [Android 15](https://developer.android.com/about/versions/15/behavior-changes-all#minimum-target-api-level)).
  Fyne's debug target 29 clears both install floors.
- Conflict to note: Fyne's wiki still says `Android (4.1+)`
  ([Fyne Supported Platforms wiki](https://github.com/fyne-io/fyne/wiki/Supported-Platforms)),
  while the current tools source sets `minSdkVersion=21` (Android 5.0). The tools PR that
  raised the value says it was `updating minimum to 21 (Android 5) as a step towards
  matching upstream support constraints`
  ([tools PR #151](https://github.com/fyne-io/tools/pull/151)). The manifest value is the
  effective install constraint for a build made with tools `v1.7.4`.

### Debug signing, manifest, and direct install

- The debug APK writer decodes the embedded PEM key with `pem.Decode([]byte(debugCert))` and
  creates the APK writer with that private key
  ([tools build_androidapp.go:397-410](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go#L397-L410)).
  The source comment above the key says `A random uninteresting private key. Must be
  consistent across builds so newer app versions can be installed`
  ([build_androidapp.go:580-610](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go#L580-L610)).
  This is a compiled-in key, not `$HOME/.android/debug.keystore`; the writer source itself
  contrasts Android Studio's JKS debug keystore with the generated key used for normal builds
  ([tools writer.go:51-65](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/writer.go#L51-L65)).
- Fyne's APK writer adds `META-INF/MANIFEST.MF`, `META-INF/CERT.SF`, and `META-INF/CERT.RSA`
  using SHA-1 digests
  ([tools writer.go:7-45](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/writer.go#L7-L45)).
  That is the v1/JAR signing layout; the debug writer adds no v2/v3 block. Confirm the final
  scheme with `apksigner verify --verbose` on the generated APK before relying on it.
- Open first-party issue: tools [#154](https://github.com/fyne-io/tools/issues/154) tracks
  deterministic debug signing and debug/release API parity, so debug signing behavior can
  change in a later tools release.
- Fyne's generated manifest sets `package="{{.JavaPkgPath}}"`,
  `android:versionCode="{{.Build}}"`, `android:versionName="{{.Version}}"`,
  `android:label="{{.Name}}"`, `android:debuggable="{{.Debug}}"`, and activity
  `org.golang.app.GoNativeActivity`
  ([tools AndroidManifest.xml:1-19](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/templates/data/AndroidManifest.xml#L1-L19)).
  The template also declares `WRITE_EXTERNAL_STORAGE`, `READ_EXTERNAL_STORAGE`, and
  `INTERNET`
  ([AndroidManifest.xml:25-27](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/templates/data/AndroidManifest.xml#L25-L27)).
- Android's adb documentation gives `adb install path_to_apk`
  ([adb docs, "Install an app"](https://developer.android.com/tools/adb#install)) and
  documents `pm install -r` as `Reinstall an existing app, keeping its data`
  ([adb docs, `pm` table](https://developer.android.com/tools/adb#pm)). It also documents
  `am start [options] intent` and the explicit component form `-n component`
  ([adb docs, `am` table](https://developer.android.com/tools/adb#am)).
  With Fyne's default activity, the component is
  `APP_ID/org.golang.app.GoNativeActivity`; a custom manifest can change it.
- Verify the generated APK's signature scheme with
  `apksigner verify --verbose Canon_D.apk`; `apksigner` is an Android SDK build-tools command
  ([Android `apksigner` docs](https://developer.android.com/studio/command-line/apksigner)).
- ADB requires Developer options and USB debugging; on Android 4.2.2+ the device shows an RSA
  authorization dialog, and `adb devices` confirms the connection
  ([adb docs](https://developer.android.com/tools/adb#enable)). Android 11+ also supports
  wireless debugging with pairing
  ([adb docs, "Connect to a device over Wi-Fi"](https://developer.android.com/tools/adb#wireless-adb-android-11)).
- Direct installation from a website or email is different from ADB: users must opt into
  `Install unknown apps` on Android 8.0+ (API 26+)
  ([Android alternative distribution](https://developer.android.com/distribute/marketing-tools/alternative-distribution#unknown-sources)).
- Android developer verification is explicitly not an ADB blocker: `As a developer, you are
  free to install apps without verification with ADB` and `The waiting period does not apply
  to ADB installs`
  ([developer verification FAQ](https://developer.android.com/developer-verification/guides/faq)).
  The same FAQ says September 30, 2026 enforcement is limited to listed stores, while a
  global rollout begins in 2027
  ([developer verification FAQ](https://developer.android.com/developer-verification/guides/faq)).
- Fyne's own `fyne install --os android` reuses an existing target APK if present before
  running `adb install`; verbose output distinguishes `Using existing package` from
  `Packaging ... (no existing package)`
  ([tools install.go:299-334](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/install.go#L299-L334)).
  This can install a stale APK after source changes; tools issue
  [#168](https://github.com/fyne-io/tools/issues/168) was closed by the reuse change.

### Exact commands and flags

Documented/install commands:

```sh
go install fyne.io/tools/cmd/fyne@v1.7.4
```

The `v1.7.4` pin comes from the Go proxy
([proxy](https://proxy.golang.org/fyne.io/tools/@latest)); Fyne's docs use `@latest`
([Quick Start](https://docs.fyne.io/started/quick/)). When implementation adds the library,
use `go get fyne.io/fyne/v2@v2.8.3`, matching the verified latest library release
([proxy](https://proxy.golang.org/fyne.io/fyne/v2/@latest)).

Environment shape (paths are examples, not first-party required values):

```sh
export ANDROID_HOME=/path/to/Android/Sdk
export ANDROID_NDK_HOME=/path/to/android-ndk
export PATH="$PATH:$ANDROID_HOME/platform-tools:$HOME/go/bin"
```

`ANDROID_HOME`/`ANDROID_NDK_HOME` behavior is verified in
[tools env.go:253-270](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/env.go#L253-L270)
and [tools util/mobile.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/util/mobile.go);
`adb` on `PATH` is required by [Fyne mobile docs](https://docs.fyne.io/started/mobile/).
If Fyne reports `toolchain not installed, run \`gomobile init\``, initialize the gomobile
toolchain as the error states
([tools env.go:57](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/env.go#L57)).

Documented Android package command (from [Fyne mobile docs](https://docs.fyne.io/started/mobile/)):

```sh
fyne package -os android -app-id com.example.myapp -icon mobileIcon.png
```

Canon D debug command (recommendation; see app-ID decision below):

```sh
cd /path/to/canon-d-fyne-android-build/cmd/canond
fyne package -os android -app-id com.yzh44yzh.canond -name "Canon D" -icon Icon.png
adb devices -l
adb install -r Canon_D.apk
adb shell am start -n com.yzh44yzh.canond/org.golang.app.GoNativeActivity
```

Smaller arm64-only smoke-test variant (recommendation; output is still `Canon_D.apk`):

```sh
CGO_ENABLED=1 fyne package -os android/arm64 -app-id com.yzh44yzh.canond -name "Canon D" -icon Icon.png
```

`-os android/arm64` is supported by the CLI's per-architecture build path
([tools env.go](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/mobile/env.go));
the output name still comes from `androidPkgName`
([build_androidapp.go:474-492](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go#L474-L492)).

Basis for the exact flags and behavior:

- `-app-id` is canonical with aliases `id` and `appID`; its usage says `set app-id in reversed
  domain notation for android, darwin, and windows targets`
  ([tools flags.go:52-59](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/flags.go#L52-L59)).
  Android app IDs must contain a dot, must not contain `-`, and dot segments may not start
  with `_` or a digit
  ([tools package.go:418-451](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package.go#L418-L451)).
- `-icon` is a package flag; `Icon.png` in the source directory is the default and a missing
  icon is an error
  ([tools package.go:328-332](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package.go#L328-L332)).
- `-name` sets the application name; if absent, the tool derives it from the module path's
  last element
  ([tools package.go:313-327](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package.go#L313-L327),
  [calculateExeName:353-372](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package.go#L353-L372)).
- Output name is `androidPkgName(name) + ".apk"` for non-release Android builds
  ([tools build.go:44-51](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build.go#L44-L51)).
  `androidPkgName` replaces every non-alphanumeric rune with `_`
  ([build_androidapp.go:474-492](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go#L474-L492)),
  so `-name "Canon D"` produces `Canon_D.apk`.
- There is no `-debug` package flag; `--release` is described as `enable installation in
  release mode, disable debug, etc`, so omit it
  ([tools flags.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/flags.go),
  [tools package.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package.go)).
  `fyne release` sets distribution/release flags and requires Android release parameters such
  as a keystore
  ([tools release.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/release.go)).
- Mobile builds do not support `--source-dir` (`parameter --source-dir is currently not
  supported for mobile builds. Change directory to the main package and try again`) and ignore
  `--executable` (`Parameter --executable is ignored for mobile builds.`)
  ([tools package.go:275-323](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package.go#L275-L323)).
  Run `fyne package` from the main-package directory.
- `fyne install --os android` is an alternative that packages when needed and calls
  `adb install`
  ([tools install.go:299-334](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/install.go#L299-L334)),
  but the explicit package + `adb install -r` sequence avoids the stale-package risk.

### App metadata and required assets

- `FyneApp.toml` is optional; its `[Details]` table supports `Icon`, `Name`, `ID`, `Version`,
  and `Build`, and command-line parameters override those values
  ([Fyne App Metadata docs](https://docs.fyne.io/started/metadata/)). A minimal file for the
  Android main package (final ID/name are still decision questions below) is:

  ```toml
  Website = "https://github.com/yzh44yzh/canon_d"

  [Details]
  Name = "Canon D"
  ID = "com.yzh44yzh.canond"
  Version = "1.0.0"
  Build = 1
  Icon = "Icon.png"

  [Migrations]
  fyneDo = true
  ```

  The top-level `Website` key and `Details` keys are shown in Fyne's metadata docs, and
  `[Migrations] fyneDo = true` is the v2.8 opt-in described in the goroutine docs
  ([Fyne App Metadata docs](https://docs.fyne.io/started/metadata/),
  [Fyne goroutines](https://docs.fyne.io/started/goroutines/)).
- The manifest maps build/version/name directly: `android:versionCode="{{.Build}}"`,
  `android:versionName="{{.Version}}"`, `android:label="{{.Name}}"`
  ([tools AndroidManifest.xml:1-13](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/templates/data/AndroidManifest.xml#L1-L13)).
  Tool defaults are `defaultAppBuild = 1` and `defaultAppVersion = "0.0.1"`
  ([tools package.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package.go)).
  Android requires `versionCode` to increase for an upgrade
  ([Android versioning](https://developer.android.com/tools/publishing/versioning)).
  This CLI version increments and saves `Details.Build` after a successful packaging run when
  metadata is used, so choose one source of truth for versioning explicitly
  ([tools package.go](https://raw.githubusercontent.com/fyne-io/tools/v1.7.4/cmd/fyne/internal/commands/package.go)).
- If the Go code calls `app.NewWithID(...)`, keep its ID identical to `Details.ID`; the app
  package documents `NewWithID` as setting the application ID explicitly
  ([Fyne app package](https://docs.fyne.io/api/v2/app/pkg/)).
- The launcher icon comes from `-icon` or `Icon.png`; non-PNG icons are normalized to PNG
  ([tools package.go:328-340](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package.go#L328-L340),
  [normaliseIcon:392-415](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package.go#L392-L415)).
  Fyne's desktop packaging docs add: `All of these commands also support a default icon file
  of Icon.png`
  ([Fyne packaging docs](https://docs.fyne.io/started/packaging/)).
- Adaptive icon layers are detected from metadata or the conventions
  `Icon-foreground.png`, `Icon-background.png`, and `Icon-monochrome.png`
  ([tools build_androidapp.go:214-241](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go#L214-L241)).
  Generated resources include `mipmap-xxxhdpi/ic_launcher.png`, `ic_launcher_round.png`, and
  adaptive XML under `mipmap-anydpi-v26`
  ([tools android_resources.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/android_resources.go)).
  No first-party source reviewed gives a required launcher-icon pixel size; adaptive layers
  are optional.
- A custom `AndroidManifest.xml` in the app's main package is used if present; otherwise Fyne
  generates one from its template
  ([tools build_androidapp.go:67-84](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go#L67-L84)).
- An `assets/` subdirectory in the main package is copied into the APK; hidden files are
  skipped and `assets/icon.png` is treated specially rather than copied as a regular asset
  ([tools build_androidapp.go:243-290](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go#L243-L290)).
  The decks do **not** need to go there because `go:embed` puts the bytes inside the Go binary.

### Embedded Markdown decks

- Go's `embed` documentation states: `Go source files that import "embed" can use the
  //go:embed directive to initialize a variable of type string, []byte, or FS with the
  contents of files read from the package directory or subdirectories at compile time`
  ([Go embed source:5-9](https://github.com/golang/go/blob/master/src/embed/embed.go#L5-L9)).
  Go 1.16 introduced `//go:embed`
  ([Go 1.16 release notes](https://go.dev/doc/go1.16#embed)).
- Directive constraints are explicit: patterns are relative to the package directory;
  `Patterns may not contain '.' or '..' or empty path elements`; each pattern must match at
  least one file or non-empty directory
  ([Go embed source:64-99](https://github.com/golang/go/blob/master/src/embed/embed.go#L64-L99)).
  Therefore the existing root-level `decks/*.md` files cannot be embedded from `cmd/tui` with
  `../../decks/*.md`; a Go package must live beside the files, or the files must be
  copied/generated into the embedding package.
- `embed.FS` is `a read-only collection of files` and implements `fs.FS`, with `Open`,
  `ReadDir`, and `ReadFile`
  ([Go embed source, `FS` type](https://github.com/golang/go/blob/master/src/embed/embed.go),
  [pkg.go.dev](https://pkg.go.dev/embed#FS)).
- Fyne/gomobile Android builds invoke the normal Go toolchain: Fyne's mobile build calls
  `go build` with `-buildmode=c-shared`
  ([tools build_androidapp.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go#L103-L108),
  [tools build.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build.go#L320-L361)).
  A Fyne maintainer wrote that embed `works just as well for Fyne as any other Go application`
  ([Fyne issue #4702 comment](https://github.com/fyne-io/fyne/issues/4702#issuecomment-1980260352)).
  No first-party issue was found that reports `go:embed` itself failing in a Fyne Android
  build; [golang/go#61656](https://github.com/golang/go/issues/61656) contains `//go:embed`
  code but fails later with duplicate native symbols, not with embedding. This is
  mechanism-level evidence, not a published Fyne Android embed test.
- Fyne's own `fyne bundle` command generates `//go:embed` plus a `fyne.StaticResource`
  ([Fyne bundling docs](https://docs.fyne.io/extend/bundle/),
  [tools bundle.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/bundle.go)).
  A `fyne.Resource` is only `a name and byte array content`
  ([Fyne resource.go](https://github.com/fyne-io/fyne/blob/develop/resource.go#L10-L17)).
  The docs say `The Fyne toolkit prefers the use of "fyne bundle" as it has various benefits`
  ([Fyne bundling docs](https://docs.fyne.io/extend/bundle/)). Use it for UI resources that
  must be `fyne.Resource`; plain `embed.FS` is sufficient for Markdown source data.

### Constraints and bugs that affect implementation

- Fyne docs state: `Since Fyne v2.6.0 all Fyne events and callbacks happen on the same
  goroutine`; app-created goroutines must call `fyne.Do` when invoking Fyne APIs and
  `fyne.DoAndWait` when waiting for completion. In v2.8 the model is opt-in via
  `[Migrations] fyneDo = true`; it becomes default in v2.9.0
  ([Fyne Using Goroutines](https://docs.fyne.io/started/goroutines/)).
- `fyne.Lifecycle` exposes exactly `SetOnEnteredForeground`, `SetOnExitedForeground`,
  `SetOnStarted`, and `SetOnStopped`
  ([Fyne API docs](https://docs.fyne.io/api/v2/fyne/lifecycle/)). The first-party comments say
  `SetOnExitedForeground hooks into the app losing input focus and going into the background`
  and `SetOnStopped hooks into an event that says the app is no longer running`
  ([Fyne API docs](https://docs.fyne.io/api/v2/fyne/lifecycle/)). Android's own lifecycle docs
  warn that `onPause` is `the first indication that the user is leaving your activity` but
  `does not always mean the activity is being destroyed`
  ([Android activity lifecycle](https://developer.android.com/guide/components/activities/activity-lifecycle)).
- Fyne app-run docs require one runloop: `App.Run()` or `Window.ShowAndRun()` must be called
  at the end of `main()` setup, only one runloop should exist, and `App.Quit()` is not
  supported on mobile
  ([Fyne Application and RunLoop](https://docs.fyne.io/started/apprun/)).
- Android 15 says native-code apps need rebuilding for 16 KB-page devices
  ([Android 15 behavior changes](https://developer.android.com/about/versions/15/behavior-changes-all#support-for-16-kb-page-sizes)).
  Fyne APKs contain native `.so` files
  ([tools build_androidapp.go:424-442](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go#L424-L442)),
  but no first-party source was found that proves Go 1.25/Fyne `.so` alignment. Treat it as an
  unresolved device-test item, not a verified blocker.
- Fyne issue [#5233](https://github.com/fyne-io/fyne/issues/5233) is open: users report Play
  Protect warnings that the app was built for an older Android version; a maintainer says apps
  target the latest API but support very old versions and that minSDK is hard-coded because
  changing it can change compiled-in Android shims. Issue [#5872](https://github.com/fyne-io/fyne/issues/5872)
  was closed as a duplicate of the SDK-setting discussion; [#1127](https://github.com/fyne-io/fyne/issues/1127)
  explains `The minimum is what it will run on` and `The targetSDK is used to make APIs
  available...`.
- Resume/lock regressions were reported and fixed on `develop`: [#6408](https://github.com/fyne-io/fyne/issues/6408)
  (`EGL_BAD_SURFACE on resume from-Recents`, Android 9; fix announced for the next release) and
  [#6276](https://github.com/fyne-io/fyne/issues/6276) (Android lock/unlock freeze). Add
  background/resume regression tests.
- Input bugs: [#6587](https://github.com/fyne-io/fyne/issues/6587) (soft-keyboard backspace
  repeat deletes one character) is open with PR [#6680](https://github.com/fyne-io/fyne/pull/6680)
  open; hardware-keyboard duplication [#6271](https://github.com/fyne-io/fyne/issues/6271)
  was fixed on `develop`. Pointer/touch accuracy issues
  [#5256](https://github.com/fyne-io/fyne/issues/5256) and [#6585](https://github.com/fyne-io/fyne/issues/6585)
  are open.
- Storage bugs are open: [#3255](https://github.com/fyne-io/fyne/issues/3255) reports
  `storage.CanWrite() always return true?` and `storage.Writer() cannot create a file to
  write?` on Android; [#4490](https://github.com/fyne-io/fyne/issues/4490) reports a
  directory-creation error that did not occur from `main`; [#4499](https://github.com/fyne-io/fyne/issues/4499)
  reports unusable Android folder-dialog URIs. Test on-device persistence before relying on
  Fyne storage APIs.
- UI performance: [#1062](https://github.com/fyne-io/fyne/issues/1062) tracks Android text and
  scroll responsiveness; [#6683](https://github.com/fyne-io/fyne/issues/6683) shows PNG icons
  being decoded repeatedly in a long scroll list while SVG was cached. Keep deck/statistics
  lists lightweight.
- Go/gomobile version-specific reports: [#69109](https://github.com/golang/go/issues/69109)
  is open for x86_64 emulator `dlopen failed` TLS crashes, and
  [#77484](https://github.com/golang/go/issues/77484) is open for missing `.gopclntab` in
  gomobile `.so` files. Record Go/Fyne/tools versions in test notes and include a physical
  device plus an emulator in validation. Android's own guidance says
  `Always test your Android app on a real device before releasing it to users`
  ([Android on-device testing](https://developer.android.com/studio/run/device)); the emulator
  is intended for API-level and screen-size coverage
  ([Android emulator docs](https://developer.android.com/studio/run/emulator)).

## Recommendations for Canon D

These are choices for implementation, not upstream requirements.

1. Pin `fyne.io/fyne/v2 v2.8.3` and install `fyne.io/tools/cmd/fyne@v1.7.4`; keep Go 1.25.0
   ([versions](https://proxy.golang.org/fyne.io/fyne/v2/@latest),
   [tools](https://proxy.golang.org/fyne.io/tools/@latest),
   [local go.mod](../../go.mod#L3)).
2. Target a debug APK only: plain `fyne package -os android`, no `--release`/`fyne release`.
   Expect `minSdkVersion=21` and `targetSdkVersion=29`
   ([tools build.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build.go#L35-L60),
   [binres.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/binres/binres.go#L269-L290)).
3. Install SDK Platform Tools, SDK Platform 29, a build-tools version containing `aapt2`, and a
   current Android Studio NDK side-by-side package (Fyne's only stated floor is `>= r19c`).
   Add SDK Platform 36 only if later producing release/distribution builds
   ([Fyne Quick Start](https://docs.fyne.io/started/quick/),
   [tools env.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/env.go#L166-L173),
   [tools build.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build.go#L58-L60)).
4. Make a dedicated main package for the Android app (for example `cmd/canond/`) and run
   `fyne package` from that directory; mobile builds reject `--source-dir` and ignore
   `--executable`
   ([tools package.go:275-323](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package.go#L275-L323)).
   Keep the existing TUI entry point and domain packages unchanged.
5. Add `Icon.png` (square PNG) to that main package, plus an optional `FyneApp.toml` with
   `Icon`, `Name`, `ID`, `Version`, and `Build`
   ([Fyne metadata docs](https://docs.fyne.io/started/metadata/),
   [tools package.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package.go#L328-L340)).
   No first-party pixel size was found; 512×512 is a practical design choice, not a verified
   requirement.
6. Put the eight Markdown decks in a Go package that can legally embed them. Either add a Go
   file in `decks/` using `//go:embed *.md`, or copy/generate the eight files into the app's
   embed package with their filenames preserved (card identity includes the deck filename,
   per [GLOSSARY.md](../../GLOSSARY.md#L23-L25)). Do not use `../../decks/*.md`
   ([Go embed rules](https://github.com/golang/go/blob/master/src/embed/embed.go#L64-L99)).
7. Read decks with `embed.FS.ReadFile`; do not require `fyne bundle` for Markdown. Use
   `fyne bundle`/`fyne.StaticResource` only for images or fonts that Fyne APIs need as
   resources
   ([Fyne bundling docs](https://docs.fyne.io/extend/bundle/),
   [Fyne resource.go](https://github.com/fyne-io/fyne/blob/develop/resource.go#L10-L17)).
8. For background saves, register both `SetOnExitedForeground` and `SetOnStopped`: the former
   matches Android backgrounding, the latter is the terminal stop signal
   ([Fyne lifecycle docs](https://docs.fyne.io/api/v2/fyne/lifecycle/)). Persist at existing
   round checkpoints as well, because process death cannot be relied on to call either hook.
9. Enable `[Migrations] fyneDo = true` in v2.8 and route all app-created goroutine UI updates
   through `fyne.Do`/`fyne.DoAndWait`
   ([Fyne goroutine docs](https://docs.fyne.io/started/goroutines/)).
10. Install with `adb install -r Canon_D.apk` and launch with
    `adb shell am start -n APP_ID/org.golang.app.GoNativeActivity`; enable Developer options,
    USB or wireless debugging, and accept the RSA prompt
    ([adb docs](https://developer.android.com/tools/adb#install),
    [Fyne manifest](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/templates/data/AndroidManifest.xml#L9-L19)).
    Prefer explicit `fyne package` + `adb install` over `fyne install --os android` to avoid
    stale APK reuse
    ([tools install.go:299-334](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/install.go#L299-L334)).
11. Validate on a physical device and at least one emulator, with explicit resume-from-Recents,
    lock/unlock, soft-keyboard entry, long-list scrolling, and statistics-persistence checks
    ([Android device testing](https://developer.android.com/studio/run/device),
    [Fyne #6408](https://github.com/fyne-io/fyne/issues/6408),
    [#6276](https://github.com/fyne-io/fyne/issues/6276),
    [#6587](https://github.com/fyne-io/fyne/issues/6587),
    [#6683](https://github.com/fyne-io/fyne/issues/6683),
    [#3255](https://github.com/fyne-io/fyne/issues/3255)).

## Unresolved decision questions

These need a product/architecture decision or a build test; the research did not resolve them.

1. **Application ID.** Both `com.yzh44yzh.canond` and `com.github.yzh44yzh.canon_d` satisfy
   Fyne's validation: each has a dot, no `-`, no empty segment, and no segment starting with
   `_` or a digit (an underscore *inside* a segment is allowed)
   ([tools package.go:418-451](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/commands/package.go#L418-L451)).
   The Go module path `com.github.yzh44yzh/canon_d` is not a valid app ID because it contains
   `/`. Fyne warns that a unique ID is important because it is `unwise to change these after
   your first release`
   ([Fyne mobile docs](https://docs.fyne.io/started/mobile/)). Choose one before the first
   installed APK.
2. **Display name and output filename.** `-name "Canon D"` produces the app label `Canon D` but
   the APK filename `Canon_D.apk`, because non-alphanumeric characters become `_`
   ([tools build.go:44-51](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build.go#L44-L51),
   [build_androidapp.go:474-492](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build_androidapp.go#L474-L492)).
   Confirm whether build scripts should use that filename or a lowercase `-name`.
3. **Minimum public support claim.** Source sets Android 5.0/API 21, but Fyne's wiki still says
   Android 4.1+
   ([tools binres/sdk.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/binres/sdk.go#L14-L15),
   [wiki](https://github.com/fyne-io/fyne/wiki/Supported-Platforms)). Decide whether Canon D
   documents the manifest floor or the wiki claim.
4. **Exact NDK/build-tools/JDK versions.** Fyne pins none beyond NDK `>= r19c`
   ([tools env.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/env.go#L166-L173)).
   Choose exact versions in implementation and record them with build instructions.
5. **Debug signature scheme.** The source writer emits v1-style `CERT.RSA`
   ([tools writer.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/writer.go#L7-L45)),
   and issue [#154](https://github.com/fyne-io/tools/issues/154) is open. Verify the generated
   APK with `apksigner verify --verbose` before depending on v2+ behavior on newer devices.
6. **16 KB page devices.** Android 15 requires rebuild/support for native-code apps
   ([Android 15](https://developer.android.com/about/versions/15/behavior-changes-all#support-for-16-kb-page-sizes));
   no first-party statement confirms Fyne/Go `.so` alignment. Test on a 16 KB emulator/device.
7. **Fyne release discrepancy.** The module proxy/tag/changelog show `v2.8.3`, while the GitHub
   Releases API still listed `v2.8.2` first
   ([proxy](https://proxy.golang.org/fyne.io/fyne/v2/@latest),
   [releases](https://github.com/fyne-io/fyne/releases)). An earlier local research draft pinned
   `v2.8.2`; this revision updates the recommendation to `v2.8.3`. Recheck at implementation
   time.
8. **Deck embedding smoke test.** No first-party document explicitly says "go:embed is verified
   for Fyne Android"; the conclusion rests on Go toolchain behavior and Fyne's use of
   `go build`. Prove it with the first debug APK build
   ([Go embed docs](https://github.com/golang/go/blob/master/src/embed/embed.go#L5-L9),
   [tools build.go](https://github.com/fyne-io/tools/blob/v1.7.4/cmd/fyne/internal/mobile/build.go#L320-L361)).

## Source verification notes

- All cited external links were fetched successfully on 2026-11-27, including Fyne docs, raw
  Fyne tools/Fyne source at the pinned tags, Go docs and source, first-party Fyne/Go GitHub
  issues, and Android developer documentation.
- GitHub's Releases API was observed to lag the tag/module proxy for Fyne v2.8.3; the artifact
  therefore cites the module proxy, tag, and changelog together.
- This task changed only this artifact; no GitHub issues were touched.
