const stepContainer = document.getElementById("step-container");
const step1Form = document.getElementById("step1-form");
const step2Form = document.getElementById("step2-form");
const step3Form = document.getElementById("step3-form");
const passwordInput = document.getElementById("password");
const confirmPasswordInput = document.getElementById("confirmPassword");
const emailInput = document.getElementById("email");
const stepMap = {
  1: "verifyemail",
  2: "verifyotp",
  3: "reset",
  4: "login",
};
let token = "";
// Live validation functions
const validateEmail = () => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (emailRegex.test(emailInput.value)) {
    emailInput.classList.remove("error");
    return true;
  } else {
    emailInput.classList.add("error");
    return false;
  }
};

const validatePassword = () => {
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
  if (passwordInput.value === confirmPasswordInput.value) {
    confirmPasswordInput.classList.remove("error");
    return true;
  } else {
    confirmPasswordInput.classList.add("error");
    return false;
  }
};

// Live validation on input events
emailInput.addEventListener("input", validateEmail);
passwordInput.addEventListener("input", validatePassword);
confirmPasswordInput.addEventListener("input", validateConfirmPassword);

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

// Function to handle API calls and step transitions
const handleStep = async (stepNumber, formData) => {
  try {
    const jsonData = formDataToJSON(formData);
    const apiUrl = `/v1/forgotpassword/${stepMap[stepNumber]}`;
    const response = await fetch(apiUrl, {
      method: "POST",
      body: jsonData,
      headers: { Authorization: `Bearer ${token}` },
    });

    if (response.ok) {
      const data = await response.json();
      token = data["token"];
      const currentForm = stepContainer.querySelector(
        `.step-form:nth-child(${stepNumber})`
      );
      currentForm.classList.remove("active");
      currentForm.classList.add("hidden");
      alert(data["message"]);
      if (stepMap[stepNumber + 1] != "login") {
        const nextForm = stepContainer.querySelector(
          `.step-form:nth-child(${stepNumber + 1})`
        );
        nextForm.classList.remove("hidden");
        nextForm.classList.add("active");
      } else {
        window.location.href = "login.html";
      }
    } else {
      // Show error message
      // Implement your preferred error handling mechanism (e.g., snackbar)
      console.error("API error:", response.statusText);
      // Example: showErrorSnackbar("Something went wrong. Please try again.");
    }
  } catch (error) {
    console.error("Error during API call:", error);
    // Handle network errors or other issues
  }
};

// Step 1 form submission
step1Form.addEventListener("submit", (event) => {
  event.preventDefault();
  const formData = new FormData(step1Form);
  handleStep(1, formData);
});

// Step 2 form submission
step2Form.addEventListener("submit", (event) => {
  event.preventDefault();
  const formData = new FormData(step2Form);
  handleStep(2, formData);
});

// Step 2 form submission
step3Form.addEventListener("submit", (event) => {
  event.preventDefault();
  const formData = new FormData(step3Form);
  handleStep(3, formData);
});

// Initially show Step 1
step1Form.classList.add("active");
