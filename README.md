# Sitemap builder
 The program generates xml of the sitemap in the standard sitemap protocol.
 It crawls through webpages for given domain by using breadth first search algorythm.
 Then it writes the xml to a given file
# Installation:
    git clone https://github.com/Hellmick/sitemap-builder.git
    cd sitemap-builder
    go build
# Usage:
    ./sitemap-builder -u <URL> 
    Parameters:
      -d <depth of bfs> 
      -f <filename>
   
