const html = require("choo/html");

const style = require("../lib/style");

const THREE = require("three");

module.exports = function(state, emit) {
  const { username, id } = state.params;

  if (state.artworks.user != username) {
    emit("getArtworks", username);
    return null;
  }

  if (!state.userInfo) {
    emit("getUserInfo");
    return null;
  }

  if (state.artworks.pending) return null;

  const artworks = state.artworks.artworks.filter(a => a.id == id);

  if (!artworks.length) {
    emit("pushState", "/");
    emit("render");
    return null;
  }

  const artwork = artworks[0];
  emit("DOMTitleChange", `_ | ${artwork.name} by ${username}`);

  function showArtworkZoom() {
    emit("showArtworkZoom", true);
  }

  function hideArtworkZoom() {
    emit("showArtworkZoom", false);
  }

  function star() {
    emit("star", artwork.id);
  }

  function unstar() {
    emit("unstar", artwork.id);
  }

  function zoomedInModal() {
    if (!state.artworks.showZoom) return null;
    return html`
      <div class="user-artwork-container-zoomed"
           onclick=${hideArtworkZoom}>
        ${artworkRender()}
      </div>
    `;
  }

  function starButton() {
    if (
      artwork.stars &&
      artwork.stars.filter(u => u.id == state.userInfo.id).length > 0
    ) {
      return html`
        <button onclick=${unstar}>Unstar</button>
      `;
    }
    return html`
      <button onclick=${star}>Star</button>
    `;
  }

  function runShader(url) {
    fetch(url, {
      method: "GET",
      mode: "cors"
    })
      .then(res => {
        return res.text();
      })
      .then(fragment => {
        let container;
        let camera, scene, renderer;
        let uniforms;

        init();
        animate();

        function init() {
          container = document.getElementById("glsl-container");
          container.innerHTML = "";

          camera = new THREE.Camera();
          camera.position.z = 1;

          scene = new THREE.Scene();

          let geometry = new THREE.PlaneBufferGeometry(2, 2);

          uniforms = {
            u_time: { type: "f", value: 1.0 },
            u_resolution: { type: "v2", value: new THREE.Vector2() },
            u_mouse: { type: "v2", value: new THREE.Vector2() }
          };

          let material = new THREE.ShaderMaterial({
            uniforms: uniforms,
            vertexShader: "void main() {gl_Position = vec4(position, 1.0);}",
            fragmentShader: fragment
          });

          let mesh = new THREE.Mesh(geometry, material);
          scene.add(mesh);

          renderer = new THREE.WebGLRenderer();
          renderer.setPixelRatio(window.devicePixelRatio);

          container.appendChild(renderer.domElement);

          onWindowResize();
          window.addEventListener("resize", onWindowResize, false);

          document.onmousemove = function(e) {
            uniforms.u_mouse.value.x = e.pageX;
            uniforms.u_mouse.value.y = e.pageY;
          };
        }

        function onWindowResize() {
          const style = getComputedStyle(container);
          renderer.setSize(parseFloat(style.width), parseFloat(style.height));
          uniforms.u_resolution.value.x = renderer.domElement.width;
          uniforms.u_resolution.value.y = renderer.domElement.height;
        }

        function animate() {
          requestAnimationFrame(animate);
          render();
        }

        function render() {
          uniforms.u_time.value += 0.05;
          renderer.render(scene, camera);
        }
      });
  }

  function artworkRender(onclick) {
    switch (artwork.type) {
      case "image":
        return html`<img src="${artwork.url}" onclick=${onclick}>`;
      case "video":
        return html`
          <video controls onclick=${onclick}>
            <source src="${artwork.url}">
          </video>`;
      case "shader":
        runShader(artwork.url);
        return html`<div id="glsl-container"></div>`;
      default:
        return null;
    }
  }

  return html`
    <body class="${style}">
      ${zoomedInModal()}
      <div class="user-artwork">
        <div class="user-artwork-container">
          ${artworkRender(showArtworkZoom)}
        </div>
        <div class="right">
          ${starButton()}
        </div>
        <div class="tombstone">
          <h3>
            ${artwork.name} by <a href="/${username}">${username}</a>
            <div class="right">
              <span class="label">${artwork.views} views</span>
              <span class="label">${
                artwork.stars ? artwork.stars.length : 0
              } stars</span>
            </div>
          </h3>
          <p>${artwork.description}</p>
          <p class="artwork-price">Price: $${artwork.price.toFixed(2)}</p>
        </div>
      </div>
    </body>
  `;
};
