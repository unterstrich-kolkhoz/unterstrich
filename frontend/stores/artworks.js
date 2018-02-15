const authFetch = require("../lib/fetch");

module.exports = function(state, emitter) {
  state.artworks = {
    user: "",
    pending: false,
    showModal: false,
    showZoom: false,
    artworks: [],
    new: {
      name: "",
      description: "",
      url: "",
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

  emitter.on("createArtwork", () => {
    state.artworks.new.price = parseFloat(state.artworks.new.price);
    console.log(JSON.stringify(state.artworks.new));
    authFetch(state, emitter, "/artworks/", {
      method: "POST",
      body: JSON.stringify(state.artworks.new)
    }).then(res => {
      if (res.status == 200) {
        state.artworks.showModal = false;
        state.artworks.user = "";
        emitter.emit("render");
      }
    });
  });
};
