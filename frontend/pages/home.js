const html = require("choo/html");

const style = require("../lib/style");

module.exports = function(state, emit) {
  return html`
    <body class=${style}>
      <div class="welcome">
        <h1><span id="logo">_</span> Unterstrich</h1>
      </div>
      <div class="right">
        ${greet()}
      </div>
      <p>We foster artistic expression.</p>
      <p>We help you get paid.</p>
      <p>We make it easy.</p>
    </body>
  `;

  function goTo(route) {
    return () => emit("pushState", route);
  }

  function greet() {
    if (state.login.token) {
      return html`<p>Hello, ${state.login.username}!</p>`;
    }
    return html`
      <div>
        <button onclick=${goTo("/login")}>Login</button>
        <button onclick=${goTo("/signup")}>Sign Up</button>
      </div>
    `;
  }
};
