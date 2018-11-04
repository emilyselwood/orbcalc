from astropy import units as u
from astropy.time import Time

from poliastro.bodies import *
from poliastro.twobody import Orbit
from poliastro.plotting import OrbitPlotter, plot

def process(name, orbit, path):
    f = open(name.replace(" ", "_") + ".csv","w+")
    for i in range(1, 366):

        r = orbit.propagate(i * u.day).state.r
        f.write("{},{},{},{},{}\n".format(name, i, r[0].value, r[1].value, r[2].value))
    f.close()


testObjects = {
    "1996 PW": Orbit.from_classical(
        Sun, 
        a = 3.79035922723884e+10 * u.km,
        ecc = 0.9901593 * u.one,
        inc = 0.5228416517687837 * u.rad,
        raan = 2.519967809619083 * u.rad,
        argp = 3.169512336568096 * u.rad,
        nu = 0.03539440456581901 * u.rad
    )
}


if __name__ == "__main__":
    for name, orbit in testObjects.items():
        process(name, orbit, "./")