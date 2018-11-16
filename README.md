# orbcalc

A library to do orbital mechanics in go

Currently very basic and work in progress, the basic orbital propogation with the mean motion method should work for most cases.
Hyperbolic and parabolic orbits have not been tested, if you find bugs please let us know.

Example in main.go which reads in the MPC orbit file propogates them forward by one day and then writes the position vectors to a file

There is a lot still to do:

* Reference frame transformations.
* Built in data for earth orbits and major planets
* Benchmarking
* Documentation

If you want to help with these please feel free to get in contact.

## Reason

This project is designed to alow you to work out the position in space of an object after some time given the normal orbital elements.

The main usecase is to be able to plot the locations of asteroids over time.

### Design Goals

1) Be Accurate
1) Be Fast
1) Be Easy To Use

## Contributing

Fantastic. We welcome an help you can give. We especially welcome bug reports and case studies of uses. If you have managed to successfully use this project
please let us know. If you have found a pain point please let us know, we can probably make it easier to use. If you are not sure if something is a bug please
rase it any way. Worst case it is something we need to document better.

If you want to provide code support to the project we use the "usual" github process, issues, forks and pull requests.

### Building from source

Prerequistits:

* Golang 1.11+

```bash
git clone git@github.com:wselwood/orbcalc.git
cd orbcalc
go build
```

We use the Go module system which should take care of the dependencies for you.

## Thanks

This project owes a great debt of thanks to the [poliastro project](https://github.com/poliastro/poliastro) for the algorithms and examples of how things should be done.