const html = require("choo/html");

module.exports = function(state, emit) {
  if (!state.artworks.showModal) {
    return null;
  }

  function progress() {
    if (!state.artworks.pending) return null;

    return html`
      <progress value="${100 * state.artworks.progress}" max="100">
      </progress>
    `;
  }

  function accept() {
    switch (state.artworks.new.type) {
      case "image":
        return "image/*";
      case "video":
        return "video/*";
      case "shader":
        return ".glsl,.frag,.vert";
      default:
        return "";
    }
  }

  return html`
    <div class="modal" onclick=${disable}>
      <div class="modal-content" onclick=${silence}>
        <h3>Create an artwork</h3>
        ${progress()}
        <input placeholder="Name"
               value=${state.artworks.new.name}
               onchange=${update("name")}
               required>
        <textarea onchange=${update("description")}
                  placehold="description"
                  rows="4" cols="30"
                  required>
          ${state.artworks.new.name}
        </textarea>
        <select onchange=${updateType}>
          <option value="image"
                  ${state.artworks.new.type == "image" ? "selected" : ""}>
            Image
          </option>
          <option value="video"
                  ${state.artworks.new.type == "video" ? "selected" : ""}>
            Video
          </option>
          <option value="shader"
                  ${state.artworks.new.type == "shader" ? "selected" : ""}>
            Shader
          </option>
        </select>
        <input id="upload-file"
               type="file"
               accept=${accept()}
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
        <button onclick=${disable}>
          Cancel
        </button>
      </div>
    </div>
  `;

  function disable() {
    emit("showArtworkModal", false);
  }

  function silence(e) {
    //e.preventDefault();
    e.stopPropagation();
  }

  function update(key) {
    return e => {
      emit("updateNewArtwork", { key, value: e.target.value });
    };
  }

  function updateType(e) {
    update("type")(e);
    emit("render");
  }

  function submit() {
    const file = document.getElementById("upload-file").files[0];
    emit("createArtwork", { file });
  }
};
