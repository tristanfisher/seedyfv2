# seedyfv2

wip

Code for CFBF / CDFV2 exploration.


## quick overview

`Compound File Binary` (aka `Object Linking and Embedding` `(OLE)` or `Component Object Models` `(COM)`) is a Microsoft
format that embeds multiple files into a single file.  The claimed purpose (see file "MS-CFB") is "improved efficiency" 
and ease of transport over using multiple files.

This file's specification reflects its origin during the dominance of 32-bit systems.  This is worth noting as it explains 
32-bit uints and sector boundaries.

This file is a virtual filesystem with branching implementations of `sectors` - abbreviated `SECT`, which are units within the compound file.  There are 2 special kinds of sectors:

1. File header

    The header is always at offset 0, accessible via `StructuredStorageHeader`.

2. Range lock

    A "range lock" exists at the offset of 2**32 bytes (32-bit max addressable - 4294967296)

Other sectors can exist anywhere else in the file.  Other sectors include:

- File Allocation Tables (FAT)
- Mini file allocation tables (MiniFAT)
- Double-indirect FAT (DIFAT) (represents storage of FAT sectors, final 4 bytes reserved for DIFAT chaining)
- Streams
- Range Locks (concurrency control for byte ranges within the file.  not required to be allocated on 512 byte CFBs (`header.MajorVersion`) (open a PR if you have a good explanation as to why))

`Sectors` are arranged in `sector chains`, which are linked lists.  

- Sectors can be in non-consecutive locations.
- Sectors that are unallocated or free are not part of a sector chain
- Sentinel/magic values are used to denote chain termination or free sectors



## documentation / spec

The [Library of Congress](https://www.loc.gov/preservation/digital/formats/fdd/fdd000380.shtml) has a good index of information.

Some useful documents have been vendored to the [./doc](./doc) folder.