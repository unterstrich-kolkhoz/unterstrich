const authFetch = require("../lib/fetch");
const authAxios = require("../lib/axios");

module.exports = function(state, emitter) {
  state.artworks = {
    user: "",
    pending: false,
    showModal: false,
    showZoom: false,
    artworks: [],
    progress: 0.0,
    new: {
      name: "",
      description: "",
      type: "image",
      price: 0.0
    }
  };

  emitter.on("getArtworks", username => {
    state.artworks.user = username;
    state.artworks.pending = true;

    // TODO: /username/artworks
    authFetch(state, emitter, "/artworks/", {
      method: "GET"
    }).then(res => {
      if (res.status == 200) {
        res.json().then(json => {
          state.artworks.pending = false;
          state.artworks.artworks = json;
          emitter.emit("render");
        });
      }
    });
  });

  emitter.on("showArtworkModal", show => {
    state.artworks.showModal = show;
    emitter.emit("render");
  });

  emitter.on("showArtworkZoom", show => {
    state.artworks.showZoom = show;
    emitter.emit("render");
  });

  emitter.on("updateNewArtwork", ({ key, value }) => {
    state.artworks.new[key] = value;
  });

  function uploadFile(artwork_id, file) {
    let data = new FormData();
    data.append("upload", file);
    authAxios(state, emitter, {
      method: "post",
      url: `/artworks/${artwork_id}/upload`,
      data: data,
      onUploadProgress: function(e) {
        state.artworks.progress = e.loaded / e.total;
        emitter.emit("render");
      }
    }).then(res => {
      if (res.status == 200) emitter.emit("render");
      state.artworks.showModal = false;
      state.artworks.user = "";
      state.artworks.pending = false;
    });
  }

  emitter.on("createArtwork", ({ file }) => {
    state.artworks.pending = true;
    state.artworks.new.price = parseFloat(state.artworks.new.price);
    authFetch(state, emitter, "/artworks/", {
      method: "POST",
      body: JSON.stringify(state.artworks.new)
    }).then(res => {
      if (res.status == 200) {
        res.json().then(json => uploadFile(json.id, file));
      }
    });
  });
};
