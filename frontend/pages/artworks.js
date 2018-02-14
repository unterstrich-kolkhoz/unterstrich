const html = require("choo/html");

module.exports = function(state, emit, username) {
  if (state.artworks.user != username) {
    emit("getArtworks", username);
  }

  return html`
    <div class="artwork-list">
      ${state.artworks.artworks.map(renderArtwork)}
    </div>
  `;

  function renderArtwork(artwork) {
    if (!artwork.url) return null;
    return html`
      <div class="artwork">
        <a class="img-link" href="/${username}/artworks/${artwork.id}">
          <img src="${artwork.thumbnail || artwork.url}"
               onerror="${e => (e.target.style.display = "none")}">
        </a>
      </div>
    `;
  }
};
