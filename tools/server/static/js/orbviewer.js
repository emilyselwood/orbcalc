if ( WEBGL.isWebGLAvailable() === false ) {
    document.body.appendChild( WEBGL.getWebGLErrorMessage() );
}

let container;
let camera, scene, renderer;
const scale = 1000000;
let controls;
let stats;
let sprite;
let raycaster;
let mouse;
let sphere;
let pointsSet = [];
let dragFlag = 0;
let controlState = 0;
let tween;
let orbitLine = null;
let server = true;

// size of the solar system
// TODO: pull this from somewhere and make the server generate it.
const minX = -1.0993153024260256e+10 / scale;
const maxX = 1.1259105381765476e+10 / scale;
const minY = -8.336972753734525e+09 / scale;
const maxY = 1.1216725295000006e+10 / scale;
const minZ = -5.482463379824468e+09 / scale;
const maxZ = 4.381383003697839e+09 / scale;

// color pallets for pride flags.
const Rainbow = [
    [231 / 255, 0, 0],
    [255 / 255, 140 / 255, 0],
    [255 / 255, 239 / 255, 0],
    [0, 129 / 255, 31 / 255],
    [0, 68 / 255, 255 / 255]
];

const Trans = [
    [85 / 255, 205 / 255, 252 / 255],
    [247 / 255, 168 / 255, 184 / 255],
    [255 / 255, 255 / 255, 255 / 255],
    [247 / 255, 168 / 255, 184 / 255],
    [85 / 255, 205 / 255, 252 / 255]
];

// TODO: check that this is correct
const Bi = [
    [217 / 255, 0, 111 / 255],
    [217 / 255, 0, 111 / 255],
    [116 / 255, 77 / 255, 152 / 255],
    [0, 51 / 255, 171 / 255],
    [0, 51 / 255, 171 / 255]
];

// TODO: Lesbian
// TODO: Pan

// basic color pallet for asteroids.
const White = [[1, 1, 1]];

const colourmaps = [
    Rainbow,
    Trans,
    Bi,
    White
];
// pick a random colour map for the asteroids
let colourMap = null;

init();

function init() {

    // pick a colour pallet for the asteroids.
    const url = new URL(window.location.href);
    if (url.searchParams.has('f')) {
        const f = parseInt(url.searchParams.get('f'), 10);
        if(!Number.isNaN(f) && f >= 0 && f < colourmaps.length) {
            colourMap = colourmaps[f];
        }
    }
    if (colourMap === null) {
        colourMap = colourmaps[Math.floor(Math.random() * colourmaps.length)];
    }

    document.addEventListener( 'mousemove', onDocumentMouseMove, false );
    container = document.getElementById( 'container' );
    container.addEventListener('click', onCanvasClick, false);
    container.addEventListener("mousedown", function(){
        dragFlag = 0;
    }, false);

    renderer = new THREE.WebGLRenderer( { antialias: true, logarithmicDepthBuffer: true } );

    renderer.vr.enabled = true;
    // HACK because vr does not play well with orbit controls
    if (!renderer.vr._origGetCamera) renderer.vr._origGetCamera = renderer.vr.getCamera;

    renderer.setPixelRatio( window.devicePixelRatio );
    renderer.setSize( window.innerWidth, window.innerHeight );
    container.appendChild( renderer.domElement );

    stats = new Stats();
    container.appendChild( stats.dom );

    let vrButton = WEBVR.createButton(renderer);
    if (vrButton) {
        document.body.appendChild(vrButton);
    }

    camera = new THREE.PerspectiveCamera( 90, window.innerWidth / window.innerHeight, 3, 100000 );
    camera.position.z = 60;

    window.addEventListener('vrdisplaypresentchange', () => {
        camera.position.z = 60;
    });

    controls = new THREE.OrbitControls( camera );

    scene = new THREE.Scene();
    scene.background = new THREE.Color( 0x000005 );
    scene.fog = new THREE.Fog( 0x000005, 90000, 100000 );
    camera.lookAt(scene.position);

    raycaster = new THREE.Raycaster();
    raycaster.params.Points.threshold = 10;

    mouse = new THREE.Vector2();

    THREE.Cache.enabled = true;

    sprite = new THREE.TextureLoader().load( 'img/particle2.png' );

    for (let i = 0; i < 16; i++) {
        loadAsteroidBatch(i);
    }

    const majorPlanets = ["Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune"];
    const colours = [new THREE.Color(1,1,0), new THREE.Color(0,1,0), new THREE.Color(0,0,1),
        new THREE.Color(1,0,0), new THREE.Color(1,0,1), new THREE.Color(0,1,1),
        new THREE.Color(0,1,0), new THREE.Color(1,0,0)];

    for (let i = 0; i < majorPlanets.length; i++) {
        loadMajorPlanet(majorPlanets[i], colours[i])
    }

    createSun();

    // Create ray cast target sphere:
    sphere = new THREE.Mesh(
        new THREE.SphereBufferGeometry( 0.1, 12, 12 ),
        new THREE.MeshBasicMaterial( { color: 0xff0000 } )
    );
    scene.add( sphere );

    animate();
}

function loadMajorPlanet(name, color) {
    loadData("data/" +name+ ".csv", 
    new THREE.LineBasicMaterial( { vertexColors: THREE.VertexColors } ),
    color,
    THREE.Line,
    false);
}

function loadAsteroidBatch(batch) {
    loadData(
        "data/data-" + batch + ".csv",
        new THREE.PointsMaterial( {
            size: 6,
            vertexColors: THREE.VertexColors,
            map: sprite,
            blending: THREE.AdditiveBlending,
            depthTest: true,
            transparent: false,
            alphaTest: 0.5,
            fog: false,
            lights: false,
            sizeAttenuation: false
        } ),
        colourMap,
        THREE.Points,
        true
    );
}

function createSun() {
    const positions = [0, 0, 0];
    const colors = [1,1,0];

    const geometry = new THREE.BufferGeometry();
    geometry.addAttribute( 'position', new THREE.Float32BufferAttribute( positions, 3 ) );
    geometry.addAttribute( 'color', new THREE.Float32BufferAttribute( colors, 3 ) );
    geometry.computeBoundingSphere();
    
    const points = new THREE.Points( geometry, new THREE.PointsMaterial( { size: 10, vertexColors: THREE.VertexColors, map: sprite, blending: THREE.AdditiveBlending, depthTest: false, transparent: true } ));
    scene.add( points );
}

function loadData(name, mat, color, T, raytarget) {
    let loader = new THREE.FileLoader();
    
    //load a text file and output the result to the console
    loader.load(
        // resource URL
        name,
        // onLoad callback
        function ( data ) {
            const bits = createGeom(data, color);
            const positions = bits[0];
            const colors = bits[1];
            const ids = bits[2];

            const geometry = new THREE.BufferGeometry();
            geometry.addAttribute( 'position', new THREE.Float32BufferAttribute( positions, 3 ) );
            geometry.addAttribute( 'color', new THREE.Float32BufferAttribute( colors, 3 ) );
            geometry.computeBoundingSphere();

            const points = new T( geometry, mat);
            if (ids.length > 0) {
                points.userData = {IDS: ids};
            }
            if(raytarget){
                pointsSet.push(points);
            }
            scene.add( points );
            console.log( "loaded " + name)
        },

        // onProgress callback
        function ( xhr ) {
            //c onsole.log( name + ": " + (xhr.loaded / xhr.total * 100) + '% loaded ' );
        },

        // onError callback
        function ( err ) {
            console.error( 'An error happened loading ' + name + ' ' + err);
        }
    );
}

function createGeom(data, color) {
    let positions = [];
    let colors = [];
    let ids = [];
    const lines = data.split(/\r?\n/);
    const n = lines.length;
    for (let i = 0; i < n; i++) {
        if (lines[i] !== "") {
            let parts = lines[i].split(",");
            let id = "";
            let x = 0.0;
            let y = 0.0;
            let z = 0.0;

            if (parts.length === 4) {
                id = parts[0];
                x = parseFloat(parts[1]);
                y = parseFloat(parts[2]);
                z = parseFloat(parts[3]);
            } else if (parts.length === 3) {
                x = parseFloat(parts[0]);
                y = parseFloat(parts[1]);
                z = parseFloat(parts[2]);
            } else {
                console.log("could not decode line " + i);
                continue;
            }
            
            if (isNaN(x) || isNaN(y) || isNaN(z)) {
                console.log("could not decode " + lines[i] + " line " + i);
                continue;
            }

            x = x / scale;
            y = y / scale;
            z = z / scale;
            positions.push( x, z, y ); // swap z and y around so we get more intuitive controls
            if (Array.isArray(color)) { // if our color is an array we probably need to pick one of the entries.
                let c = mapToColour(color, x, y, z);
                colors.push(c[0], c[1], c[2]);
            } else { // otherwise it should be a three.Color object.
                colors.push(color.r, color.g, color.b);
            }
            if (id !== "") {
                ids.push(id);
            }

        }
    }
    return [positions, colors, ids];
}

function raycastCheck() {
    raycaster.setFromCamera( mouse, camera );
    const intersections = raycaster.intersectObjects(pointsSet );
    const intersection = ( intersections.length ) > 0 ? intersections[ 0 ] : null;
    
    if ( intersection !== null) {
        sphere.position.copy( intersection.point );
        sphere.scale.set( 1, 1, 1 );

        const objectID = intersection.object.userData.IDS[intersection.index];
        console.log("clicked on " + objectID );
        const linkTag = document.getElementById("asteroidLink");
        linkTag.innerText = objectID;
        linkTag.href = "https://www.minorplanetcenter.net/db_search/show_object?utf8=✓&object_id=" + objectID;
        if (server) {
            const loader = new THREE.FileLoader();
            loader.load("/obj/" + objectID.replace(/ /g, '+'),
                function (data) {
                    if (orbitLine !== null) {
                        scene.remove(orbitLine)
                    }

                    const response = JSON.parse(data);

                    let positions = [];
                    let colors = [];
                    for (let i = 0; i < response.Orbit.length; i++) {
                        positions.push(response.Orbit[i].X / scale, response.Orbit[i].Z / scale, response.Orbit[i].Y / scale);
                        colors.push(0, 0, 255);
                    }

                    positions.push(positions[0], positions[1], positions[2]);
                    colors.push(0, 0, 255);

                    const geometry = new THREE.BufferGeometry();
                    geometry.addAttribute('position', new THREE.Float32BufferAttribute(positions, 3));
                    geometry.addAttribute('color', new THREE.Float32BufferAttribute(colors, 3));
                    geometry.computeBoundingSphere();

                    orbitLine = new THREE.Line(geometry, new THREE.LineBasicMaterial({color: 0x77FF77, linewidth: 10}));
                    scene.add(orbitLine);
                    console.log("loaded " + objectID);
                },
                function(xhr) {

                },
                function(error) {
                    // If we cant talk to the server for some reason we won't try again.
                    // It is likely we are running in a static environment.
                    server = false;
                }
            )
        }
    }
}

function needsResize(canvas) {
    if (canvas.lastWidth !== canvas.clientWidth || canvas.lastHeight !== canvas.clientHeight) {
        canvas.width = canvas.lastWidth = canvas.clientWidth;
        canvas.height = canvas.lastHeight = canvas.clientHeight;
        return true;
    }
}

function onDocumentMouseMove( event ) {
    event.preventDefault();
    mouse.x = ( event.clientX / window.innerWidth ) * 2 - 1;
    mouse.y = - ( event.clientY / window.innerHeight ) * 2 + 1;
    dragFlag = 1;                
}

function onCanvasClick() {
    if (dragFlag === 0) {
        raycastCheck();
    }  
    dragFlag = 0
}

function resize() {
    if (needsResize(container)) {
        const w = container.clientWidth;
        const h = container.clientHeight;
        camera.aspect = w / h;
        camera.updateProjectionMatrix();
        renderer.setSize(w, h, false);
    }
}

function animate(time) {
    if (!renderer.domElement.parentElement) {
        return;
    }
    requestAnimationFrame(animate);
    resize();
    // if vr is enabled three will handle the controls for us.
    if (renderer.vr.isPresenting()) {
        renderer.vr.getCamera = renderer.vr._origGetCamera;
    } else {
        if (controlState === 0) {
            controls.update();
        } else {
            tween.update(time);
        }
        renderer.vr.getCamera = () => camera;
    }

    stats.update();
    render();
}

function render() {
    renderer.render( scene, camera );
}

function toggleTour() {
    if (controlState === 0) {
        controlState = 1;
        setupMove();
    } else {
        controlState = 0;
        tween.stop();
        controls.update();
    }

}

function setupMove() {
    let targetPos = pickPosition();
    let startPos = {
        x: camera.position.x,
        y: camera.position.y,
        z: camera.position.z
    };

    console.log("moving from: " + startPos.x + "," + startPos.y + "," +startPos.z + " to: " + targetPos.x + "," + targetPos.y + "," +targetPos.z );

    tween = new TWEEN.Tween(startPos)
        .to(targetPos, 10000)
        .easing(TWEEN.Easing.Quadratic.Out)
        .onUpdate(function() {
            camera.position.set(startPos.x, startPos.y, startPos.z);
            camera.lookAt(scene.position);
            controls.update();
        })
        .onComplete(setupMove)
        .start();

}

function pickPosition() {
    // pick a number between +-90 normal distribution
    let lat = (randomNumber(45) + randomNumber(45)) - 90;

    // pick a number between +-180 liner distribution
    let lon = randomNumber(360) - 180;

    // convert these two angles to a point on a sphere somewhere near the edge of the solar system.
    const R = 1000;
    return {
        x: R * Math.cos(lat) * Math.cos(lon),
        y: R * Math.cos(lat) * Math.sin(lon),
        z: R * Math.sin(lat)
    };
}

function randomNumber(max) {
    return Math.floor(Math.random() * max)
}

function mapToColour(map, x, y, z) {
    let bucket = Math.round(((x - minX) / (maxX - minX)) * (map.length-1));
    if (bucket < 0 || bucket >= map.length) {
        console.log("could not find bucket for x: " + x + " got bucket " + bucket);
        return map[0]
    }
    return map[bucket];
}
