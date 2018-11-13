from astropy import units as u
from astropy.time import Time

from poliastro.bodies import *
from poliastro.twobody import Orbit
from poliastro.plotting import OrbitPlotter, plot
from poliastro.twobody.classical import coe2rv
from poliastro.core.elements import rv_pqw

def process(name, orbit, path):
    f = open(name.replace(" ", "_") + ".csv","w+")
    r, v = coe2rv(orbit.attractor.k.to(u.km ** 3 / u.s ** 2).value,
                      orbit.p.to(u.km).value,
                      orbit.ecc.value,
                      orbit.nu.to(u.rad).value)
    print(name, "r:", r)
    print(name, "v:", v)
    for i in range(0, 366):
        if i == 0:
            e = orbit.epoch
            r = orbit.state.r
        else :
            o = orbit.propagate(i * u.day)
            e = o.epoch
            r = o.state.r
        f.write("{},{},{},{},{}\n".format(name, e, r[0].value, r[1].value, r[2].value))
    f.close()


testObjects = {
    "1996 PW": Orbit.from_classical(
        Sun, 
        a = 3.79035922723884e+10 * u.km,
        ecc = 0.9901593 * u.one,
        inc = 0.5228416517687837 * u.rad,
        raan = 2.519967809619083 * u.rad,
        argp = 3.169512336568096 * u.rad,
        nu = 0.03539440456581901 * u.rad,
        epoch = Time('2018-01-01T00:00:00Z', scale='utc', format='isot')
    ),
    "1": Orbit.from_classical(
        Sun, 
        a = 4.1394459238740003e+08 * u.km,
        ecc = 0.0755347 * u.one,
        inc = 0.1848916288429445 * u.rad,
        raan = 1.4016725260132445 * u.rad,
        argp = 1.2761023695175595 * u.rad,
        nu = 6.147582300011738 * u.rad,
        epoch = Time('2018-01-01T00:00:00Z', scale='utc', format='isot')
    )
}


if __name__ == "__main__":
    for name, orbit in testObjects.items():
        process(name, orbit, "./")