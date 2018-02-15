const html = require("choo/html");

module.exports = function(state, emit, username) {
  if (state.artworks.user != username && state.artworks.pending == false) {
    emit("getArtworks", username);
  }

  return html`
    <div class="artwork-list">
      ${state.artworks.artworks.map(renderArtwork)}
    </div>
  `;

  function renderArtwork(artwork) {
    if (!artwork.thumbnail)
      return html`<div class="artwork-placeholder">Thumbnail not yet generated</div>`;
    return html`
      <div class="artwork">
        <a class="img-link" href="a/${username}/${artwork.id}">
          <img src="${artwork.thumbnail}"
               onerror="${e => (e.target.style.display = "none")}">
        </a>
      </div>
    `;
  }
};
