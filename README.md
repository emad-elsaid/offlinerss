OfflineRSS
==========

Downloads RSS feed URLs and split each entry to an XML file and write them to
`~/rss/INBOX`.

Offlinerss will not write the entry to the inbox directory if it exists in any
subdirectory in `~/rss`. This means if you read an article and you want to move
it to another directory you can move it to `~/rss/archive` for example and it
won't be created again when you run `offlinerss`.

## Usage

Clone or copy `offlinerss` to any directory in your `PATH`, put a config file
that include your RSS feeds URLs in `~/rss/config.yml` as follows

```yaml
urls:
   - https://domain.tld/path/to/feed.rss
   - https://domain.tld/path/to/another/feed.atom
```

Run `offlinerss` when you want to update your RSS feeds.

## How it works

- offlinerss will read `~/rss/config.yml`
- Will read each URL
- For each entry it will write its XML to a file in `~/rss/INBOX/` with a name
  in this format `sha1(feed url)-sha1(entry id/guid).rss`
- If the file `sha1(feed url)-sha1(entry id/guid).rss` exist under any
  subdirector in `~/rss` it will not be written again.
- The rest of the RSS feed is written to `~/rss/.meta` to a file with a name in
  this format `sha1(feed url).rss`


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
