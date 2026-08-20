# Bible Terminal — Text Sources

Bible Terminal provides offline access to King James Bible texts embedded
directly in the executable.

## King James Version

Text dataset derived from:

aruljohn/Bible-kjv

https://github.com/aruljohn/Bible-kjv

This upstream repository contains the 66-book King James Version in JSON
format and is distributed under the MIT License.

## King James Bible 1611

1611 text dataset derived from:

aruljohn/Bible-kjv-1611

https://github.com/aruljohn/Bible-kjv-1611

This source contains the original-spelling KJV 1611 text represented as
80 books, including the books commonly grouped as the Apocrypha.

## 1611 Typeface — future graphical project

The following project is NOT the source of Bible Terminal's 1611 text:

ctrlcctrlv/kjv1611

https://github.com/ctrlcctrlv/kjv1611

It is a digital restoration of the typeface found in the 1611 King James
Bible. It is being investigated separately for the future Bible Illuminated
graphical project.

## Reference terminal Bible projects

Other terminal Bible implementations examined during the project's design
include:

layeh/kjv
https://github.com/layeh/kjv

LukeSmithxyz/kjv
https://github.com/LukeSmithxyz/kjv

Bible Terminal is its own Go implementation and does not require these
programs at runtime.

## Embedded application data

The generated datasets are:

- internal/data/kjv.tsv
- internal/data/kjv1611.tsv

They are compiled into the Bible Terminal executable using Go's embed
facility.

Current generated record counts:

- KJV: 31,102
- KJV 1611: 36,702
- Total: 67,804

## Reproducibility

The source-data conversion program is:

scripts/build-data.go

The intention is for every embedded text dataset to remain traceable to its
upstream source.
