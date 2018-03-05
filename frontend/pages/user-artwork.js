const html = require("choo/html");

const { defaultVertex, defaultFragment, runShader } = require("../lib/shader");
const style = require("../lib/style");

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

  function zoomedInModal() {
    if (!state.artworks.showZoom) return null;
    return html`
      <div class="user-artwork-container-zoomed"
           onclick=${hideArtworkZoom}>
        ${artworkRender()}
      </div>
    `;
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
      case "fragment-shader":
        fetch(artwork.url, {
          method: "GET",
          mode: "cors"
        })
          .then(res => res.text())
          .then(shader => runShader(defaultVertex, shader));
        return html`<div id="glsl-container" onclick=${onclick}></div>`;
      case "vertex-shader":
        fetch(artwork.url, {
          method: "GET",
          mode: "cors"
        })
          .then(res => res.text())
          .then(shader => runShader(shader, defaultFragment));
        return html`<div id="glsl-container" onclick=${onclick}></div>`;
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
