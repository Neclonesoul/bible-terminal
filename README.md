# Bible Terminal

**The King James Bible in your terminal.**

Bible Terminal is a small, offline-first Bible reader written in Go, built
around two editions:

- King James Version
- King James Bible 1611

It is designed for fast reading, searching and direct comparison between the
two texts without requiring a browser, account, network connection or
database server.

## Features

- Fully offline
- Self-contained executable
- KJV and KJV 1611
- 1611 spelling and typography preserved where supported by the source
- Side-by-side or stacked comparison
- Full-text search
- Interactive terminal reader
- Responsive phone and desktop layouts
- Book and chapter navigation
- Keyboard controls
- Embedded Bible datasets
- Works particularly well in Termux
- Linux-friendly

## Quick examples

Read a chapter:

```sh
bible Psalm 23
Read a verse:
bible John 3:16
Read the 1611 edition:
bible --1611 John 3:16
Compare editions:
bible --compare John 3:16
Search the KJV:
bible --search "still waters"
Search the 1611 text:
bible --1611 --search "beleeue"
Launch the interactive reader:
cat > README.md <<'EOF'
# Bible Terminal

**The King James Bible in your terminal.**

Bible Terminal is a small, offline-first Bible reader written in Go, built
around two editions:

- King James Version
- King James Bible 1611

It is designed for fast reading, searching and direct comparison between the
two texts without requiring a browser, account, network connection or
database server.

## Features

- Fully offline
- Self-contained executable
- KJV and KJV 1611
- 1611 spelling and typography preserved where supported by the source
- Side-by-side or stacked comparison
- Full-text search
- Interactive terminal reader
- Responsive phone and desktop layouts
- Book and chapter navigation
- Keyboard controls
- Embedded Bible datasets
- Works particularly well in Termux
- Linux-friendly

## Quick examples

Read a chapter:

```sh
bible Psalm 23
Read a verse:
bible John 3:16
Read the 1611 edition:
bible --1611 John 3:16
Compare editions:
bible --compare John 3:16
Search the KJV:
bible --search "still waters"
Search the 1611 text:
bible --1611 --search "beleeue"
Launch the interactive reader:
bible
Interactive controls
Key
Action
↑ / k
Scroll up
↓ / j
Scroll down
← / h
Previous chapter
→ / l
Next chapter
b
Book selector
v
Switch KJV / 1611
c
Comparison mode
g
Go to reference
/
Search
q
Quit
Responsive comparison
On narrow terminals, comparison is displayed vertically:
KJV · 16
For God so loved the world...

────────────────────────────

1611 · 16
¶ For God so loued þ world...
On sufficiently wide terminals the same comparison becomes two columns:
KJV                         │ KJV 1611
────────────────────────────┼────────────────────────────
For God so loved the world  │ ¶ For God so loued þ world
that he gave...             │ that he gaue...


