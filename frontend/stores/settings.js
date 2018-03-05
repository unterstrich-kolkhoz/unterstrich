const authFetch = require("../lib/fetch");

module.exports = function(state, emitter) {
  state.settings = {
    tab: "info",

    info: {
      id: 0,
      email: "",
      username: "",
      name: "",

      line1: "",
      line2: "",
      zip: "",
      city: "",
      state: "",
      country: "",

      ello: "",
      github: "",
      website: ""
    }
  };

  emitter.on("settingsTab", tab => {
    state.settings.tab = tab;
    emitter.emit("render");
  });

  emitter.on("fetchSettings", () => {
    authFetch(state, emitter, "/me", {
      method: "GET"
    }).then(res => {
      if (res.status == 200) {
        res.json().then(json => {
          for (let key in state.settings.info) {
            state.settings.info[key] = json[key];
          }
          emitter.emit("render");
        });
      }
    });
  });

  emitter.on("submitSettings", () => {
    authFetch(state, emitter, `/users/${state.settings.info.id}`, {
      method: "PUT",
      body: JSON.stringify(state.settings.info)
    }).then(res => {
      if (res.status == 200) {
        res.json().then(json => {
          emitter.emit("fetchSettings");
          emitter.emit("render");
        });
      }
    });
  });

  emitter.on("updateSettings", ({ key, value }) => {
    state.settings.info[key] = value;
  });

  emitter.emit("fetchSettings");
};
