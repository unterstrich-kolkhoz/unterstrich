/* Page: signup */

const html = require("choo/html");

const page = require("../lib/page");

function signup(state, emit) {
  emit("DOMTitleChange", "_ | signup");
  if (state.login.token) {
    emit("pushState", state.login.username);
    emit("render");
  }

  return page(html`
    <div class="signup">
      <h1>signup</h1>
      <input id="signup-email" placeholder="email"
             value=${state.signup.email}
             required
             onchange=${update("email")}>
      <input id="signup-username" placeholder="username"
             value=${state.signup.username}
             required
             onchange=${update("username")}>
      <input id="signup-password" type="password"
             placeholder="password" value=${state.signup.password}
             required
             onchange=${update("password")}>
      <button type="submit"
              onclick=${submitSignup}>
        Sign Up
      </button>
    </div>
  `);

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
        username: state.signup.username,
        password: state.signup.password,
        email: state.signup.email,
        is_artist: state.signup.is_artist,
        is_curator: state.signup.is_curator
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
