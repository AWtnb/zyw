# fuzZY Walker

Traverse directory and select path with fuzzy-finder.


```
> .\zyw.exe -h

  -all
        switch to search including file
  -src string
        source directory
  -exclude string
        path to skip searching (comma-separated)
  -offset 0
        Specify the directory to start file traversal, by the number of layers from the current directory.
        `0` for the current directory, `1` for the parent directory, `2` for its parent directory, and so on.
        If this value is negative, the path is traversed back to the directory containing the file `.root`. If no `.root` file is found, the current directory is used as the root of the search. (default -1)
```

The directory where the [`.root`](.root) is placed should be the starting point of the traversal.
