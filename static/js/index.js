function getCookie(name) {
  let cookieValue = null;
  if (document.cookie && document.cookie !== "") {
    const cookies = document.cookie.split(";");
    for (let i = 0; i < cookies.length; i++) {
      const cookie = cookies[i].trim();
      // Does this cookie string begin with the name we want?
      if (cookie.substring(0, name.length + 1) === name + "=") {
        cookieValue = decodeURIComponent(cookie.substring(name.length + 1));
        break;
      }
    }
  }
  return cookieValue;
}

// export default setCookie;
// setCookie("access_token", "your_access_token", 7); // Set cookie for 7 days

// const accessToken = getCookie("access_token");

window.onload(checkAuth());

function checkAuth() {
  const accessToken = getCookie("access_token");
  if (!accessToken) {
    alert("Please login!");
    window.location.href = "login.html";
  }
}
