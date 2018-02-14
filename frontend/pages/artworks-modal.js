const html = require("choo/html");

module.exports = function(state, emit) {
  if (!state.artworks.showModal) {
    return null;
  }

  return html`
    <div class="modal">
      <div class="modal-content">
        <h3>Create an artwork</h3>
        <select onchange=${update("type")}>
          <option value="image"
                  ${state.artworks.new.type == "image" ? "selected" : ""}>
            Image
          </option>
          <option value="video"
                  ${state.artworks.new.type == "video" ? "selected" : ""}>
            Video
          </option>
        </select>
        <input placeholder="URL"
               value=${state.artworks.new.url}
               onchange=${update("url")}
               required>
        <input type="number"
               min="0"
               step="any"
               placeholder="Price"
               value=${state.artworks.new.price}
               onchange=${update("price")}
               required>
        <button type="submit" onclick=${submit}>
          Submit
        </button>
      </div>
    </div>
  `;

  function update(key) {
    return e => {
      emit("updateNewArtwork", { key, value: e.target.value });
    };
  }

  function submit() {
    emit("createArtwork");
  }
};
