module.exports = function(state, emitter) {
  state.signup = {
    username: "",
    password: "",
    is_artist: false,
    is_curator: false,
    email: ""
  };

  emitter.on("updateSignup", ({ key, value }) => {
    state.signup[key] = value;
  });

  emitter.on("updateSignupBool", ({ key, value }) => {
    state.signup[key] = value;
  });
};
