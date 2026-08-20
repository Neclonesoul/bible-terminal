::: {align="center"}
# ✠ BIBLE TERMINAL ✠

### VERBUM TUUM VERITAS

**THY WORD IS TRUTH**

*A fast, offline King James Bible reader for the command line.*

**KJV · KJV 1611 · COMPARISON · SEARCH · NO ACCOUNT · NO CLOUD**
:::

------------------------------------------------------------------------

> **"Sanctify them through thy truth: thy word is truth."**
>
> --- John 17:17, KJV

Bible Terminal puts the King James Bible where a text belongs: **in your
possession**.

No browser.\
No account.\
No advertisements.\
No recommendation algorithm.\
No telemetry.\
No Internet connection required.

Just the text.

Read the standard KJV, enter the original-spelling world of the 1611
edition, or place both texts together and compare them verse by verse.

Built in Go. Designed for the terminal. Small enough to carry anywhere.

------------------------------------------------------------------------

## `> THE WORD`

``` text
$ bible John 3:16

THE HOLY BIBLE · KJV · JOHN 3
────────────────────────────────────────────────────────────────────────

16  For God so loved the world, that he gave his only begotten Son,
    that whosoever believeth in him should not perish, but have
    everlasting life.
```

And the 1611 text:

``` text
$ bible --1611 John 3:16

THE HOLY BIBLE · 1611 · JOHN 3
────────────────────────────────────────────────────────────────────────

16  ¶ For God so loued þ world, that he gaue his only begotten Sonne:
    that whosoeuer beleeueth in him, should not perish, but haue
    euerlasting life.
```

------------------------------------------------------------------------

# KJV ↔ 1611

Bible Terminal was built around a simple idea:

**comparison should be effortless.**

``` sh
bible --compare John 3:16
```

On a wide terminal:

``` text
KJV                                      │ KJV 1611
─────────────────────────────────────────┼────────────────────────────────────────
16  For God so loved the world, that he  │ 16  ¶ For God so loued þ world, that
    gave his only begotten Son, that     │     he gaue his only begotten Sonne:
    whosoever believeth in him should    │     that whosoeuer beleeueth in him,
    not perish, but have everlasting     │     should not perish, but haue
    life.                                │     euerlasting life.
```

On narrow terminals and phones, comparison automatically becomes a
readable stacked layout.

No websites to switch between. No tabs. No copy-and-paste.

**One command. Two witnesses.**

------------------------------------------------------------------------

# INTERACTIVE READER

Run:

``` sh
bible
```

and Bible Terminal becomes an interactive reader.

``` text
┌──────────────────────────────────────────────────────────────┐
│                       THE HOLY BIBLE                         │
│                    KING JAMES VERSION                        │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│                         PSALMS XXIII                         │
│                                                              │
│  1  The LORD is my shepherd; I shall not want.              │
│                                                              │
│  2  He maketh me to lie down in green pastures: he          │
│     leadeth me beside the still waters.                      │
│                                                              │
├──────────────────────────────────────────────────────────────┤
│ PSALMS 23 · KJV · 0%                                         │
│ ↑↓ scroll · ←→ chapter · b books · v · c · g · / · q        │
└──────────────────────────────────────────────────────────────┘
```

### Controls

      Key     Function
  ----------- ----------------------------
   `↑` / `k`  Scroll up
   `↓` / `j`  Scroll down
   `←` / `h`  Previous chapter
   `→` / `l`  Next chapter
      `b`     Book selector
      `v`     KJV ↔ 1611
      `c`     Comparison mode
      `g`     Go directly to a reference
      `/`     Search
      `q`     Quit

------------------------------------------------------------------------

# SEARCH THE SCRIPTURES

Search the entire KJV locally:

``` sh
bible --search "fear not"
```

Search the 1611 text:

``` sh
bible --1611 --search "beleeue"
```

Search happens entirely on your machine.

The Bible database is embedded inside the executable.

**Internet required: 0 bytes.**

------------------------------------------------------------------------

# DOWNLOAD

Download the binary for your machine from the **Releases** section.

## Windows

**Most Windows PCs**

`bible-terminal-v0.1.0-windows-amd64.exe`

→ Intel or AMD 64-bit Windows.

**Windows ARM**

`bible-terminal-v0.1.0-windows-arm64.exe`

→ Windows ARM PCs, including Snapdragon systems.

## macOS

**Apple Silicon**

`bible-terminal-v0.1.0-macos-arm64`

→ M1 and newer Apple Silicon Macs.

**Intel Mac**

`bible-terminal-v0.1.0-macos-amd64`

→ Intel-based Macs.

## Linux

**Most Linux PCs**

`bible-terminal-v0.1.0-linux-amd64`

→ Standard Intel/AMD x86-64 Linux systems.

**Linux ARM64**

`bible-terminal-v0.1.0-linux-arm64`

→ ARM64 Linux, Raspberry Pi and compatible SBCs.

## Android / Termux

`bible-terminal-v0.1.0-android-arm64`

→ ARM64 Android devices running Termux.

Bible Terminal was developed and tested extensively from Android itself.

------------------------------------------------------------------------

# INSTALL

## Linux

``` sh
chmod +x bible-terminal-v0.1.0-linux-amd64
sudo install bible-terminal-v0.1.0-linux-amd64 /usr/local/bin/bible
bible
```

## macOS

For Apple Silicon:

``` sh
chmod +x bible-terminal-v0.1.0-macos-arm64
sudo install bible-terminal-v0.1.0-macos-arm64 /usr/local/bin/bible
bible
```

## Android / Termux

``` sh
chmod +x bible-terminal-v0.1.0-android-arm64
install bible-terminal-v0.1.0-android-arm64 "$PREFIX/bin/bible"
bible
```

## Windows

Download `bible-terminal-v0.1.0-windows-amd64.exe`, rename it to
`bible.exe` if desired, and run from PowerShell or Windows Terminal:

``` powershell
.\bible.exe John 3:16
```

------------------------------------------------------------------------

# BUILD IT YOURSELF

Don't trust binaries?

Good.

Read the source and compile it yourself.

Requirements:

-   Go
-   Git

``` sh
git clone https://github.com/Neclonesoul/bible-terminal.git
cd bible-terminal

go mod download
go build -o bible ./cmd/bible

./bible
```

The application does not require a database server or external Bible
API.

------------------------------------------------------------------------

# WHAT IS INSIDE?

``` text
bible-terminal/
│
├── cmd/bible/          command-line application
│
├── internal/
│   ├── data/           embedded Bible texts
│   ├── search/         offline full-text search
│   ├── tui/            interactive terminal reader
│   └── version/        release information
│
├── scripts/
│   └── build-data.go   reproducible dataset builder
│
├── SOURCES.md          textual provenance
├── LICENSE
└── README.md
```

The Bible datasets are compiled directly into the executable using Go's
embedding facilities.

The program therefore carries its library with it.

------------------------------------------------------------------------

# THE TEXTS

Bible Terminal currently contains:

``` text
King James Version       31,102 verse records
King James Bible 1611    36,702 verse records
──────────────────────────────────────────────
Total                    67,804 records
```

The 1611 dataset contains the additional books represented by its
80-book source.

Text provenance and upstream sources are documented in **SOURCES.md**.

Nothing should enter an archival text project without provenance.

------------------------------------------------------------------------

# WHY A TERMINAL BIBLE?

Because plain text is powerful.

A terminal Bible is:

-   fast;
-   searchable;
-   scriptable;
-   portable;
-   low-bandwidth;
-   usable over SSH;
-   usable on old hardware;
-   usable on a phone;
-   independent of browsers;
-   independent of cloud services;
-   easy to archive;
-   easy to inspect;
-   difficult to make obsolete.

The interface is intentionally austere.

The subject is not.

------------------------------------------------------------------------

# PRINCIPLES

### THE TEXT FIRST

The interface exists to reveal the text, not compete with it.

### OFFLINE FIRST

Once installed, reading Scripture should not depend upon a server
remaining online.

### VERIFY

The project records the provenance of its source texts and provides
SHA-256 checksums for release binaries.

### COMPARE

Differences should be visible, not buried behind menus.

### OWN YOUR TOOLS

The source is available. The data is local. The executable is
self-contained.

Build it yourself if you wish.

------------------------------------------------------------------------

# VERIFY YOUR DOWNLOAD

Every release includes `SHA256SUMS`.

On Linux, macOS or Termux:

``` sh
sha256sum -c SHA256SUMS
```

A downloaded executable should be verifiable independently of the page
from which it came.

------------------------------------------------------------------------

# EXAMPLES

``` sh
# Open the reader
bible

# Psalm 23
bible Psalm 23

# John 3:16
bible John 3:16

# Original-spelling 1611
bible --1611 John 3:16

# Compare the editions
bible --compare John 3:16

# Search KJV
bible --search "the truth shall make you free"

# Search 1611
bible --1611 --search "beleeue"
```

------------------------------------------------------------------------

# ROADMAP

### Bible Terminal

The terminal application will remain focused on speed, portability and
textual comparison.

Future work may include:

-   improved reference parsing;
-   bookmarks;
-   reading history;
-   piping and scripting improvements;
-   packaged installation methods;
-   automated cross-platform releases.

### Bible Illuminated

A separate graphical companion is planned: **Bible Illuminated**.

Where Bible Terminal is deliberately austere, Bible Illuminated will
explore the visual world of the early printed King James Bible:

-   1611-inspired typography;
-   ornamental borders;
-   decorated initials;
-   historical page composition;
-   parallel KJV / 1611 study;
-   modern responsive reading;
-   offline-first architecture.

Two interfaces.

One textual foundation.

------------------------------------------------------------------------

::: {align="center"}
## ✠ VERBUM TUUM VERITAS ✠

**"Thy word is true from the beginning."**

Psalm 119:160

`KJV // 1611 // OFFLINE // OPEN SOURCE`
:::
