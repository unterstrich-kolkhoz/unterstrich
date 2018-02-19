const html = require("choo/html");

const page = require("../lib/page");

module.exports = function(state, emit) {
  function tab(name) {
    return () => emit("settingsTab", name);
  }

  function update(key) {
    return e => emit("updateSettings", { key, value: e.target.value });
  }

  function updateBool(key) {
    return e => emit("updateSettings", { key, value: e.target.checked });
  }

  function submitInfo() {
    emit("submitSettings");
  }

  function infoTab() {
    if (state.settings.tab != "info") return null;
    return html`
      <div class="settings-contents">
        <h4>Info</h4>
          <input id="settings-email"
                 type="email"
                 placeholder="email"
                 value=${state.settings.info.email}
                 required
                 onchange=${update("email")}>
          <input id="settings-firstname"
                 placeholder="Name"
                 value=${state.settings.info.name}
                 onchange=${update("name")}>
      <label for="settings-artist">Artist?</label>
      <input id="settings-artist"
             type="checkbox"
             ${state.settings.info.is_artist ? "checked" : ""}
              onchange=${updateBool("is_artist")}>
      <label for="settings-curator">Curator?</label>
      <input id="settings-curator"
             type="checkbox"
             ${state.settings.info.is_curator ? "checked" : ""}
             onchange=${updateBool("is_curator")}>
      </div>
    `;
  }

  function addressTab() {
    if (state.settings.tab != "address") return null;
    return html`
      <div class="settings-contents">
        <h4>Address</h4>
        <input id="settings-line1"
               placeholder="Address Line 1"
               value=${state.settings.info.line1}
               onchange=${update("line1")}>
        <input id="settings-line1"
               placeholder="Address Line 1"
               value=${state.settings.info.line2}
               onchange=${update("line2")}>
        <input id="settings-zip"
               placeholder="Zip Code"
               value=${state.settings.info.zip}
               onchange=${update("zip")}>
        <input id="settings-city"
               placeholder="City"
               value=${state.settings.info.city}
               onchange=${update("city")}>
        <input id="settings-state"
               placeholder="State"
               value=${state.settings.info.state}
               onchange=${update("state")}>
        <input id="settings-country"
               placeholder="Country"
               value=${state.settings.info.country}
               onchange=${update("country")}>
      </div>
    `;
  }

  function socialTab() {
    if (state.settings.tab != "social") return null;
    return html`
      <div class="settings-contents">
        <h4>Social</h4>
        <input id="settings-website"
               placeholder="Website"
               value=${state.settings.info.website}
               onchange=${update("website")}>
        <input id="settings-ello"
               placeholder="Ello"
               value=${state.settings.info.ello}
               onchange=${update("ello")}>
        <input id="settings-github"
               placeholder="Github"
               value=${state.settings.info.github}
               onchange=${update("github")}>
      </div>
    `;
  }

  return page(html`
    <div>
      <h3>Settings</h3>
      <div class="settings-tabs">
        <button onclick=${tab("info")}>General Info</button>
        <button onclick=${tab("address")}>Address</button>
        <button onclick=${tab("social")}>Social</button>
      </div>
        ${infoTab()}
        ${addressTab()}
        ${socialTab()}
        <button type="submit"
                onclick=${submitInfo}>
          Set
        </button>
    </div>
  `);
};
