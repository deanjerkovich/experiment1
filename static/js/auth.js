// Authentication Module
class AuthManager {
    constructor() {
        this.currentUser = null;
        this.initializeEventListeners();
    }

    // Initialize all authentication-related event listeners
    initializeEventListeners() {
        // Wait for DOM to be ready
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => this.setupEventListeners());
        } else {
            this.setupEventListeners();
        }
    }

    // Setup event listeners after DOM is ready
    setupEventListeners() {
        const loginForm = document.getElementById('loginForm');
        const registerForm = document.getElementById('registerForm');
        const passwordChangeForm = document.getElementById('passwordChangeForm');

        if (loginForm) {
            loginForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.handleLogin();
            });
        }

        if (registerForm) {
            registerForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.handleRegistration();
            });
        }

        if (passwordChangeForm) {
            passwordChangeForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.handlePasswordChange();
            });
        }
    }

    // Handle user login
    async handleLogin() {
        const username = document.getElementById('loginUsername').value;
        const password = document.getElementById('loginPassword').value;

        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password }),
            });

            let data;
            try {
                data = await response.json();
            } catch (jsonError) {
                // If response is not JSON, get the text
                const text = await response.text();
                console.error('Non-JSON response:', text);
                this.showMessage('Login failed: ' + text, 'error');
                return;
            }

            if (data.success) {
                this.currentUser = data.data;
                this.showProfile(data.data);
                this.showMessage('Login successful!', 'success');
            } else {
                this.showMessage(data.message || 'Login failed', 'error');
            }
        } catch (error) {
            console.error('Login error:', error);
            this.showMessage('An error occurred during login', 'error');
        }
    }

    // Handle user registration
    async handleRegistration() {
        const username = document.getElementById('registerUsername').value;
        const email = document.getElementById('registerEmail').value;
        const password = document.getElementById('registerPassword').value;

        try {
            const response = await fetch('/api/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, email, password }),
            });

            const data = await response.json();

            if (data.success) {
                this.showMessage(data.message, 'success');
                // Switch to login form after successful registration
                setTimeout(() => {
                    this.showLoginForm();
                    document.getElementById('loginUsername').value = data.data.username || '';
                    document.getElementById('loginPassword').focus();
                }, 1500);
                document.getElementById('registerForm').reset();
            } else {
                this.showMessage(data.message || 'Registration failed', 'error');
            }
        } catch (error) {
            console.error('Registration error:', error);
            this.showMessage('An error occurred during registration', 'error');
        }
    }

    // Handle password change
    async handlePasswordChange() {
        const currentPassword = document.getElementById('currentPassword').value;
        const newPassword = document.getElementById('newPassword').value;
        const confirmPassword = document.getElementById('confirmPassword').value;

        if (newPassword !== confirmPassword) {
            this.showMessage('New passwords do not match', 'error');
            return;
        }

        if (newPassword.length < 6) {
            this.showMessage('New password must be at least 6 characters', 'error');
            return;
        }

        try {
            const response = await fetch('/api/change-password', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ currentPassword, newPassword }),
            });

            const data = await response.json();

            if (data.success) {
                this.showMessage(data.message, 'success');
                document.getElementById('passwordChangeForm').reset();
            } else {
                if (response.status === 401) {
                    this.showMessage('Authentication required. Please login again.', 'error');
                    this.forceShowLogin();
                } else {
                    this.showMessage(data.message || 'Password change failed', 'error');
                }
            }
        } catch (error) {
            console.error('Password change error:', error);
            this.showMessage('An error occurred while changing password', 'error');
        }
    }

    // Handle user logout
    async handleLogout() {
        try {
            const response = await fetch('/api/logout', { method: 'POST' });
            const data = await response.json();

            if (data.success) {
                this.showMessage(data.message, 'success');
                this.currentUser = null;
                this.forceShowLogin();
            } else {
                this.forceShowLogin();
            }
        } catch (error) {
            console.error('Logout error:', error);
            this.showMessage('An error occurred during logout', 'error');
            this.forceShowLogin();
        }
    }

    // Show user profile
    showProfile(user) {
        console.log('Showing profile for user:', user);
        this.currentUser = user;

        // Hide forms and show profile
        this.hideAllForms();
        this.showProfileSection();
        this.updateProfileInfo(user);
        this.startSessionMonitoring();
    }

    // Hide all authentication forms
    hideAllForms() {
        const elements = ['loginForm', 'registerForm'];
        elements.forEach(id => {
            const element = document.getElementById(id);
            if (element) {
                element.classList.add('hidden');
                element.style.display = 'none';
            }
        });
    }

    // Show profile section
    showProfileSection() {
        const profile = document.getElementById('profile');
        const topNav = document.getElementById('topNav');

        if (profile) {
            profile.style.display = 'block';
            profile.classList.remove('hidden');
        }

        if (topNav) {
            topNav.classList.remove('hidden');
            topNav.style.display = 'flex';
        }

        // Initialize and reset navigation to default state
        if (window.navigationManager) {
            console.log('Calling navigationManager.show() and reset()');
            window.navigationManager.show();
            window.navigationManager.reset();
        } else {
            console.log('navigationManager not available');
        }

        // Show profile content by default
        this.showProfileContent();
    }

    // Show profile content (default view)
    showProfileContent() {
        // This will be handled by the NavigationManager
        // For now, just ensure the profile content is visible
        const profileContent = document.getElementById('profileContent');
        if (profileContent) {
            profileContent.classList.add('active');
            profileContent.style.display = 'block';
        }
    }

    // Update profile information display
    updateProfileInfo(user) {
        const profileInfo = document.getElementById('profileInfo');
        if (profileInfo) {
            profileInfo.innerHTML = `
                <p><strong>Username:</strong> ${user.username}</p>
                <p><strong>Email:</strong> ${user.email}</p>
                <p><strong>Created:</strong> ${new Date(user.created).toLocaleDateString()}</p>
            `;
        }
    }

    // Force show login screen
    forceShowLogin() {
        console.log('Force showing login screen');
        
        // Hide all authenticated content
        const elements = ['profile', 'topNav', 'registerForm'];
        elements.forEach(id => {
            const element = document.getElementById(id);
            if (element) {
                element.classList.add('hidden');
                element.style.display = 'none';
            }
        });

        // Show login form
        const loginForm = document.getElementById('loginForm');
        if (loginForm) {
            loginForm.classList.remove('hidden');
            loginForm.style.display = 'block';
        }

        // Reset forms and clear messages
        this.resetForms();
        this.hideMessage();
        this.currentUser = null;
    }

    // Reset all forms
    resetForms() {
        const forms = ['loginForm', 'registerForm', 'passwordChangeForm'];
        forms.forEach(formId => {
            const form = document.getElementById(formId);
            if (form) form.reset();
        });
    }

    // Show login form
    showLoginForm() {
        this.hideAllForms();
        const loginForm = document.getElementById('loginForm');
        if (loginForm) {
            loginForm.classList.remove('hidden');
            loginForm.style.display = 'block';
        }
    }

    // Start monitoring session validity
    startSessionMonitoring() {
        // Check session every minute
        setInterval(async () => {
            if (this.currentUser) {
                try {
                    const response = await fetch('/api/profile');
                    if (!response.ok || !(await response.json()).success) {
                        console.log('Session expired, redirecting to login');
                        this.forceShowLogin();
                        this.showMessage('Your session has expired. Please login again.', 'error');
                    }
                } catch (error) {
                    console.log('Authentication check failed, redirecting to login');
                    this.forceShowLogin();
                    this.showMessage('Authentication failed. Please login again.', 'error');
                }
            }
        }, 60000);
    }

    // Safe message display that doesn't depend on MessageManager
    showMessage(text, type = 'success') {
        try {
            if (window.MessageManager) {
                window.MessageManager.show(text, type);
            } else {
                // Fallback message display
                const messageElement = document.getElementById('message');
                if (messageElement) {
                    messageElement.textContent = text;
                    messageElement.className = `message ${type}`;
                    messageElement.classList.remove('hidden');
                    
                    // Auto-hide after 5 seconds
                    setTimeout(() => {
                        messageElement.classList.add('hidden');
                    }, 5000);
                }
            }
        } catch (error) {
            console.error('Error showing message:', error);
            // Last resort: console log
            console.log(`[${type.toUpperCase()}] ${text}`);
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

    // Check if user is already logged in on page load
    async checkExistingSession() {
        try {
            const response = await fetch('/api/profile');
            if (response.ok) {
                const data = await response.json();
                if (data.success) {
                    this.showProfile(data.data);
                    return;
                }
            }
        } catch (error) {
            console.log('No existing session');
        }
        
        // No valid session, show login
        this.forceShowLogin();
    }
}

// Export for use in other modules
window.AuthManager = AuthManager;
