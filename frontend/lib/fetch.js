module.exports = (state, emitter, url, details) => {
  details.headers = new Headers(
    Object.assign({}, details.headers, {
      Authorization: `Bearer ${state.login.token}`
    })
  );
  return fetch(url, details).then(res => {
    if (res.status == 401) {
      state.login.token = "";
      emitter.emit("pushState", "/login");
      emitter.emit("render");
    }
    return res;
  });
};
