/* Page: login */

const html = require("choo/html");
const page = require("../lib/page");
const { setLocal } = require("../lib/storage");

function login(state, emit) {
  return page(html`
      <div class="login">
        <h1>login</h1>
        <div id="form-errors" ${state.has_login_error ? "" : "hidden"}>
          Wrong email or password.
        </div>
        <input id="login-username" placeholder="username"
               value=${state.login.username}
               required
               onchange=${update("username")}>
        <input id="login-password" type="password"
               placeholder="password" value=${state.login.password}
               onchange=${update("password")}
               required
               onkeyup=${maybeSubmit}>
        <button type="submit"
                onclick=${submitLogin}>
            Login
        </button>
      </div>
  `);

  function maybeSubmit(e) {
    if (e.keyCode == 13) {
      submitLogin();
    }
    emit("updateLogin", { key: "password", value: e.target.value });
  }

  function update(key) {
    return e => emit("updateLogin", { key, value: e.target.value });
  }

  function submitLogin() {
    fetch("/login", {
      method: "POST",
      body: JSON.stringify({
        username: state.login.username,
        password: state.login.password
      })
    })
      .then(res => {
        if (res.status == 200) {
          res.json().then(json => {
            emit("updateLogin", { key: "password", value: "" });
            emit("updateLogin", { key: "token", value: json.token });
            setLocal("username", state.login.username);
            setLocal("token", json.token);
            emit("pushState", `/${state.login.username}`);
          });
        } else {
          emit("loginError", res.status);
        }
      })
      .catch(res => {
        emit("loginError", res.status);
      });
  }
}

module.exports = login;
