# orbcalc

A library to do orbital mechanics in go

Currently very basic and work in progress, the basic orbital propogation with the mean motion method should work for most cases. 
Hyperbolic and parabolic orbits have not been tested, if you find bugs please let us know.

Current Example in main.go which reads in the MPC orbit file propogates them forward by one day and then writes the position vectors to a file

There is a lot still to do:

* Reference frame transformations.
* Built in data for earth orbits and major planets
* Benchmarking
* Documentation

If you want to help with these please feel free to get in contact.