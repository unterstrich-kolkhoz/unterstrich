module.exports = {
  getLocal: localStorage.getItem.bind(localStorage),
  setLocal: localStorage.setItem.bind(localStorage),
  removeLocal: localStorage.removeItem.bind(localStorage)
};
