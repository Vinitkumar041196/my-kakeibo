const signupForm = document.getElementById("signup-form");
const loginForm = document.getElementById("login-form");

const signupEmailInput = document.getElementById("signup-email");
const loginEmailInput = document.getElementById("login-email");

const signupPasswordInput = document.getElementById("signup-password");
const loginPasswordInput = document.getElementById("login-password");
const confirmPasswordInput = document.getElementById("confirmPassword");

const showSignupButton = document.getElementById("show-signup");
const showLoginButton = document.getElementById("show-login");

const formTitle = document.getElementById("form-title");
const loginFormTitle = "Log In To Better Money Management!";
const signupFormTitle = "Sign Up To Better Money Management!";
// Initially show the signup form
loginForm.classList.remove("hidden");
formTitle.textContent = loginFormTitle;
document.getElementById("show-login").classList.add("highlighted");

// Button click events
showSignupButton.addEventListener("click", () => {
  signupForm.classList.remove("hidden");
  // signupForm.classList.add("highlighted");
  loginForm.classList.add("hidden");
  formTitle.textContent = signupFormTitle;
  document.getElementById("show-login").classList.remove("highlighted");
  document.getElementById("show-signup").classList.add("highlighted");
});

showLoginButton.addEventListener("click", () => {
  signupForm.classList.add("hidden");
  loginForm.classList.remove("hidden");
  // loginForm.classList.add("highlighted");
  formTitle.textContent = loginFormTitle;
  document.getElementById("show-signup").classList.remove("highlighted");
  document.getElementById("show-login").classList.add("highlighted");
});

function setCookie(name, value, minutes) {
  let expires = "";
  if (minutes) {
    const date = new Date();
    date.setTime(date.getTime() + minutes * 60 * 60 * 1000);
    expires = "; expires=" + date.toUTCString();
  }
  document.cookie =
    name + "=" + encodeURIComponent(value) + expires + "; path=/";
}

// Live validation functions
const validateEmail = (e) => {
  const emailInput = e.target;
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (emailRegex.test(emailInput.value)) {
    emailInput.classList.remove("error");
    return true;
  } else {
    emailInput.classList.add("error");
    return false;
  }
};

const validatePassword = (e) => {
  const passwordInput = e.target;
  const passwordRegex =
    /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@\-_])[^\s]{8,20}$/;
  //console.log(passwordRegex.test(passwordInput.value));
  if (passwordRegex.test(passwordInput.value)) {
    passwordInput.classList.remove("error");
    return true;
  } else {
    passwordInput.classList.add("error");
    return false;
  }
};

const validateConfirmPassword = () => {
  if (signupPasswordInput.value === confirmPasswordInput.value) {
    confirmPasswordInput.classList.remove("error");
    return true;
  } else {
    confirmPasswordInput.classList.add("error");
    return false;
  }
};

function formDataToJSON(formData) {
  const obj = {};
  for (const [key, value] of formData.entries()) {
    // Handle nested objects/arrays in FormData
    if (Array.isArray(value)) {
      obj[key] = value.map((v) => (v instanceof File ? v : v.toString()));
    } else if (value instanceof File) {
      obj[key] = value; // Keep File objects as-is
    } else {
      obj[key] = value.toString(); // Convert other values to strings
    }
  }
  return JSON.stringify(obj);
}

// Live validation on input events
loginEmailInput.addEventListener("input", validateEmail);
signupEmailInput.addEventListener("input", validateEmail);
loginPasswordInput.addEventListener("input", validatePassword);
signupPasswordInput.addEventListener("input", validatePassword);
confirmPasswordInput.addEventListener("input", validateConfirmPassword);

// Form submission handling
signupForm.addEventListener("submit", (event) => {
  event.preventDefault(); // Prevent default form submission

  // Check all validations again before submitting
  const formData = new FormData(signupForm);
  const jsonData = formDataToJSON(formData);

  fetch("/v1/signup", {
    method: "POST",
    body: jsonData,
  })
    .then((response) => {
      if (response.ok) {
        // Handle successful signup, e.g., redirect to login page
        window.location.href = "login.html";
      } else {
        // Handle signup error, e.g., display an error message
        alert("Signup failed. Please try again.");
      }
    })
    .catch((error) => {
      console.error("Error during signup:", error);
    });
});

loginForm.addEventListener("submit", (event) => {
  event.preventDefault(); // Prevent default form submission

  const formData = new FormData(loginForm);
  const jsonData = formDataToJSON(formData);

  fetch("/v1/login", {
    method: "POST",
    body: jsonData,
  })
    .then((response) => {
      if (response.ok) {
        const data = response.json();
        console.log(data);
        setCookie("access_token", data["token"], 15 ); // Set cookie for 15 mins
        // Handle successful signup, e.g., redirect to login page
        window.location.href = "index.html";
      } else {
        // Handle signup error, e.g., display an error message
        alert("Signup failed. Please try again.");
      }
    })
    .catch((error) => {
      console.error("Error during signup:", error);
    });
});
