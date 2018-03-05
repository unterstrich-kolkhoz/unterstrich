module.exports = function(state, emitter) {
  state.signup = {
    username: "",
    password: "",
    email: ""
  };

  emitter.on("updateSignup", ({ key, value }) => {
    state.signup[key] = value;
  });

  emitter.on("updateSignupBool", ({ key, value }) => {
    state.signup[key] = value;
  });
};
