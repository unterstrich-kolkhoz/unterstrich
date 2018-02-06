const html = require("choo/html");

const style = require("../lib/style");

module.exports = function(state, emit) {
  if (!state.username) emit("pushState", "/");

  return html`
    <body class=${style}>
      <div class="welcome">
        <h1><span id="logo">_</span> ${state.username}</h1>
      </div>
    </body>
  `;
};
