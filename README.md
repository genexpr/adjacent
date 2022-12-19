## Overview

This program displays a translation of an English word (or text) in other languages of a single linguistic group.

The goal is to explore and showcase the similarity of roots in languages that are considered closely related.

Some Wikipedia pages have tables for a word in multiple related languages. This program allows one to query for any word
and get a similar overview.

**You would need a token for the DeepL Translation API to use this.**

Set your API token as the value for the `DEEPL_TRANSLATE_TOKEN` environment variable.

## Usage

Choose from the Slavic, Germanic or Romance groups.

```
$ adjacent -text=sky -group=slavic
```

The above will display how "sky" is written in a few Slavic languages,
such as Czech, Polish, Ukrainian, and so on.

For help:

```
$ adjacent -h
```
