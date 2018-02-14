const authFetch = require("../lib/fetch");

module.exports = function(state, emitter) {
  emitter.on("getUserInfo", () => {
    authFetch(state, emitter, "/me", {
      method: "GET"
    }).then(res => {
      if (res.status == 200) {
        res.json().then(json => {
          state.userInfo = json;
          emitter.emit("render");
        });
      }
    });
  });
};
