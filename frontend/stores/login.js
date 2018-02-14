const { getLocal } = require("../lib/storage");

module.exports = function(state, emitter) {
  state.login = {
    username: getLocal("username") || "",
    password: "",
    token: getLocal("token") || "",
    has_login_error: false
  };

  emitter.on("updateLogin", ({ key, value }) => {
    state.login[key] = value;
  });

  emitter.on("loginError", status_code => {
    state.login.has_login_error = true;
    emitter.emit("render");
  });

  if (state.login.username && state.login.token) {
    emitter.on("DOMContentLoaded", function() {
      emitter.emit("render");
    });
  }
};
