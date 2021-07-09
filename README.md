# genserverogz

A Sauerbraten map file parser/trimmer. Feed it a gunzipped .ogz file and it gives you a file containing only the header, game identifier, and entities sections. (You'll have to gzip that file back up in order to use it as an .ogz file on the server. See: [genall.sh](./genall.sh))

It specifically omits the map vars, extras, MRU ("most recently used"), lightmaps, blendmaps, vslots and octree geometry data from the original file. This reduces the file size by ~99.6%.

You can also use it as a parser of the processed sections, but it can't give you information on vslots, extras, MRU, lightmaps, blendmaps and the octree since it doesn't parse them.

## Usage

```
$ genserverogz -help
Reads uncompressed OGZ file on stdin and writes a minimal uncompressed OGZ file containing only the entity data to stdout.
Specify any print flag to only print the requested fields to stdout instead of the shrunk map data.
Usage of genserverogz:
  -ents
        print map entities
  -game
        print game identifier
  -vars
        print map vars (version 29+ only)
  -version
        print file format version
```

## Example

To analyse a map file:

```
$ gunzip --suffix=ogz --stdout ~/sauerbraten-code/packages/base/ot.ogz | genserverogz -version -vars -game -ents
OGZ file format version: 33
map variables: (26)
  skylight = 10196625
  lerpsubdiv = 2
  lightlod = 0
  waterfog = 20
  lighterror = 1
  lerpsubdivsize = 4
  lerpangle = 44
  skyboxcolour = 15466495
  causticscale = 40
  lightprecision = 42
  maptitle = Temple Ot by Nieb & ot4ku
  sunlight = 10525590
  lavacolour = 16728064
  sunlightyaw = 120
  sunlightscale = 2
  fog = 3000
  skytexturelight = 0
  skytexture = 0
  sunlightpitch = 55
  fogcolour = 10984320
  watercolour = 48840
  ambient = 2432801
  yawsky = 120
  blurskylight = 2
  waterfallcolour = 0
  skybox = penguins/yonder
game: fps
map entities: (89)
     0: type: 19, attrs:   1  -1   0   0   0, pos: 744.000000 304.000000 580.000000
     1: type: 20, attrs: 173   1   1   0   0, pos: 592.000000 624.000000 548.000000
     2: type:  1, attrs: 250  80  74  70   0, pos: 664.000000 568.000000 576.000000
     3: type:  1, attrs: 250  80  74  70   0, pos: 688.000000 432.000000 576.000000
     4: type: 19, attrs:   2  -1   0   0   0, pos: 432.000000 296.000000 516.000000
     5: type: 20, attrs: 320   2   1   0   0, pos: 432.000000 304.000000 580.000000
     6: type:  1, attrs: 250  80  74  70   0, pos: 408.000000 344.000000 556.000000
     7: type:  3, attrs: 269   0   0   0   0, pos: 440.000000 528.000000 548.000000
     8: type:  3, attrs: 128   0   0   0   0, pos: 688.000000 576.000000 548.000000
     9: type:  3, attrs: 320   0   0   0   0, pos: 496.000000 416.000000 532.000000
    10: type:  3, attrs:  18   0   0   0   0, pos: 664.000000 296.000000 580.000000
    11: type:  3, attrs: 297   0   0   0   0, pos: 384.000000 320.000000 516.000000
    12: type:  1, attrs:  50  10 100 200   0, pos: 432.000000 296.000000 528.000000
    13: type:  1, attrs:  60  10 100 200   0, pos: 744.000000 304.000000 592.000000
    14: type: 23, attrs:  30   0   0   0   0, pos: 744.000000 496.000000 548.000000
    15: type:  1, attrs:  40  10 100 200   0, pos: 736.000000 496.000000 556.000000
    16: type:  1, attrs:  15 250 250 240   0, pos: 664.000000 568.000000 588.000000
    17: type:  1, attrs:  15 250 250 250   0, pos: 408.000000 344.000000 564.000000
    18: type: 10, attrs:   0   0   0   0   0, pos: 604.000000 300.000000 580.000000
    19: type: 10, attrs:   0   0   0   0   0, pos: 456.000000 528.000000 628.000000
    20: type: 11, attrs:   0   0   0   0   0, pos: 668.000000 420.000000 612.000000
    21: type:  8, attrs:   0   0   0   0   0, pos: 436.000000 588.000000 548.000000
    22: type:  8, attrs:   0   0   0   0   0, pos: 444.000000 372.000000 516.000000
    23: type:  8, attrs:   0   0   0   0   0, pos: 660.000000 596.000000 548.000000
    24: type: 14, attrs:   0   0   0   0   0, pos: 496.000000 528.000000 532.000000
    25: type: 11, attrs:   0   0   0   0   0, pos: 484.000000 380.000000 580.000000
    26: type: 17, attrs:   0   0   0   0   0, pos: 664.000000 496.000000 612.000000
    27: type: 16, attrs:   0   0   0   0   0, pos: 584.000000 330.000000 532.000000
    28: type:  9, attrs:   0   0   0   0   0, pos: 380.000000 404.000000 516.000000
    29: type: 14, attrs:   0   0   0   0   0, pos: 724.000000 344.000000 548.000305
    30: type: 14, attrs:   0   0   0   0   0, pos: 652.000000 552.000000 564.000000
    31: type:  1, attrs: 100  60  57  56   0, pos: 592.000000 624.000000 560.000000
    32: type:  2, attrs:  90  11   0   0   0, pos: 363.894836 364.000000 527.999451
    33: type:  2, attrs:  90  11   0   0   0, pos: 363.896301 332.000000 527.987000
    34: type:  1, attrs:  15 255 250 250   0, pos: 688.000000 432.000000 588.000000
    35: type:  1, attrs:  12 250 225 200   0, pos: 444.000000 372.000000 516.000000
    36: type:  1, attrs:  12 250 225 200   0, pos: 380.000000 404.000000 516.000000
    37: type:  1, attrs:  12 250 225 200   0, pos: 436.000000 588.000000 548.000000
    38: type:  1, attrs:  12 250 225 200   0, pos: 484.000000 380.000000 580.000000
    39: type:  1, attrs:  12 250 225 200   0, pos: 604.000000 300.000000 580.000000
    40: type:  1, attrs:  12 250 225 200   0, pos: 660.000000 596.000000 548.000000
    41: type:  1, attrs:  12 250 225 200   0, pos: 668.000000 420.000000 612.000000
    42: type:  2, attrs: 195  14   0   0   0, pos: 600.657227 584.657593 511.551697
    43: type:  2, attrs: 195  14   0   0   0, pos: 624.110901 392.000000 511.956787
    44: type:  2, attrs:  90  11   0   0   0, pos: 363.908356 348.000000 528.007812
    45: type:  1, attrs:  60  60  59  58   0, pos: 736.000000 496.000000 608.000000
    46: type:  3, attrs: 132   0   0   0   0, pos: 676.000000 588.000000 548.000000
    47: type:  3, attrs: 359   0   0   0   0, pos: 500.000000 308.000000 580.000000
    48: type:  3, attrs: 114   0   0   0   0, pos: 700.000000 556.000000 612.000000
    49: type:  3, attrs: 226   0   0   0   0, pos: 452.000000 572.000000 548.000000
    50: type:  3, attrs: 308   0   0   0   0, pos: 400.000000 320.000000 516.000000
    51: type:  3, attrs:  90   0   0   0   0, pos: 716.000000 460.000000 612.000000
    52: type:  3, attrs: 282   0   0   0   0, pos: 444.000000 364.000000 580.000000
    53: type:  3, attrs:  88   0   0   0   0, pos: 686.000000 436.000000 548.000000
    54: type:  3, attrs: 221   0   0   0   0, pos: 496.000000 584.000000 532.000000
    55: type:  3, attrs:  25   0   0   0   0, pos: 632.000000 296.000000 580.000000
    56: type:  3, attrs: 301   0   0   0   0, pos: 384.000000 336.000000 516.000000
    57: type:  1, attrs:  12 250 225 200   0, pos: 676.000000 340.000000 548.000000
    58: type: 12, attrs:   0   0   0   0   0, pos: 676.000000 340.000000 548.000000
    59: type: 14, attrs:   0   0   0   0   0, pos: 426.000000 392.000000 596.000000
    60: type:  3, attrs:  20   0   0   0   0, pos: 648.000000 296.000000 580.000000
    61: type:  3, attrs: 135   0   0   0   0, pos: 684.000000 584.000000 548.000000
    62: type:  3, attrs: 245   0   0   0   0, pos: 444.000000 564.000000 548.000000
    63: type:  3, attrs: 259   0   0   0   0, pos: 436.000000 540.000000 548.000000
    64: type:  3, attrs: 290   0   0   0   0, pos: 444.000000 376.000000 580.000000
    65: type:  3, attrs: 353   0   0   0   0, pos: 486.000000 306.000000 580.000000
    66: type:  3, attrs:  65   0   0   0   0, pos: 686.000000 414.000000 548.000000
    67: type:  3, attrs:   8   0   0   0   0, pos: 470.000000 318.000000 516.000000
    68: type:  3, attrs: 105   0   0   0   0, pos: 692.000000 564.000000 548.000000
    69: type:  6, attrs:   0  80   0   0   0, pos: 432.000000 296.000000 512.000000
    70: type:  1, attrs:  30 200 180 160   0, pos: 732.000000 340.000000 580.000000
    71: type:  1, attrs:  30 200 180 160   0, pos: 732.000000 348.000000 580.000000
    72: type:  1, attrs:  30 200 180 160   0, pos: 732.000000 356.000000 580.000000
    73: type:  1, attrs:  30 200 180 160   0, pos: 732.000000 364.000000 580.000000
    74: type:  7, attrs:  50   0   0   0   0, pos: 734.000000 340.000000 584.000000
    75: type:  7, attrs:  50   0   0   0   0, pos: 734.000000 348.000000 584.000000
    76: type:  7, attrs:  50   0   0   0   0, pos: 734.000000 356.000000 584.000000
    77: type:  7, attrs:  50   0   0   0   0, pos: 734.000000 364.000000 584.000000
    78: type:  1, attrs:  30 200 180 160   0, pos: 652.000000 380.000000 580.000000
    79: type:  1, attrs:  30 200 180 160   0, pos: 644.000000 380.000000 580.000000
    80: type:  7, attrs:  50   0   0   0   0, pos: 652.000000 382.000000 584.000000
    81: type:  7, attrs:  50   0   0   0   0, pos: 644.000000 382.000000 584.000000
    82: type:  3, attrs:  98   0   0   0   0, pos: 721.892029 352.422546 548.000000
    83: type:  1, attrs:  12 250 225 200   0, pos: 496.000000 528.000000 532.000000
    84: type:  2, attrs: 270 200   0   0   0, pos: 431.764313 294.000000 528.000000
    85: type:  2, attrs:   0 200   0   0   0, pos: 746.000000 304.003906 592.010986
    86: type:  2, attrs:   0 194   0   0   0, pos: 736.000000 496.000000 545.000000
    87: type:  1, attrs:  12 250 225 200   0, pos: 724.000000 344.000000 548.000000
    88: type:  1, attrs:  12 250 225 200   0, pos: 664.000000 496.000000 612.000000
```

To only generate trimmed map files for use on a server use a script like [genall.sh](./genall.sh).

## Roadmap

We could save even more space by rewriting the input OGZ into a new one, omitting all map variables, MRU info and extra data.
