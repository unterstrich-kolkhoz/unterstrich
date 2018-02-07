const getLocalItem = require("../lib/storage").getLocalItem;

module.exports = function(state, emitter) {
  state.login = {
    username: "",
    password: "",
    token: getLocalItem("token") || "",
    has_login_error: false
  };

  emitter.on("updateLogin", ({ key, value }) => {
    state.login[key] = value;
  });

  emitter.on("loginError", status_code => {
    state.login.has_login_error = true;
    emitter.emit("render");
  });

  let username = getLocalItem("username");
  let token = getLocalItem("token");
  if (username && token) {
    state.login.username = username;
    state.login.token = token;
    emitter.on("DOMContentLoaded", function() {
      emitter.emit("render");
    });
  }
};
