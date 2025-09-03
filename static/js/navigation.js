// Navigation Module
class NavigationManager {
    constructor() {
        this.initializeNavigation();
        this.initializeEventListeners();
    }

    // Initialize navigation structure
    initializeNavigation() {
        // Set up navigation items
        this.navItems = [
            { id: 'profile', icon: '👤', label: 'Profile', handler: () => this.showProfileContent() },
            { id: 'base64', icon: '🔐', label: 'Base64 Tools', handler: () => this.showBase64Content() },
            { id: 'logout', icon: '🚪', label: 'Logout', handler: () => this.handleLogout() }
        ];

        // Create navigation menu
        this.createNavigationMenu();
    }

    // Create the navigation menu HTML
    createNavigationMenu() {
        const topNav = document.getElementById('topNav');
        if (!topNav) {
            console.log('topNav element not found in createNavigationMenu');
            return;
        }

        console.log('Creating navigation menu...');

        // Clear any existing content and event listeners
        topNav.innerHTML = '';
        if (this.handleNavClick) {
            topNav.removeEventListener('click', this.handleNavClick);
            console.log('Removed existing click handler');
        }

        topNav.innerHTML = this.navItems.map((item, index) => `
            <div class="nav-item ${index === 0 ? 'active' : ''}" 
                 data-nav-id="${item.id}">
                <span class="nav-icon">${item.icon}</span>
                ${item.label}
            </div>
        `).join('');

        console.log('Navigation menu HTML created with', this.navItems.length, 'items');
        
        // Set up event listeners after creating the menu
        this.setupClickHandlers();
    }

        // Set up click handlers for navigation items
    setupClickHandlers() {
        const topNav = document.getElementById('topNav');
        if (!topNav) {
            console.log('topNav not found in setupClickHandlers');
            return;
        }

        console.log('Setting up navigation click handlers...');

        // Remove any existing listeners to avoid duplicates
        if (this.handleNavClick) {
            topNav.removeEventListener('click', this.handleNavClick);
            console.log('Removed existing click handler');
        }
        
        // Add the click handler
        this.handleNavClick = (e) => {
            console.log('Navigation click detected:', e.target);
            const target = e.target.closest('.nav-item');
            if (!target) {
                console.log('No nav-item found in click target');
                return;
            }
            const navId = target.getAttribute('data-nav-id');
            console.log('Navigation ID clicked:', navId);
            switch (navId) {
                case 'profile':
                    console.log('Showing profile content');
                    this.showProfileContent();
                    break;
                case 'base64':
                    console.log('Showing base64 content');
                    this.showBase64Content();
                    break;
                case 'logout':
                    console.log('Handling logout');
                    this.handleLogout();
                    break;
                default:
                    console.log('Unknown navigation ID:', navId);
                    break;
            }
        };
        
        topNav.addEventListener('click', this.handleNavClick);
        console.log('Navigation click handlers set up successfully');
    }

    // Initialize event listeners
    initializeEventListeners() {
        // Add keyboard shortcuts
        document.addEventListener('keydown', (e) => {
            if (e.ctrlKey && e.shiftKey && e.key === 'L') {
                console.log('Emergency logout triggered');
                this.handleLogout();
            }
        });
    }

    // Show profile content
    showProfileContent() {
        this.updateActiveNavItem('profile');
        this.showContent('profileContent');
    }

    // Show base64 content
    showBase64Content() {
        this.updateActiveNavItem('base64');
        this.showContent('base64Content');
    }

    // Update active navigation item
    updateActiveNavItem(activeId) {
        this.navItems.forEach(item => {
            const navItem = document.querySelector(`[data-nav-id="${item.id}"]`);
            if (navItem) {
                if (item.id === activeId) {
                    navItem.classList.add('active');
                } else {
                    navItem.classList.remove('active');
                }
            }
        });
    }

    // Show specific content section
    showContent(contentId) {
        // Hide all content sections
        const allContent = document.querySelectorAll('.nav-content');
        allContent.forEach(content => {
            content.classList.remove('active');
            content.style.display = 'none';
        });

        // Show selected content
        const selectedContent = document.getElementById(contentId);
        if (selectedContent) {
            selectedContent.classList.add('active');
            selectedContent.style.display = 'block';
        }
    }

    // Handle logout
    handleLogout() {
        if (window.authManager) {
            window.authManager.handleLogout();
        }
    }

    // Show navigation menu
    show() {
        console.log('NavigationManager.show() called');
        const topNav = document.getElementById('topNav');
        if (topNav) {
            // Clear any existing content first
            topNav.innerHTML = '';
            // Ensure navigation menu is created
            this.createNavigationMenu();
            topNav.classList.remove('hidden');
            topNav.style.display = 'flex';
            console.log('Navigation menu shown');
        } else {
            console.log('topNav element not found');
        }
    }

    // Hide navigation menu
    hide() {
        const topNav = document.getElementById('topNav');
        if (topNav) {
            // Clear content and event listeners when hiding
            topNav.innerHTML = '';
            if (this.handleNavClick) {
                topNav.removeEventListener('click', this.handleNavClick);
            }
            topNav.classList.add('hidden');
            topNav.style.display = 'none';
        }
    }

    // Reset navigation to default state
    reset() {
        console.log('NavigationManager.reset() called');
        // Recreate the navigation menu to ensure clean state
        this.createNavigationMenu();
        this.updateActiveNavItem('profile');
        this.showContent('profileContent');
        console.log('Navigation reset to profile');
    }
}

// Export for use in other modules
window.NavigationManager = NavigationManager;
