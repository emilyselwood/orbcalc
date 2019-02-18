if ( WEBGL.isWebGLAvailable() === false ) {
    document.body.appendChild( WEBGL.getWebGLErrorMessage() );
}

var container;
var camera, scene, renderer;
var scale = 1000000;
var controls;
var stats;
var sprite;
var raycaster;
var mouse;
var sphere;
var pointsSet = [];
var dragFlag = 0;
var controlState = 0;
var tween;

init();

function init() {
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

    var vrButton = WEBVR.createButton(renderer);
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

    raycaster = new THREE.Raycaster();
    raycaster.params.Points.threshold = 10;

    mouse = new THREE.Vector2();

    THREE.Cache.enabled = true;

    var loader = new THREE.TextureLoader();
    sprite = loader.load( 'img/particle2.png' );

    for (var i = 0; i < 16; i++) {
        loadAsteroidBatch(i);
    }

    var majorPlanets = ["Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune"];
    var colours = [new THREE.Color(1,0,0), new THREE.Color(0,1,0), new THREE.Color(0,0,1),
        new THREE.Color(1,1,0), new THREE.Color(1,0,1), new THREE.Color(0,1,1),
        new THREE.Color(0,1,0), new THREE.Color(1,0,0)];

    for (var i = 0; i < majorPlanets.length; i++) {
        loadMajorPlanet(majorPlanets[i], colours[i])
    }

    createSun();


    // Create ray cast target sphere:
    var sphereGeometry = new THREE.SphereBufferGeometry( 0.1, 12, 12 );
    var sphereMaterial = new THREE.MeshBasicMaterial( { color: 0xff0000 } );
    sphere = new THREE.Mesh( sphereGeometry, sphereMaterial );
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
        new THREE.Color(0.9,0.9,1),
        THREE.Points,
        true
    );
}

function createSun() {
    var positions = [0, 0, 0];
    var colors = [1,1,0];

    var geometry = new THREE.BufferGeometry();
    geometry.addAttribute( 'position', new THREE.Float32BufferAttribute( positions, 3 ) );
    geometry.addAttribute( 'color', new THREE.Float32BufferAttribute( colors, 3 ) );
    geometry.computeBoundingSphere();
    
    var points = new THREE.Points( geometry, new THREE.PointsMaterial( { size: 10, vertexColors: THREE.VertexColors, map: sprite, blending: THREE.AdditiveBlending, depthTest: false, transparent: true } ));
    scene.add( points );
}

function loadData(name, mat, color, T, raytarget) {
    var loader = new THREE.FileLoader();
    
    //load a text file and output the result to the console
    loader.load(
        // resource URL
        name,
        // onLoad callback
        function ( data ) {
            var bits = createGeom(data, color);
            var positions = bits[0];
            var colors = bits[1];
            var ids = bits[2];

            var geometry = new THREE.BufferGeometry();
            geometry.addAttribute( 'position', new THREE.Float32BufferAttribute( positions, 3 ) );
            geometry.addAttribute( 'color', new THREE.Float32BufferAttribute( colors, 3 ) );
            geometry.computeBoundingSphere();

            var points = new T( geometry, mat);
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
            console.log( name + ": " + (xhr.loaded / xhr.total * 100) + '% loaded ' );
        },

        // onError callback
        function ( err ) {
            console.error( 'An error happened loading ' + name );
        }
    );
}

function createGeom(data, color) {
    var positions = [];
    var colors = [];
    var ids = [];
    var lines = data.split(/\r?\n/);
    var n = lines.length;
    for (var i = 0; i < n; i++) {
        if (lines[i] !== "") {
            parts = lines[i].split(",");
            var id = "";
            var x = 0.0;
            var y = 0.0;
            var z = 0.0;

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
            } else {

                x = x / scale;
                y = y / scale;
                z = z / scale;
                positions.push( x, z, y ); // swap z and y around so we get more intuitve controls
                colors.push( color.r, color.g, color.b );
                if (id !== "") {
                    ids.push(id);
                }
            }
        }
    }
    return [positions, colors, ids];
}

function raycastCheck() {
    raycaster.setFromCamera( mouse, camera );
    var intersections = raycaster.intersectObjects(pointsSet );
    intersection = ( intersections.length ) > 0 ? intersections[ 0 ] : null;
    
    if ( intersection !== null) {
        sphere.position.copy( intersection.point );
        sphere.scale.set( 1, 1, 1 );

        var objectID = intersection.object.userData.IDS[intersection.index];
        console.log("clicked on " + objectID );
        var linkTag = document.getElementById("asteroidLink");
        linkTag.innerText = objectID;
        linkTag.href = "https://www.minorplanetcenter.net/db_search/show_object?utf8=✓&object_id=" + objectID
        
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

function onCanvasClick(event) {
    if (dragFlag === 0) {
        raycastCheck();
    }  
    dragFlag = 0
}

function resize() {
    if (needsResize(container)) {
        var w = container.clientWidth;
        var h = container.clientHeight;
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
    var startPos = {
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