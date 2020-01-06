# ORCID Public Data Extractor

[![Actions Status](https://github.com/nad2000/orcid-pub-data-extractor/workflows/Go/badge.svg)](https://github.com/nad2000/orcid-pub-data-extractor/actions)

This is a simple tool that can help with filtering and extracting specific
activity records. Currently it supports filtering activity records that are
related to a specific country (country code).

	./orcid-pub-data-extractor  -h
	NAME:
	   extract-orcid - extract filtered data from ORCID profile activity public data

	USAGE:
	   main [global options] command [command options] FILE

	VERSION:
	   1.1.0

	COMMANDS:
	   help, h  Shows a list of commands or help for one command

	GLOBAL OPTIONS:
	   --country value, -c value  the country the record is related to (default: "NZ")
	   --type value, -t value     the record type: emp[ployment], edu[cation], work, fund[ing], peer[-review] ...
	   --output value, -o value   the output destination directory (default: "/home/rcir178/orcid-pub-data-extractor")
	   --search value, -s value   the search string
	   --regex value, -r value    the search regular expression (https://github.com/google/re2/wiki/Syntax)
	   --help, -h                 show help
	   --version, -v              print the version

