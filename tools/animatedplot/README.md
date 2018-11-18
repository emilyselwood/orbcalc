# Animated Plot

This creates a series of images with the orbits moved by one day at a time.

```bash
go build
./animatedplot -in /data/MPCORB.DAT -out /data/frames
cd /data/frames
ffmpeg -f image2 -r 60 -i frame_%05d.png -c:v libx264 -s 1000x1000 ../out.avi
```