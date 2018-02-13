# gget
GGet is a library which mimics the functionality of wget.
GGet will read in a file, download and save all resources to the filesystem.

## To Build
`go build -o gget src/cmd/gget/main.go`

## Background
This is a maintained fork of the original and was adapted from this usecase.

```
I had a special application that required the fetching of a lot of files listed in a text file.
Normally in cases like this I would use wget -i, but since the urls contained unicode characters needing URL-encoding, 
Wget messed up when naming files and folders by failing to decode it back for a proper filename.

So here is my very go-ish solution to the problem.
```

## Input
Example input:

````
http://sigmaitc.se/wp-content/uploads/sites/4/2014/02/Panelbild_jobs-800x333.jpg
http://sigmaitc.se/wp-content/uploads/sites/4/2015/03/sigma-compleo.jpg
http://sigmaitc.se/wp-content/uploads/sites/4/2014/10/mailbanner-kille.jpg
http://sigmaitc.se/wp-content/uploads/sites/4/2014/02/medarbetare_tomas_hellman.jpg
http://sigmaitc.se/wp-content/uploads/sites/4/2014/02/medarbetare_lars_littorin.jpg

http://www.kammarkollegiet.se/sites/default/files/2014-06-12, Remissyttrande - Översvämningsmyggor vid Nedre Dalälven (PDF, 1.01 MB).pdf
````