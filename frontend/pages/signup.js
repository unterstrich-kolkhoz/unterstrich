/* Page: signup */

const html = require("choo/html");

const style = require("../lib/style");

function signup(state, emit) {
  return html`
    <body class=${style}>
      <div class="signup">
        <h1>signup</h1>
        <input id="signup-email" placeholder="email"
               value=${state.loginemail}
               onchange=${update("loginemail")}>
        <input id="signup-username" placeholder="username"
               value=${state.loginusername}
               onchange=${update("loginusername")}>
        <input id="signup-password" type="password"
               placeholder="password" value=${state.loginpassword}
               onchange=${update("loginpassword")}>
        <label for="signup-artist">Artist?</label>
        <input id="signup-artist" type="checkbox"
               value=${state.loginis_artist}
               onchange=${updateBool("loginis_artist")}>
        <label for="signup-curator">Curator?</label>
        <input id="signup-curator" type="checkbox"
               value=${state.loginis_curator}
               onchange=${updateBool("loginis_curator")}>
        <button onclick=${submitSignup}>Sign Up</button>
      </div>
    </body>
  `;

  function update(key) {
    return e =>
      emit("updateSignup", {
        key,
        value: e.target.value
      });
  }

  function updateBool(key) {
    return e =>
      emit("updateSignupBool", {
        key,
        value: e.target.checked
      });
  }

  function submitSignup() {
    fetch("/users", {
      method: "POST",
      body: JSON.stringify({
        username: state.loginusername,
        password: state.loginpassword,
        email: state.loginemail,
        is_artist: state.loginis_artist,
        is_curator: state.loginis_curator
      })
    })
      .then(res => {
        if (res.status == 200) {
          res.json().then(json => {
            emit("update", { key: "password", value: "" });
            emit("pushState", `/`);
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

module.exports = signup;
