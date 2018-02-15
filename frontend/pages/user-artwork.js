const html = require("choo/html");

const style = require("../lib/style");

module.exports = function(state, emit) {
  const { username, id } = state.params;

  if (state.artworks.user != username) {
    emit("getArtworks", username);
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
        <img src="${artwork.url}" onclick=${showArtworkZoom}>
      </div>
    `;
  }

  return html`
    <body class="${style}">
      ${zoomedInModal()}
      <div class="user-artwork">
        <div class="user-artwork-container">
          <img src="${artwork.url}" onclick=${showArtworkZoom}>
        </div>
        <div class="tombstone">
          <h3>
            ${artwork.name} by <a href="/${username}">${username}</a>
          </h3>
          <p>${artwork.description}</p>
          <p class="artwork-price">Price: $${artwork.price.toFixed(2)}</p>
        </div>
      </div>
    </body>
  `;
};
