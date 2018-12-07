# Validation

This folder contains some bits to check our algorithms against other astronomical packages.

Currently we only have the [Poliastro](http://docs.poliastro.space/en/latest/)

In the tools folder above this there is a diff tool to compaire the output of the go
program in here with the other systems in the sub folders.

The output from our library can be generated using the following commands:

```bash
go build
./verification
```

This will generate a series of CSV files with the object id.

New verification packages should output the same format as far as possible.
The diff tool is able to handle differences in the number format as long as they are parseable as float64 numbers