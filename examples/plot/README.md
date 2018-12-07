# Plot

The plot tool uses gonum/plot to draw a graph of an orbit list. It was mainly created to be able to debug the verification code.

The tool takes a path to a CSV position list and the path to an output file to put the png in.

## Usage

```bash
cd tools/plot
go build
./plot -in ../verification/1.csv -out ./out_1.png
```