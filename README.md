# Jasper Converter

This is a Jasper report converter from version 6 JRXML to version 7 JRXML files using GO.

## Description

This tool converts Jasperreport v6 JRXML files into Jasperreport v7 files.

It can convert either a single file or all files in a directory.
The converted files are stored in the destination directory preserving the
directory structure of the source directory.

To convert a single JRXML file, use following command

```sh
> bin/darwin_arm64/convert -f tests6x/OrphanFooterReport.36.jrxml -d output 
Load and parsing JRXML ... tests6x/OrphanFooterReport.36.jrxml
```

## API

You can use the API of JasperGO to call it inside an application. Either use the path as input

```GO
ConvertFileToPath("./test", "./output")
```

or by providing the corresponding `io.Reader` structure.

To convert a single file use

```GO
var r io.Reader
s,_:=ConvertReader("OrphanFooterReport.36.jrxml", r)
... work on output s
```
