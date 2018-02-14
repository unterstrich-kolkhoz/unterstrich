const authFetch = require("../lib/fetch");

module.exports = function(state, emitter) {
  state.artworks = {
    user: "",
    showModal: false,
    artworks: [],
    new: {
      url: "",
      type: "image",
      price: 0.0
    }
  };

  emitter.on("getArtworks", username => {
    state.artworks.user = username;

    // TODO: /username/artworks
    authFetch(state, emitter, "/artworks/", {
      method: "GET"
    }).then(res => {
      if (res.status == 200) {
        res.json().then(json => {
          console.log(json);
          state.artworks.artworks = json;
          emitter.emit("render");
        });
      }
    });
  });

  emitter.on("showArtworkModal", show => {
    state.artworks.showModal = show;
  });

  emitter.on("updateNewArtwork", ({ key, value }) => {
    state.artworks.new[key] = value;
  });

  emitter.on("createArtwork", () => {
    console.log(state.artworks.new);
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
