const html = require("choo/html");

const style = require("../lib/style");

const getLocalItem = require("../lib/storage").getLocalItem;

module.exports = function(state, emit) {
  if (!state.userInfo) {
    emit("getUserInfo");
    return null;
  }

  return html`
    <body class=${style}>
      <div class="welcome">
        <h1><span id="logo">_</span> ${state.userInfo.username}</h1>
        ${artistBadge()}
        ${curatorBadge()}
      </div>
    </body>
  `;

  function artistBadge() {
    if (state.userInfo.is_artist) {
      return html`<span class="badge">Artist</span>`;
    }
    return null;
  }

  function curatorBadge() {
    if (state.userInfo.is_curator) {
      return html`<span class="badge">Curator</span>`;
    }
    return null;
  }
};
