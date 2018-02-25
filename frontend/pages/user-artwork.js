const html = require("choo/html");

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
        <img src="${artwork.url}" onclick=${showArtworkZoom}>
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

  return html`
    <body class="${style}">
      ${zoomedInModal()}
      <div class="user-artwork">
        <div class="user-artwork-container">
          <img src="${artwork.url}" onclick=${showArtworkZoom}>
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
