// Forms Utility Module
class FormsManager {
    constructor() {
        this.initializeEventListeners();
    }

    // Initialize form-related event listeners
    initializeEventListeners() {
        // Form toggle links
        document.addEventListener('click', (e) => {
            if (e.target.textContent === 'Register here') {
                e.preventDefault();
                this.showRegisterForm();
            } else if (e.target.textContent === 'Login here') {
                e.preventDefault();
                this.showLoginForm();
            }
        });
    }

    // Show register form
    showRegisterForm() {
        this.hideAllForms();
        const registerForm = document.getElementById('registerForm');
        if (registerForm) {
            registerForm.classList.remove('hidden');
            registerForm.style.display = 'block';
        }
        this.hideMessage();
    }

    // Show login form
    showLoginForm() {
        this.hideAllForms();
        const loginForm = document.getElementById('loginForm');
        if (loginForm) {
            loginForm.classList.remove('hidden');
            loginForm.style.display = 'block';
        }
        this.hideMessage();
    }

    // Hide all forms
    hideAllForms() {
        const forms = ['loginForm', 'registerForm'];
        forms.forEach(formId => {
            const form = document.getElementById(formId);
            if (form) {
                form.classList.add('hidden');
                form.style.display = 'none';
            }
        });
    }

    // Reset all forms
    resetAllForms() {
        const forms = ['loginForm', 'registerForm', 'passwordChangeForm'];
        forms.forEach(formId => {
            const form = document.getElementById(formId);
            if (form) form.reset();
        });
    }

    // Focus on first input in a form
    focusFirstInput(formId) {
        const form = document.getElementById(formId);
        if (form) {
            const firstInput = form.querySelector('input, textarea');
            if (firstInput) firstInput.focus();
        }
    }

    // Validate form inputs
    validateForm(formId) {
        const form = document.getElementById(formId);
        if (!form) return false;

        const inputs = form.querySelectorAll('input[required], textarea[required]');
        let isValid = true;

        inputs.forEach(input => {
            if (!input.value.trim()) {
                isValid = false;
                input.classList.add('error');
            } else {
                input.classList.remove('error');
            }
        });

        return isValid;
    }

    // Clear form validation errors
    clearValidationErrors(formId) {
        const form = document.getElementById(formId);
        if (form) {
            const inputs = form.querySelectorAll('.error');
            inputs.forEach(input => input.classList.remove('error'));
        }
    }

    // Safe message hiding
    hideMessage() {
        try {
            if (window.MessageManager) {
                window.MessageManager.hide();
            } else {
                const messageElement = document.getElementById('message');
                if (messageElement) {
                    messageElement.classList.add('hidden');
                }
            }
        } catch (error) {
            console.error('Error hiding message:', error);
        }
    }
}

// Export for use in other modules
window.FormsManager = FormsManager;
