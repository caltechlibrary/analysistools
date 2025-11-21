---
title: analysistools
abstract: "An experimental package for Archival text analysis on the file system. It
provide a phrase check tool and some file type analysis based on file extensions."
authors:
  - family_name: Doiel
    given_name: R. S.
    id: https://orcid.org/0000-0003-0900-6903



repository_code: https://github.com/caltechlibrary/analysistools
version: 0.0.2


programming_language:
  - Go

keywords:
  - file system
  - analysis
  - text views

date_released: 2025-11-20
---

About this software
===================

## analysistools 0.0.2

Swapped filepath.Walk() with filepath.WalkDir() for improved performance.
Refined checking process by implementing a tokenizer that returns words and new lines.
Upgrade methods for matching to handle prefix and suffix asterix as well as as exact match.
Improved proximity checking and testing.

### Authors

- R. S. Doiel, <https://orcid.org/0000-0003-0900-6903>






An experimental package for Archival text analysis on the file system. It
provide a phrase check tool and some file type analysis based on file extensions.


- GitHub: <https://github.com/caltechlibrary/analysistools>
- Issues: <https://github.com/caltechlibrary/analysistools>

### Programming languages

- Go




### Software Requirements

- Go >= 1.25
- CMTools >= 0.0.40


### Software Suggestions

- Pandoc &gt;&#x3D; 3
- GNU Make &gt;&#x3D; 3
- PageFind &gt;&#x3D; 1.4.0


