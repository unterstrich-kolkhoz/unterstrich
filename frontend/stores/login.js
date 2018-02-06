const getLocalItem = require("../lib/storage").getLocalItem;

module.exports = function(state, emitter) {
  state.username = "";
  state.password = "";
  state.token = getLocalItem("token") || "";
  state.has_login_error = false;

  emitter.on("updateLogin", ({ key, value }) => {
    console.log(key);
    console.log(value);
    state[key] = value;
  });

  emitter.on("loginError", status_code => {
    state.has_login_error = true;
    emitter.emit("render");
  });

  let username = getLocalItem("username");
  let token = getLocalItem("token");
  if (username && token) {
    state.username = username;
    console.log(state.username);
    state.token = token;
    emitter.on("DOMContentLoaded", function() {
      emitter.emit("render");
    });
  }
};
