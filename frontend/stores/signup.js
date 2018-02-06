module.exports = function(state, emitter) {
  state.loginusername = "";
  state.loginpassword = "";
  state.loginis_artist = false;
  state.loginis_curator = false;
  state.loginemail = "";

  emitter.on("updateSignup", ({ key, value }) => {
    state[key] = value;
  });

  emitter.on("updateSignupBool", ({ key, value }) => {
    state[key] = value;
  });
};
