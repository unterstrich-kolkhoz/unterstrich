/* Page: login */

const html = require("choo/html");
const style = require("../lib/style");
const setLocalItem = require("../lib/storage").setLocalItem;

function login(state, emit) {
  return html`
    <body class=${style}>
      <div class="login">
        <h1>login</h1>
        <div id="form-errors" ${state.has_login_error ? "" : "hidden"}>
          Wrong email or password.
        </div>
        <input id="login-username" placeholder="username"
               value=${state.login.username}
               onchange=${update("username")}>
        <input id="login-password" type="password"
               placeholder="password" value=${state.login.password}
               onchange=${update("password")}
               onkeyup=${maybeSubmit}>
        <button onclick=${submitLogin}>Login</button>
      </div>
    </body>
  `;

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
            setLocalItem("username", state.login.username);
            setLocalItem("token", json.token);
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
