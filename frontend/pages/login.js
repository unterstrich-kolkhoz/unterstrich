/* Page: login */

const html = require("choo/html");

const style = require("../lib/style");

function login(state, emit) {
  return html`
    <body class=${style}>
      <div class="login">
        <h1>login</h1>
        <div id="form-errors" ${state.has_login_error ? "" : "hidden"}>
          Wrong email or password.
        </div>
        <input id="login-username" placeholder="username"
               value=${state.username} onchange=${updateUser}>
        <input id="login-password" type="password"
               placeholder="password" value=${state.password}
               onchange=${updatePass} onkeypress=${maybeSubmit}>
        <button onclick=${submitLogin}>Login</button>
      </div>
    </body>
  `;

  function updatePass(e) {
    emit("updatePass", { value: e.target.value });
  }

  function maybeSubmit(e) {
    if (e.keyCode == 13) {
      submitLogin();
    }
  }

  function updateUser(e) {
    emit("updateName", { value: e.target.value });
  }

  function submitLogin() {
    fetch("/login", {
      method: "POST",
      body: JSON.stringify({
        username: state.username,
        password: state.password
      })
    })
      .then(res => {
        if (res.status == 200) {
          res.json().then(json => {
            emit("updatePass", { value: "" });
            emit("updateToken", json.token);
            emit("pushState", `/${state.username}`);
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
