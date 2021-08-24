OfflineRSS
==========

A rewrite of the [offlinerss](https://rubygems.org/gems/offlinerss) ruby gem in Go.

Downloads RSS feed URLs and split each entry to an XML file and write them to
`~/rss/INBOX`.

Offlinerss will not write the entry to the inbox directory if it exists in any
subdirectory in `~/rss`. This means if you read an article and you want to move
it to another directory you can move it to `~/rss/archive` for example and it
won't be created again when you run `offlinerss`.

## Installation

Using Go toolchain:

```
$ go install github.com/emad-elsaid/offlinerss@v0.1.3
```

Copy [config.example.json](config.example.json) to `~/rss/config.json` and fill it with your RSS feeds URLs

## Usage

Run `offlinerss` when you want to download your RSS feeds.

## How it works

- offlinerss will read `~/rss/config.json`
- Will read each URL in parallel
- For each entry it will write its XML to a file in `~/rss/INBOX/` with a name
  in this format `sha1(feed url)-sha1(entry content).rss`
- If the file `sha1(feed url)-sha1(entry content).rss` exist under any
  subdirectory in `~/rss` it will not be written again.
- The rest of the RSS feed is written to `~/rss/.meta` to a file with a name in
  this format `sha1(feed url).rss`

## Design guidelines

- No dependencies. Just Go! itself
- Sync in parallel as the process is mostly IO bound.
- No sub packages.
- No unnedeeded abstractions

## Benefits of using the file system as a database

- Any application can read it
- You can search with any file search tool (grep, ag, ripgrep..etc)
- You can move the files as you wish to any directory
- You can version it with Git or sync it with `rsync`
- You can replace `offlinerss` with another implementation and the data doesn't
  need any transformation, as offlinerss doesn't do any transformation at all,
  just splits it to file.
- You can view the files with any simple application that reads the file and
  renders it to HTML or any other format

## Motivation

I wrote about this script for the first time couple days ago in my blog
describing [ why I was set to create it ](https://www.emadelsaid.com/download-RSS-offline/)


## Upgrade from ruby gem to the Go implementation

The Go implementation has a small difference in the rss item file name. Instead of getting the item id/guid/link and hashing it with SHA1 I used the item itself, this was simpler to do as it doesn't require parsing the feed or the item at all.

So if you used the ruby implementation and you ran the Go implementation on the same rss directory it will write existing rss item again to INBOX as the filename is different.

So to upgrade, make sure you have read and moved all your RSS items from INBOX then run the go implementation which will write again all rss items from your feeds to INBOX. After that move all items in INBOX to another directory say "Duplicates" or "Trash".

Another difference is the `config` file. now is in JSON format instead of YAML
