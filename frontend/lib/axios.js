const axios = require("axios");

module.exports = (state, emitter, details) => {
  details.headers = Object.assign({}, details.headers, {
    Authorization: `Bearer ${state.login.token}`
  });
  return axios(details).then(res => {
    if (res.status == 401) {
      state.login.token = "";
      emitter.emit("pushState", "/login");
      emitter.emit("render");
    }
    return res;
  });
};
