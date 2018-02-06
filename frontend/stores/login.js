module.exports = function(state, emitter) {
  state.username = "";
  state.password = "";
  state.token = "";
  state.has_login_error = false;

  emitter.on("updatePass", ({ value }) => {
    state.password = value;
  });

  emitter.on("updateName", ({ value }) => {
    state.username = value;
  });

  emitter.on("updateToken", ({ value }) => {
    state.token = value;
  });

  emitter.on("loginError", status_code => {
    state.has_login_error = true;
    emitter.emit("render");
  });
};
