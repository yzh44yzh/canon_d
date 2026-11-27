# Android statistics storage: primary-source findings

Ticket: [issue #2](https://github.com/yzh44yzh/canon_d/issues/2), “Research Android storage support for card statistics”.
Sources were accessed on **2026-11-27**. Only first-party Fyne, Android/AOSP, and Go sources are used. The two
Linux man pages cited below are the kernel project's own syscall documentation (`man-pages`, hosted by man7.org),
included only because Linux syscall durability semantics are not owned by Fyne, Android SDK docs, or Go docs.
No blogs, forums, or secondary summaries are used.

Fyne source is pinned to tag `v2.8.3` (`057615e811227d3324c6b08c6f2ee6f2064b3a5f`); the same behavior was
checked at `v2.6.0` (`2fac12d01d90f5b84eaf87a0eb4d185a18e8ac9f`). Go source is pinned to `go1.25.0`.

## Verdicts at a glance

| Question | Verdict |
| --- | --- |
| Can a Fyne app get its Android app-private root from `fyne.CurrentApp().Storage().RootURI()`? | **Yes.** It resolves to a `file:` URI rooted at Android `Context.getFilesDir()` plus `/fyne`. |
| How does that URI become an OS path? | Call `root.Path()`; do not strip `file://` manually. |
| Do `os.OpenFile`, `os.CreateTemp`, `os.Rename`, `File.Sync` work there? | **Yes** on `GOOS=android`; they use Go's Unix/Linux implementations. |
| Are storage permissions or runtime prompts needed? | **No** for app-private internal files. Scoped-storage rules concern external storage. |
| Is same-directory atomic replace + `fsync` durable? | **Best available first-party recipe** on ext4/f2fs internal storage; parent-directory sync is required. |
| Is `syscall.Flock(fd, LOCK_EX|LOCK_NB)` compatible? | **Yes** for normal Go Android builds; it uses a raw Linux syscall, not libc/Bionic. |

## 1. Fyne app-private storage path on Android

### Verified

- `fyne.Storage` declares `RootURI() URI` and is documented as managing “file storage inside an application
  sandbox” where files are “unique to the current application”; `fyne.App` exposes `Storage() Storage`, and
  `CurrentApp()` returns the process app.
  [Fyne API: fyne.Storage](https://docs.fyne.io/api/v2/fyne/storage/),
  [Fyne API: fyne.App](https://docs.fyne.io/api/v2/fyne/app/).
- `fyneApp.Storage()` returns the app's `store`; `(*store).RootURI()` calls
  `storage.ParseURI(s.a.storageRoot())`, falling back to `storage.NewFileURI(s.a.storageRoot())` when parsing
  fails. Because `storage.ParseURI` rejects strings with no scheme (“invalid URI, scheme must be present”), a
  plain filesystem root necessarily takes the `NewFileURI` path.
  [Fyne v2.8.3 `app/app.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/app/app.go),
  [Fyne v2.8.3 `app/storage.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/app/storage.go),
  [Fyne v2.8.3 `storage/repository/parse.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/storage/repository/parse.go).
- On Android, `storageRoot()` returns `app.RootConfigDir()`, documented in source as
  `// we are in a sandbox, so no app ID added to this path`. `RootConfigDir()` delegates to `rootConfigDir()`,
  which on Android reads `FILESDIR` and returns `filepath.Join(filesDir, "fyne")`.
  [Fyne v2.8.3 `app/preferences_android.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/app/preferences_android.go),
  [Fyne v2.8.3 `internal/app/config.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/internal/app/config.go),
  [Fyne v2.8.3 `internal/app/config_mobile_and.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/internal/app/config_mobile_and.go).
- Before Go `main.main`, Fyne's native Android entry point calls JNI
  `getFilesDir()` → `getAbsolutePath()` and `setenv("FILESDIR", filesdir, 1)`; the Go mobile driver propagates
  that environment into the Go process. The root is therefore `<Context.getFilesDir().getAbsolutePath()>/fyne`,
  not a hard-coded `/data/data/...` literal.
  [Fyne v2.8.3 `internal/driver/mobile/app/android.c`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/internal/driver/mobile/app/android.c),
  [Fyne v2.8.3 `internal/driver/mobile/app/android.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/internal/driver/mobile/app/android.go).
- `storage.NewFileURI` constructs `url.URL{Scheme: "file", Path: path}`; `String()` serializes the full URI,
  while `Path()` returns the stored path. `fyne.URI` documents `Path()` as the RFC 3986 path and has no
  `OSPath()`/`FilePath()` method.
  [Fyne v2.8.3 `storage/repository/parse.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/storage/repository/parse.go),
  [Fyne v2.8.3 `storage/repository/uri.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/storage/repository/uri.go),
  [Fyne API: fyne.URI](https://docs.fyne.io/api/v2/fyne/uri/).
- Fyne's desktop file repository uses `u.Path()` directly with `os.Stat`, `os.Open`, `os.Create`, and `os.Mkdir`,
  which confirms that for a `file:` URI `Path()` is the native path. On Android the mobile driver registers a
  repository for `file` and `content` whose reader/writer use JNI streams built from `u.String()`; that is an
  alternative abstraction, not a reason to avoid ordinary Go file APIs on the app-private path.
  [Fyne v2.8.3 `internal/repository/file.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/internal/repository/file.go),
  [Fyne v2.8.3 `internal/driver/mobile/file_android.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/internal/driver/mobile/file_android.go).
- Fyne ensures the root and a `Documents` child exist during app setup, but it does not create arbitrary
  subdirectories. `os.MkdirAll` (or Fyne's `storage.CreateListable`) is needed for a new statistics subdirectory.
  [Fyne v2.8.3 `app/app.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/app/app.go),
  [Fyne API: storage package](https://docs.fyne.io/api/v2/storage/pkg/).

### Recommendation

Use `fyne.CurrentApp().Storage().RootURI()` as the public API, verify `Scheme() == "file"`, and call `Path()` to
get the OS directory. Persist only paths relative to that root (e.g. `statistics/history.csv`), re-resolving the
root at startup. `String()` is for serialization, not for `os.OpenFile`.

### Unresolved risks

- Fyne guarantees `<getFilesDir()>/fyne`, not any literal `/data/user/0/<pkg>/...` prefix. Android docs warn that
  the `filesDir` path may change when an app is moved to adopted storage, so persisting the absolute path is unsafe.
  [Android `Context#getFilesDir()`](https://developer.android.com/reference/android/content/Context#getFilesDir()).
- If `FILESDIR` is unset, Fyne logs a warning and falls back to `/data/data` with the source comment
  “probably won't work”. If the app was created without a unique ID, `RootURI()` logs an error and returns
  `os.TempDir()`. Both need defensive handling.
  [Fyne v2.8.3 `internal/app/config_mobile_and.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/internal/app/config_mobile_and.go),
  [Fyne v2.8.3 `app/storage.go`](https://raw.githubusercontent.com/fyne-io/fyne/v2.8.3/app/storage.go).
- `Path()` is an RFC 3986 path; treating arbitrary URIs (`content:`, remote/custom schemes) as filesystem paths
  is not justified. Only the known `file:` root from `RootURI()` supports this.
  [Fyne API: fyne.URI](https://docs.fyne.io/api/v2/fyne/uri/).

## 2. Ordinary Go file APIs on `GOOS=android`

### Verified

- Go documents that “Using GOOS=android matches build tags and files as for GOOS=linux in addition to android tags
  and files”; the builder implements this by accepting `linux` files for `android`.
  [Go build constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints),
  [Go `src/go/build/build.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/go/build/build.go).
- `os.OpenFile` is documented as uniform across operating systems and delegates on Unix to `openFileNolog`, which
  calls the platform `open(...)` via `syscall.O_CLOEXEC`. The containing directory must already exist.
  [Go `os` docs](https://pkg.go.dev/os#OpenFile),
  [Go `src/os/file.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/os/file.go),
  [Go `src/os/file_unix.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/os/file_unix.go).
- `os.CreateTemp` creates with `O_RDWR|O_CREATE|O_EXCL` and mode `0600` before umask; the caller must remove the
  file. Its `dir` must be supplied explicitly for app data: Go's Android `tempDir()` defaults to
  `/data/local/tmp` when `TMPDIR` is empty, and the `os` docs state the default temp directory is “neither
  guaranteed to exist nor have accessible permissions”.
  [Go `os.CreateTemp` docs](https://pkg.go.dev/os#CreateTemp),
  [Go `src/os/tempfile.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/os/tempfile.go),
  [Go `src/os/file_unix.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/os/file_unix.go),
  [Go `src/os/file.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/os/file.go).
- `os.Rename` calls `syscall.Rename(oldname, newname)` on Unix. Go documents: “If newpath already exists and is
  not a directory, Rename replaces it” and “Even within the same directory, on non-Unix platforms Rename is not
  an atomic operation.” Android is a Linux/Unix target here.
  [Go `os.Rename` docs](https://pkg.go.dev/os#Rename),
  [Go `src/os/file_unix.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/os/file_unix.go).
- `(*os.File).Sync` is documented as committing “the current contents of the file to stable storage”; its Unix
  implementation calls `poll.FD.Fsync`, whose POSIX implementation wraps `syscall.Fsync`.
  [Go `os.File.Sync` docs](https://pkg.go.dev/os#File.Sync),
  [Go `src/os/file_posix.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/os/file_posix.go),
  [Go `src/internal/poll/fd_fsync_posix.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/internal/poll/fd_fsync_posix.go).

### Recommendation

Pass the Fyne root path into the storage layer, create `statistics/` with `os.MkdirAll`, and call
`os.CreateTemp(statisticsDir, ".history-*.csv")`. Do not use `os.CreateTemp("", ...)` for durable data on Android.

## 3. Android runtime and storage permissions

### Verified

- Android provides app-specific internal directories for persistent files (`getFilesDir()`) and cache
  (`getCacheDir()`); “Your app doesn't require any system permissions to read and write to files in these
  directories”, and other apps cannot access them. On Android 10 (API 29)+ they are encrypted.
  [Android app-specific storage](https://developer.android.com/training/data-storage/app-specific),
  [Android data storage overview](https://developer.android.com/training/data-storage).
- `Context#getFilesDir()` returns “the absolute path to the directory on the filesystem” used by
  `openFileOutput`; “No additional permissions are required for the calling app to read or write files under the
  returned path.” AOSP's `ContextImpl` sets `mFilesDir = new File(getDataDir(), "files")`, and
  `ApplicationInfo.dataDir` is documented as the package's default persistent-data directory.
  [Android `Context#getFilesDir()`](https://developer.android.com/reference/android/content/Context#getFilesDir()),
  [AOSP `ContextImpl.java`](https://android.googlesource.com/platform/frameworks/base/+/refs/heads/main/core/java/android/app/ContextImpl.java?format=HTML),
  [Android `ApplicationInfo#dataDir`](https://developer.android.com/reference/android/content/pm/ApplicationInfo#dataDir).
- Runtime permission requests are required for *dangerous permissions* on Android 6.0 (API 23)+, not for
  internal app-private files. `READ_EXTERNAL_STORAGE` “has no effect” starting in API 33, and
  `WRITE_EXTERNAL_STORAGE` has no effect for apps targeting Android 11 (API 30) or higher. Scoped storage changes
  access to *external* storage; `requestLegacyExternalStorage` is ignored when targeting Android 11.
  [Android runtime permissions](https://developer.android.com/training/permissions/requesting),
  [Android `READ_EXTERNAL_STORAGE`](https://developer.android.com/reference/android/Manifest.permission#READ_EXTERNAL_STORAGE),
  [Android `WRITE_EXTERNAL_STORAGE`](https://developer.android.com/reference/android/Manifest.permission#WRITE_EXTERNAL_STORAGE),
  [Android use cases / scoped storage](https://developer.android.com/training/data-storage/use-cases),
  [Android `R.attr#requestLegacyExternalStorage`](https://developer.android.com/reference/android/R.attr#requestLegacyExternalStorage),
  [Android 11 storage changes](https://developer.android.com/about/versions/11/privacy/storage).
- Android's application sandbox is enforced with Linux UID/GID facilities, and “native code is as sandboxed as
  interpreted code”; NDK `ANativeActivity.internalDataPath` is the app's internal data directory. Internal storage
  is not FUSE-mounted, unlike emulated external storage.
  [Android application sandbox](https://source.android.com/docs/security/app-sandbox),
  [NDK `ANativeActivity`](https://developer.android.com/ndk/reference/struct/a-native-activity.html),
  [Android scoped-storage internals](https://source.android.com/docs/core/storage/scoped).

### Recommendation

No manifest storage permission and no runtime prompt should be added for a statistics CSV under the Fyne root.
Keep all statistics files inside app-private internal storage; do not switch to `getExternalFilesDir()` or shared
storage for crash-critical data.

## 4. Atomic replace and durability

### Verified

- On Linux, `rename(2)` atomically replaces an existing non-directory target: “there is no point at which another
  process attempting to access newpath will find it missing”. Cross-filesystem renames fail with `EXDEV`, so the
  temporary file must be created in the same directory as the final file.
  [Linux `rename(2)`](https://man7.org/linux/man-pages/man2/rename.2.html),
  [Go `os.Rename`](https://pkg.go.dev/os#Rename).
- `fsync(2)` on a file is not enough to persist its directory entry: “an explicit fsync() on a file descriptor for
  the directory is also needed”. Go's `File.Sync()` exposes that syscall but does not add directory syncing.
  [Linux `fsync(2)`](https://man7.org/linux/man-pages/man2/fsync.2.html),
  [Go `os.File.Sync`](https://pkg.go.dev/os#File.Sync).
- AOSP's ext4 code itself warns that a newly created file whose parent directory is not synced leaves “a huge
  window during which a crash may lose the file”; ext4 and f2fs directory operations in AOSP kernels include
  `.fsync` handlers. Android lists ext4 and f2fs among supported filesystems.
  [AOSP ext4 `fsync.c`](https://android.googlesource.com/kernel/common/+/refs/tags/android13-5.15-2025-05_r13/fs/ext4/fsync.c?format=HTML),
  [AOSP ext4 `dir.c`](https://android.googlesource.com/kernel/common/+/refs/tags/android13-5.15-2025-05_r13/fs/ext4/dir.c?format=HTML),
  [AOSP f2fs `dir.c`](https://android.googlesource.com/kernel/common/+/refs/tags/android13-5.15-2025-05_r13/fs/f2fs/dir.c?format=HTML),
  [Android kernel filesystem support](https://source.android.com/docs/core/architecture/android-kernel-file-system-support).

### Recommendation

For a save: create temp in the final directory; write all bytes; `tmp.Sync()`; `tmp.Close()`; `os.Rename(tmp, final)`;
open the parent directory with `os.Open(dir)` and call `dir.Sync()`; check *every* error. Hold the advisory lock
for the whole sequence, and lock a stable lock file rather than the temp file that will be renamed. Treat a
directory-sync error as “not durably saved”.

### Unresolved risks

- Filesystem and mount policy are device-specific. AOSP documentation notes f2fs `fsync_mode` differences and
  ext4 compatibility heuristics (`auto_da_alloc`); these are not portable guarantees.
  [AOSP f2fs docs](https://android.googlesource.com/kernel/common/+/refs/heads/android14-6.1/Documentation/filesystems/f2fs.rst?format=HTML),
  [AOSP ext4 docs](https://android.googlesource.com/kernel/common/+/refs/tags/android12-5.10-2025-09_r4/Documentation/admin-guide/ext4.rst?format=HTML).
- Emulated external storage commonly goes through FUSE, and Android notes FUSE passthrough depends on device and
  OS version; first-party docs do not establish equivalent directory-fsync durability through that layer. This is
  another reason to stay on app-private internal storage.
  [Android FUSE passthrough](https://source.android.com/docs/core/storage/fuse-passthrough).

## 5. `syscall.Flock` with `LOCK_EX|LOCK_NB`

### Verified

- Go's build constraints treat Android as Linux, so Linux `syscall` files compile. `syscall.Flock` is declared in
  generated Linux wrappers with `//go:build linux && amd64` / `linux && arm64`; `LOCK_EX = 0x2`, `LOCK_NB = 0x4`
  match Android kernel UAPI (`#define LOCK_EX 2`, `#define LOCK_NB 4`).
  [Go build constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints),
  [Go `syscall.Flock`](https://pkg.go.dev/syscall#Flock),
  [Go `zsyscall_linux_amd64.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/syscall/zsyscall_linux_amd64.go),
  [Go `zsyscall_linux_arm64.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/syscall/zsyscall_linux_arm64.go),
  [Go `zerrors_linux_amd64.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/syscall/zerrors_linux_amd64.go),
  [AOSP `fcntl.h`](https://android.googlesource.com/kernel/common/+/refs/tags/android17-6.18-2026-06_r76/include/uapi/asm-generic/fcntl.h?format=HTML).
- Go's `Flock` calls `Syscall(SYS_FLOCK, fd, how, 0)` directly; `Syscall` reaches runtime assembly that issues the
  hardware syscall instruction. It does **not** call libc `flock()`, so Bionic's API history is not in the call
  path. Bionic nevertheless declares `flock(int, int)` in `sys/file.h` and lists `flock(int, int)` in
  `SYSCALLS.TXT`; AOSP kernels implement `SYSCALL_DEFINE2(flock, ...)`.
  [Go `zsyscall_linux_amd64.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/syscall/zsyscall_linux_amd64.go),
  [Go `syscall_linux.go`](https://raw.githubusercontent.com/golang/go/go1.25.0/src/syscall/syscall_linux.go),
  [Go ARM64 runtime syscall assembly](https://raw.githubusercontent.com/golang/go/go1.25.0/src/internal/runtime/syscall/asm_linux_arm64.s),
  [Bionic `sys/file.h`](https://android.googlesource.com/platform/bionic/+/refs/heads/main/libc/include/sys/file.h?format=HTML),
  [Bionic `SYSCALLS.TXT`](https://android.googlesource.com/platform/bionic/+/refs/heads/main/libc/SYSCALLS.TXT?format=HTML),
  [AOSP `fs/locks.c`](https://android.googlesource.com/kernel/common/+/refs/tags/android17-6.18-2026-06_r76/fs/locks.c?format=HTML).
- Bionic describes `flock()` as performing advisory locks; with `LOCK_NB`, a conflict returns `EAGAIN` in the AOSP
  kernel implementation. AOSP's app seccomp policy is generated so syscalls in `SYSCALLS.TXT` are allowed unless
  blocklisted, and current app blocklist/allowlist files contain no `flock` entry.
  [AOSP `fs/locks.c`](https://android.googlesource.com/kernel/common/+/refs/tags/android17-6.18-2026-06_r76/fs/locks.c?format=HTML),
  [Bionic app seccomp blocklist](https://android.googlesource.com/platform/bionic/+/refs/heads/main/libc/SECCOMP_BLOCKLIST_APP.TXT?format=HTML),
  [Bionic app seccomp allowlist](https://android.googlesource.com/platform/bionic/+/refs/heads/main/libc/SECCOMP_ALLOWLIST_APP.TXT?format=HTML).
- Go's `syscall` package documentation recommends `golang.org/x/sys` for new code, but the existing call is valid;
  `x/sys/unix.Flock` has the same signature if a migration is desired.
  [Go `syscall` docs](https://pkg.go.dev/syscall),
  [ `x/sys/unix.Flock`](https://pkg.go.dev/golang.org/x/sys/unix#Flock).

### Recommendation

Keep the `LOCK_EX|LOCK_NB` call. Handle `EWOULDBLOCK`/`EAGAIN` as “another instance is saving”, release the lock
with `LOCK_UN` or by closing the fd on every path, and place the lock file under the Fyne app-private root.

### Unresolved risks

- `flock` is advisory: cooperation requires all writers to use the same protocol. Filesystems may also provide
  their own `flock` handler instead of generic inode locking.
  [AOSP `fs/locks.c`](https://android.googlesource.com/kernel/common/+/refs/tags/android17-6.18-2026-06_r76/fs/locks.c?format=HTML).
- AOSP cannot prove every OEM kernel configuration. AOSP includes a conditional `flock` stub that can return
  `ENOSYS` if file locking is configured out; verify on target devices if this matters.
  [AOSP `kernel/sys_ni.c`](https://android.googlesource.com/kernel/common/+/refs/tags/android17-6.18-2026-06_r76/kernel/sys_ni.c?format=HTML).
- Gomobile itself is documented as experimental with a deliberately limited API surface; that caveat concerns
  mobile bindings, not Go's filesystem/syscall implementations.
  [ `golang.org/x/mobile`](https://pkg.go.dev/golang.org/x/mobile),
  [ `golang.org/x/mobile/app`](https://pkg.go.dev/golang.org/x/mobile/app),
  [gomobile command docs](https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile).

## 6. Canon D local-code audit and storage-adapter boundary recommendation

### What the current code actually does

The local implementation matches the premise that platform-independent CSV logic and platform I/O can be separated;
the current `history` package is already mostly that separation, but it has no injected interface.

- `history.Open(folder)` hardcodes `<folder>/.canon_d.csv` and `<folder>/.canon_d.csv.lock`, creates the lock with
  `os.OpenFile(..., os.O_CREATE|os.O_RDWR, 0600)`, locks it, then calls `readCSV`. It requires `folder` to already
  exist; it does not call `os.MkdirAll`. It returns a concrete `*history.Store`, not an interface.
  [Canon D `internal/history/store.go`](../../internal/history/store.go).
- `readCSV` uses `encoding/csv`, requires exactly the five-column header `deck,card_header,shows,correct_answers,incorrect_answers`, rejects invalid CSV, then calls `parseRow`. Go's `encoding/csv.Reader.FieldsPerRecord` defaults to zero: the field count is inferred from the first record, so subsequent records must have the same count. With a five-column header, short/long data rows produce `ErrFieldCount` before `parseRow`.
  [Canon D `internal/history/csv.go`](../../internal/history/csv.go),
  [Go `encoding/csv.Reader`](https://pkg.go.dev/encoding/csv#Reader).
- `parseRow` uses `strconv.ParseUint(..., 10, 64)`, validates identity and `positive answers require shows > 0`; `validateIdentity` requires non-whitespace valid UTF-8, rejects control characters, rejects `.`, `..`, `/`, and `\` in deck filenames, and preserves other exact text. `writeCSV` validates every merged row before writing and sorts identities deterministically.
  [Canon D `internal/history/validation.go`](../../internal/history/validation.go),
  [Canon D `internal/history/csv.go`](../../internal/history/csv.go).
- `Store.Save(results)` merges the caller's **cumulative current-run** results over the immutable `base` loaded at `Open`, checks uint64 overflow, validates, writes, and no-ops if `results` equals the last saved map. Repeated checkpoints do not double-count. `InitialStatistics()` returns a clone of `base`, excluding current-run results. `Save` expects the same complete cumulative map later, not a delta. `Close` closes the lock file, releasing the OS lock.
  [Canon D `internal/history/store.go`](../../internal/history/store.go),
  [Canon D ADR 0005](../adr/0005-saved-answer-totals.md),
  [Canon D ADR 0008](../adr/0008-history-csv-validation.md).
- `writeCSV` creates `.canon_d-*.tmp` in the destination directory, writes standard CSV, flushes, checks writer error, calls `file.Sync()`, closes, `os.Rename`s to the final file, and defers temp cleanup. It does **not** sync the parent directory. It checks rename error and never deletes the old file before rename. Current Canon D's own contract promises same-directory temp + file sync + atomic replacement; it does **not** promise directory-entry durability across power loss.
  [Canon D `internal/history/csv.go`](../../internal/history/csv.go),
  [Canon D ADR 0005](../adr/0005-saved-answer-totals.md).
- On `GOOS=android`, `lock_unix.go` is compiled because Android satisfies the `linux` build tag and does not satisfy `_windows`; it calls `syscall.Flock(..., LOCK_EX|LOCK_NB)`. The `lock_windows.go` implementation is excluded. The current `Close` releases the lock without an explicit `LOCK_UN`.
  [Canon D `internal/history/lock_unix.go`](../../internal/history/lock_unix.go),
  [Canon D `internal/history/lock_windows.go`](../../internal/history/lock_windows.go),
  [Go build constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints),
  [Go `syscall.Flock`](https://pkg.go.dev/syscall#Flock).
- `history` has no Bubble Tea, Fyne, or terminal-package imports. The TUI decides when to call `Save`; `history` performs no Fyne calls. `Session` statistics are passed in as `map[model.CardIdentity]model.CardStatistics`.
  [Canon D `internal/history`](../../internal/history),
  [Canon D `internal/session/session.go`](../../internal/session/session.go),
  [Canon D `internal/session/statistics.go`](../../internal/session/statistics.go).

### Reuse verdict for Fyne Android

**Yes for CSV parsing, validation, card identity, cumulative merging, retry-safe checkpoints, statistics snapshots, atomic same-directory replacement, file sync, and OS locking, with no Fyne-specific rewrite.** Keep those functions in `internal/history`. The only Android-aware responsibility outside that package is resolving a real app-private directory. A minimal call-site adapter is sufficient:

```go
func openMobileStatistics(app fyne.App) (*history.Store, error) {
    root := app.Storage().RootURI()
    if root.Scheme() != "file" {
        return nil, fmt.Errorf("statistics root is not a file URI: %s", root.String())
    }

    rootPath := root.Path() // OS path; do not use root.String() as a filename.
    if err := os.MkdirAll(rootPath, 0o700); err != nil {
        return nil, err
    }

    return history.Open(rootPath) // app-private root, or a pre-created child directory.
}
```

For `app.NewWithID("com.example.canond")`, `rootPath` is `<Context.getFilesDir()>/fyne`. Create a child with `os.MkdirAll(filepath.Join(rootPath, "statistics"), 0o700)` and pass that child instead if desired. App-private internal storage needs no Android storage permission or runtime request. No absolute app-private path should be persisted; resolve `RootURI().Path()` again at startup. A stable application ID keeps the Android app's data directory stable. Independence is per app installation/UID, not synchronization; uninstall/reinstall or clearing app data can remove history.
[Fyne root chain](#1-fyne-app-private-storage-path-on-android),
[Android app-specific storage](https://developer.android.com/training/data-storage/app-specific),
[Android application sandbox](https://source.android.com/docs/security/app-sandbox),
[Fyne app package](https://docs.fyne.io/api/v2/app/pkg/).

Use the same eight embedded deck filenames on every run; card identity includes exact deck filename and hierarchical header. Do not share this history file with TUI storage in the initial mobile design.

### Minimal adapter only if needed

There is no need to introduce an Android-specific CSV implementation. If testability or a stronger durability guarantee requires abstraction, keep `history.Open(folder)` as a wrapper and inject **one unexported low-level persistence implementation** inside `history`, preserving `Save`, `InitialStatistics`, `Close`, CSV validation, merging, and identity rules. The platform operations to factor are exactly those currently hardcoded:

- directory creation;
- read open;
- create-exclusive temporary file in the final directory;
- write/close;
- **file sync**;
- **same-directory rename**;
- **parent-directory sync** (if adding the stronger crash guarantee);
- temp removal;
- stable lock-file open; lock/unlock;
- close.

Do not create a separate `MobileStore` that re-parses or re-validates CSV, and do not let Fyne `storage.Repository` replace ordinary `os` operations on the known `file:` app-private root. In particular, do not abstract away `File.Sync`, same-directory `os.Rename`, stable lock inode, or cumulative-results semantics.
[Fyne storage API](https://docs.fyne.io/api/v2/fyne/storage/),
[Go `os`](https://pkg.go.dev/os).

### Background saves and concurrency

- `history.Store` is stateful and `session.Session` has in-memory mutable state; neither provides internal locking. On mobile, serialize checkpoint requests onto the same app event goroutine or protect both with a mutex. Fyne v2.8 `fyneDo=true` gives one Fyne event goroutine, but application-created goroutines must use `fyne.Do`/`fyne.DoAndWait`.
  [Fyne goroutines](https://docs.fyne.io/started/goroutines/).
- Keep one `history.Store` for the active Practice, checkpoint synchronously with cumulative `Session.Statistics()`, do not also open a second `history.Open` on the same file, and `Close` after the final save. The existing `.lock` then protects against accidental concurrent writers.
  [Canon D `internal/session/statistics.go`](../../internal/session/statistics.go),
  [Canon D `internal/history/store.go`](../../internal/history/store.go).
- Save at the existing completed-round checkpoints and when Android backgrounds the app; also preserve already acknowledged answers in `Session.Statistics()`. Backgrounding should request one idempotent checkpoint, not construct a new `Session` or apply current results twice. After process death, reopen history and construct a new Practice from saved totals, as selected by the Wayfinder map.
- If background return races a save, wait for the in-flight checkpoint or serialize it before accepting further practice actions. Current TUI already blocks after the first save failure; mobile should surface save failure and prevent claiming durable completion.
  [Canon D ADR 0003](../adr/0003-ui-and-color-scheme.md).

### Actual reuse check performed

Reading the local source establishes reuse; additionally, inspected the real Canon D `internal/history` package and ran, without modifying repository files:

```sh
CGO_ENABLED=0 GOOS=android GOARCH=arm64 \
  go test -c -o /tmp/opencode/canon-d-history-android-arm64.test ./internal/history
```

The **Android arm64 test-package compilation succeeded**. This is not an Android runtime test, a CGO/NDK mobile-app build, a device test, or a crash-durability test.

### Current reuse limitations

- No `os.MkdirAll`, `fsync(directory)`, injected filesystem, or application-level mutex in `history`.
- `history.Open` uses the current OS and lock implementation; compiling `lock_unix.go` for `GOOS=android` does not prove every OEM kernel exposes `flock` or has identical filesystem durability. See unresolved OEM risks above.
- Android app-private root resolution is a caller responsibility. `RootURI()` can fall back to temporary storage on app/root setup errors, so check scheme and root validity rather than silently writing statistics elsewhere.
- Independent mobile statistics do not imply import/export or sharing with TUI; those remain separate product decisions.

### Overall compatibility risk

Low for ordinary app-private internal storage on stock AOSP devices; moderate only for durability claims on
unknown OEM filesystems/mount options. No unverified assumption should be baked into code: resolve the root at
runtime, persist relative names, check every sync/rename error, and test on target devices for the durability and
seccomp caveats noted above.
